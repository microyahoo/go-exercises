package main

// https://github.com/golang/go/issues/23019

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", "echo hello world && sleep 5")

	stdoutPipe, err := cmd.StdoutPipe()
	fmt.Println(err == nil)
	fmt.Println(stdoutPipe != nil)

	start := time.Now()

	err = cmd.Start()
	fmt.Println(err == nil)

	var wg sync.WaitGroup
	var stdout string

	wg.Add(1)
	go func() {
		defer wg.Done()
		buf := bufio.NewReader(stdoutPipe)
		for {
			line, err := buf.ReadString('\n')
			if len(line) > 0 {
				stdout = stdout + line + "\n"
			}
			if err != nil {
				fmt.Println(err) // NOTE: read |0: file already closed
				return
			}
		}
	}()

	// https://golang.org/pkg/os/exec/#Cmd.StdoutPipe
	// Wait will close the pipe after seeing the command exit, so most callers need not close
	// the pipe themselves. It is thus incorrect to call Wait before all reads from the pipe have completed.
	err = cmd.Wait()
	wg.Wait()
	d := time.Since(start)

	if err != nil {
		exiterr, ok := err.(*exec.ExitError)
		fmt.Println(ok)
		status, ok := exiterr.Sys().(syscall.WaitStatus)
		fmt.Println(ok)
		fmt.Println(0 != status.ExitStatus())
	}
	fmt.Println(stdout)
	fmt.Println(strings.HasPrefix(stdout, "hello world"), "Stdout: ", stdout)

	fmt.Println(d.Seconds() < 3, "Duration was: ", d)
}
