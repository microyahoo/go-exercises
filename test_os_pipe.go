package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"time"
)

func main() {
	inR, inW, _ := os.Pipe()
	outR, outW, _ := os.Pipe()
	done := make(chan struct{})

	process, _ := os.StartProcess("/bin/sh", nil, &os.ProcAttr{
		Files: []*os.File{inR, outW, outW}})

	notifyC := make(chan struct{})
	go func() {
		writer := bufio.NewWriter(inW)
		for i := 0; i < 3; i++ {
			time.Sleep(time.Second)
			writer.WriteString("cat test1.go\n")
			writer.Flush()
		}
		inW.Close()
		outW.Close()
		notifyC <- struct{}{}
	}()
	go func() {
		// https://go-review.googlesource.com/c/go/+/292130/1/src/os/exec/exec.go
		var (
			buf = make([]byte, 1000)
			err error
			nr  int
		)
		// outR.SetReadDeadline(time.Now().Add(time.Millisecond * 99))
		for {
			// A deadline is an absolute time after which I/O operations fail with an error
			// instead of blocking. The deadline applies to all future and pending I/O, not
			// just the immediately following call to Read or Write. After a deadline has
			// been exceeded, the connection can be refreshed by setting a deadline in the future.
			outR.SetReadDeadline(time.Now().Add(time.Millisecond * 99))
			nr, err = outR.Read(buf)
			// fmt.Println(nr, err)
			if err != nil && errors.Is(err, os.ErrDeadlineExceeded) {
				select {
				case <-notifyC:
					break
				default:
					if nr == 0 {
						continue
					}
					err = nil
				}
			}
			break
		}
		fmt.Println(string(buf))
		// scanner := bufio.NewScanner(outR)
		// for scanner.Scan() {
		// 	fmt.Println(scanner.Text())
		// }
		process.Signal(os.Kill)
		done <- struct{}{}
		fmt.Println("finish")
	}()

	// process.Wait()
	<-done

}
