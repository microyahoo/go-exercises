package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/go-homedir"
)

const (
	DivisionNum    = 10
	ConcurrencyNum = 2
	FilterDir      = "go/src/github.com/microyahoo/go-exercises"
	IncludeFrom    = "fig_more_complicated_ward.go"
)

func splitIncludeFromFile() (files []string, err error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	f := fmt.Sprintf("%s/%s", filepath.Join(homeDir, FilterDir), IncludeFrom)
	fmt.Println(f)
	if DivisionNum <= 1 {
		return []string{f}, nil
	}
	info, err := os.Stat(f)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if info.IsDir() || info.Size() == 0 {
		return nil, nil
	}
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	lines, err := lineCounter(file)
	if err != nil {
		return nil, err
	}
	if lines == 0 {
		return nil, nil
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	step := int64(math.Ceil(float64(lines) / DivisionNum))
	// step := int64(lines) / DivisionNum
	var i int64
	reader := bufio.NewReader(file)
	for i = 0; i < DivisionNum; i++ {
		name := fmt.Sprintf("%s-%s.%d", f, uuid.String(), i+1)
		newFile, err := os.Create(name)
		if err != nil {
			return nil, err
		}
		var j int64
		for j < step {
			s, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					fmt.Println(err)
					break
				}
				return nil, err
			}
			_, err = newFile.WriteString(s)
			if err != nil {
				return nil, err
			}
			j++
		}
		newFile.Close()
		files = append(files, name)
	}

	return files, nil
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func main() {
	fmt.Println(os.Executable())
	files, err := splitIncludeFromFile()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, f := range files {
		fmt.Println(f)
	}
	concurrencyNum := ConcurrencyNum
	if concurrencyNum > 1 && concurrencyNum > len(files) {
		concurrencyNum = len(files)
	} else if concurrencyNum <= 1 {
		concurrencyNum = 1
	}
	jobQueue := make(chan string, concurrencyNum)
	tokenBucket := make(chan struct{}, concurrencyNum)
	go func() {
		for _, file := range files {
			fmt.Printf("----> send file: %s with %s\n", file, time.Now())
			jobQueue <- file
		}
		close(jobQueue)
	}()

	group := &sync.WaitGroup{}
	group.Add(len(files))
	for cmd := range jobQueue {
		select {
		case tokenBucket <- struct{}{}:
		}
		go func(c string) {
			defer group.Done()
			defer func() {
				<-tokenBucket
			}()
			t := time.Second * time.Duration(rand.Intn(30))
			fmt.Printf("\t<--- Start to handle %s: need %s\n", c, t)
			time.Sleep(t)
			fmt.Printf("\t\t<--- Succeed to run command: %s: %s\n", c, time.Now())
		}(cmd)
	}
	group.Wait()
	fmt.Println("Finish")
}
