package file

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// DeleteStats aggregates deletion-phase statistics.
type DeleteStats struct {
	Scanned     atomic.Int64
	Deleted     atomic.Int64
	Failed      atomic.Int64
	DirsRemoved atomic.Int64
}

// ────────────────────────────────────────────────────────────────────
// Exact mode: partitioned hash set (zero false positives)
// ────────────────────────────────────────────────────────────────────

// WalkDestinationPaths traverses dstDir and writes every regular-file
// relative path into the partition writer. It also collects directory
// paths for later empty-directory cleanup. The scanned counter is
// incremented for every file encountered so a progress reporter can
// track the walk.
func WalkDestinationPaths(dstDir string, pw *PartitionWriter, scanned *atomic.Int64) (dirs []string, err error) {
	err = filepath.WalkDir(dstDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			fmt.Fprintf(os.Stderr, "[WALK-ERR] %s: %v\n", path, walkErr)
			return nil
		}
		if d.IsDir() {
			if path != dstDir {
				dirs = append(dirs, path)
			}
			return nil
		}
		if !d.Type().IsRegular() {
			return nil
		}
		rel, relErr := filepath.Rel(dstDir, path)
		if relErr != nil {
			fmt.Fprintf(os.Stderr, "[WALK-ERR] rel %s: %v\n", path, relErr)
			return nil
		}
		scanned.Add(1)
		return pw.Add(filepath.Clean(rel))
	})
	return dirs, err
}

// DeleteExtrasExact processes each partition sequentially: loads the
// source bucket into a hash set, streams the destination bucket, and
// concurrently deletes every destination file that is NOT in the source
// set. This guarantees zero false positives / false negatives.
func DeleteExtrasExact(dstDir string, srcPW, dstPW *PartitionWriter, partitions, concurrency int, dryRun bool, stats *DeleteStats) error {
	for p := 0; p < partitions; p++ {
		srcSet, err := LoadPartitionSet(srcPW.PartitionPath(p))
		if err != nil {
			return fmt.Errorf("load source partition %d: %w", p, err)
		}

		extras, err := findExtrasInPartition(dstPW.PartitionPath(p), srcSet)
		srcSet = nil
		if err != nil {
			return fmt.Errorf("diff partition %d: %w", p, err)
		}

		if len(extras) > 0 {
			deleteFiles(dstDir, extras, concurrency, dryRun, stats)
		}
	}
	return nil
}

func findExtrasInPartition(dstPartPath string, srcSet map[string]struct{}) ([]string, error) {
	f, err := os.Open(dstPartPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var extras []string
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		rel := scanner.Text()
		if rel == "" {
			continue
		}
		if _, ok := srcSet[rel]; !ok {
			extras = append(extras, rel)
		}
	}
	return extras, scanner.Err()
}

func deleteFiles(dstDir string, relPaths []string, concurrency int, dryRun bool, stats *DeleteStats) {
	ch := make(chan string, concurrency*2)
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for rel := range ch {
				abs := filepath.Join(dstDir, rel)
				if dryRun {
					log.Printf("[DRY-RUN] would delete: %s", abs)
					stats.Deleted.Add(1)
					continue
				}
				if err := os.Remove(abs); err != nil {
					if !os.IsNotExist(err) {
						stats.Failed.Add(1)
						fmt.Fprintf(os.Stderr, "[DELETE-ERR] %s: %v\n", abs, err)
					}
				} else {
					stats.Deleted.Add(1)
				}
			}
		}()
	}
	for _, rel := range relPaths {
		ch <- rel
	}
	close(ch)
	wg.Wait()
}

// ────────────────────────────────────────────────────────────────────
// Approx mode: bloom filter (fast, bounded memory, ~FPR% may survive)
// ────────────────────────────────────────────────────────────────────

// DeleteExtrasApprox walks dstDir, checks every regular file against
// the bloom filter, and removes files that are (probably) not in the
// source list. Returns collected directory paths for empty-dir cleanup.
func DeleteExtrasApprox(dstDir string, bf *BloomFilter, concurrency int, dryRun bool, stats *DeleteStats) (dirs []string, err error) {
	pathCh := make(chan string, concurrency*2)

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for absPath := range pathCh {
				rel, relErr := filepath.Rel(dstDir, absPath)
				if relErr != nil {
					stats.Failed.Add(1)
					fmt.Fprintf(os.Stderr, "[DELETE-ERR] rel path for %s: %v\n", absPath, relErr)
					continue
				}
				rel = filepath.Clean(rel)

				if bf.Contains([]byte(rel)) {
					continue
				}

				if dryRun {
					log.Printf("[DRY-RUN] would delete: %s", absPath)
					stats.Deleted.Add(1)
					continue
				}

				if rmErr := os.Remove(absPath); rmErr != nil {
					if !os.IsNotExist(rmErr) {
						stats.Failed.Add(1)
						fmt.Fprintf(os.Stderr, "[DELETE-ERR] %s: %v\n", absPath, rmErr)
					}
				} else {
					stats.Deleted.Add(1)
				}
			}
		}()
	}

	err = filepath.WalkDir(dstDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			fmt.Fprintf(os.Stderr, "[DELETE-WALK] %s: %v\n", path, walkErr)
			return nil
		}
		if d.IsDir() {
			if path != dstDir {
				dirs = append(dirs, path)
			}
			return nil
		}
		if !d.Type().IsRegular() {
			return nil
		}
		stats.Scanned.Add(1)
		pathCh <- path
		return nil
	})

	close(pathCh)
	wg.Wait()
	return dirs, err
}

// ────────────────────────────────────────────────────────────────────
// Common helpers
// ────────────────────────────────────────────────────────────────────

// CleanEmptyDirs tries to remove directories deepest-first.
// os.Remove on a non-empty directory is a no-op that returns an error
// which we silently ignore.
func CleanEmptyDirs(dirs []string) int64 {
	sort.Slice(dirs, func(i, j int) bool {
		return len(dirs[i]) > len(dirs[j])
	})
	var removed int64
	for _, dir := range dirs {
		if err := os.Remove(dir); err == nil {
			removed++
		}
	}
	return removed
}

// StartDeleteProgressReporter periodically logs delete-phase progress.
func StartDeleteProgressReporter(stats *DeleteStats, start time.Time, interval time.Duration) func() {
	ticker := time.NewTicker(interval)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				elapsed := time.Since(start)
				log.Printf("[delete-progress] elapsed=%s scanned=%d deleted=%d failed=%d",
					elapsed.Round(time.Second),
					stats.Scanned.Load(), stats.Deleted.Load(), stats.Failed.Load())
			case <-done:
				return
			}
		}
	}()
	return func() {
		ticker.Stop()
		close(done)
	}
}
