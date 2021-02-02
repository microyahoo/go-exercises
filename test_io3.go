package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
)

func main() {
	inputReader, inputWriter, _ := os.Pipe()
	outputReader, outputWriter, _ := os.Pipe()

	io.Copy(inputWriter, bytes.NewReader([]byte("hello world\n")))

	stdin := inputReader
	stdout := outputWriter
	stderr := outputWriter

	var attr = os.ProcAttr{
		Dir: "/tmp",
		Env: nil,
		Files: []*os.File{
			stdin,
			stdout,
			stderr,
		},
		Sys: nil,
	}

	var (
		process         *os.Process
		startProcessErr error
	)
	if runtime.GOOS == "linux" {
		process, startProcessErr = os.StartProcess("/usr/bin/ls", []string{"ls"}, &attr)
	} else if runtime.GOOS == "darwin" {
		process, startProcessErr = os.StartProcess("/bin/ls", []string{"ls"}, &attr)
	}
	if startProcessErr != nil {
		panic(startProcessErr)
	}

	defer func() {
		if releaseProcessErr := process.Release(); releaseProcessErr != nil {
			panic(releaseProcessErr)
		}
	}()

	outputWriter.Close() // <-- add this line

	fmt.Println("-------1-------")
	var output bytes.Buffer
	io.Copy(&output, outputReader)
	fmt.Println(output)
}
