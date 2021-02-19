package main

//https://github.com/golang/go/issues/23019

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	// "testing"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", "echo hello world && sleep 5")

	stdoutPipe, err := cmd.StdoutPipe()
	fmt.Println(err == nil)
	// assert.Nil(t, err)
	fmt.Println(stdoutPipe != nil)
	// assert.NotNil(t, stdoutPipe)

	start := time.Now()

	err = cmd.Start()
	fmt.Println(err == nil)
	// assert.Nil(t, err)

	var stdout string
	go func() {
		buf := bufio.NewReader(stdoutPipe)
		for {
			line, err := buf.ReadString('\n')
			if len(line) > 0 {
				stdout = stdout + line + "\n"
			}
			if err != nil {
				return
			}
		}
	}()

	err = cmd.Wait()
	d := time.Since(start)

	if err != nil {
		exiterr, ok := err.(*exec.ExitError)
		// require.True(t, ok)
		fmt.Println(ok)
		status, ok := exiterr.Sys().(syscall.WaitStatus)
		// require.True(t, ok)
		fmt.Println(ok)
		// assert.NotEqual(t, 0, status.ExitStatus())
		fmt.Println(0 != status.ExitStatus())
	}
	fmt.Println(stdout)
	// assert.True(t, strings.HasPrefix(stdout, "hello world"), "Stdout: %v", stdout)
	fmt.Println(strings.HasPrefix(stdout, "hello world"), "Stdout: ", stdout)

	// assert.True(t, d.Seconds() < 3, "Duration was %v", d)
	fmt.Println(d.Seconds() < 3, "Duration was: ", d)
}
