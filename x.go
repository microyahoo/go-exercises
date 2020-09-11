package main

import (
	"fmt"
)

func main() {
	fmt.Println(3)
}

// Run starts the matching process off
func (m *March) Run() error {
	m.init()

	srcDepth := fs.Config.MaxDepth
	if srcDepth < 0 {
		srcDepth = fs.MaxLevel
	}
	dstDepth := srcDepth
	if filter.Active.Opt.DeleteExcluded {
		dstDepth = fs.MaxLevel
	}

	var mu sync.Mutex // Protects vars below
	var jobError error
	var errCount int

	// Start some directory listing go routines
	var wg sync.WaitGroup         // sync closing of go routines
	var traversing sync.WaitGroup // running directory traversals
	in := make(chan listDirJob, fs.Config.Checkers)
	for i := 0; i < fs.Config.Checkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-m.Ctx.Done():
					return
				case job, ok := <-in:
					if !ok {
						return
					}
					jobs, err := m.processJob(job)
					if err != nil {
						mu.Lock()
						// Keep reference only to the first encountered error
						if jobError == nil {
							jobError = err
						}
						errCount++
						mu.Unlock()
					}
					if len(jobs) > 0 {
						traversing.Add(len(jobs))
						go func() {
							// Now we have traversed this directory, send these
							// jobs off for traversal in the background
							for _, newJob := range jobs {
								select {
								case <-m.Ctx.Done():
									// discard job if finishing
									traversing.Done()
								case in <- newJob:
								}
							}
						}()
					}
					traversing.Done()
				}
			}
		}()
	}

	// Start the process
	traversing.Add(1)
	in <- listDirJob{
		srcRemote: m.Dir,
		srcDepth:  srcDepth - 1,
		dstRemote: m.Dir,
		dstDepth:  dstDepth - 1,
		noDst:     m.NoCheckDest,
	}
	go func() {
		// when the context is cancelled discard the remaining jobs
		<-m.Ctx.Done()
		for range in {
			traversing.Done()
		}
	}()
	traversing.Wait()
	close(in)
	wg.Wait()

	if errCount > 1 {
		return errors.Wrapf(jobError, "march failed with %d error(s): first error", errCount)
	}
	return jobError
}
