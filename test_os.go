package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"runtime"
)

func main() {
	fmt.Println(runtime.GOOS)
	fmt.Println(os.Getpagesize())
	fmt.Println(0x1000)
	fmt.Println(math.Pow(2, 12))
}

func capture(f func()) (string, error) {
	rescueOut := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", fmt.Errorf("capture: %s", err.Error())
	}
	os.Stdout = w
	f()
	out, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatalf("capture: %s", err.Error())
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("capture: %s", err.Error())
	}
	os.Stdout = rescueOut
	return string(out), nil
}

func capture2(f func()) (string, error) {
	rescueOut := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", fmt.Errorf("capture: %s", err.Error())
	}
	os.Stdout = w
	outC := make(chan string)
	go func() {
		out, err := ioutil.ReadAll(r)
		if err != nil {
			log.Fatalf("capture: %s", err.Error())
		}
		outC <- string(out)
	}()
	f()
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("capture: %s", err.Error())
	}
	os.Stdout = rescueOut
	// 在stdout没有全部写到pipe之前，让outC阻塞住
	return <-outC, nil
}

// https://blog.scnace.me/post/go-os-pipe/
