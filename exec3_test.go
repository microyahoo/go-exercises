package main

import (
	"bufio"
	"context"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCmd(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", "echo hello world && sleep 5")
	assert.NotNil(t, cmd)

	stdoutPipe, err := cmd.StdoutPipe()
	assert.Nil(t, err)
	assert.NotNil(t, stdoutPipe)

	start := time.Now()

	err = cmd.Start()
	assert.Nil(t, err)

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
		require.True(t, ok)
		status, ok := exiterr.Sys().(syscall.WaitStatus)
		require.True(t, ok)
		assert.NotEqual(t, 0, status.ExitStatus())
	}
	assert.True(t, strings.HasPrefix(stdout, "hello world"), "Stdout: %v", stdout)

	assert.True(t, d.Seconds() < 3, "Duration was %v", d)
}
