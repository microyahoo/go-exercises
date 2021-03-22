package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func fatal(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
	os.Exit(1)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", "echo hello world && sleep 5")
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		fatal("no pipe: %v", err)
	}
	start := time.Now()

	if err = cmd.Start(); err != nil {
		fatal("start failed: %v", err)
	}

	go func() {
		var stdout string
		buf := bufio.NewReader(stdoutPipe)
		for {
			line, err := buf.ReadString('\n')
			if len(line) > 0 {
				stdout = stdout + line + "\n"
			}
			if err != nil {
				if !strings.HasPrefix(stdout, "hello world") {
					fatal("wrong output: %q", stdout)
				}
				return
			}
		}
	}()

	err = cmd.Wait()
	d := time.Since(start)

	fmt.Println(d)

	if err != nil {
		exiterr := err.(*exec.ExitError)
		status := exiterr.Sys().(syscall.WaitStatus)
		if status.ExitStatus() == 0 {
			fatal("wrong exit status: %v", status.ExitStatus())
		}
	}

	if d.Seconds() >= 3 {
		fatal("Cancelation took too long: %v", d)
	}
	fmt.Println("Success!")
}
