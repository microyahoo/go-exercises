package main

import (
	"fmt"
	"runtime"
	"time"
)

type test struct {
	Num int
}

func main() {
	runtime.GOMAXPROCS(1)

	t := &test{}

	go func() {
		for i := 0; i < 100; i++ {
			time.Sleep(300 * time.Microsecond)
			num := t.Num
			time.Sleep(300 * time.Microsecond)
			fmt.Printf("First: i = %d\n", i)
			t.Num = num + 1
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			num := t.Num
			time.Sleep(300 * time.Microsecond)
			fmt.Printf("Second: i = %d\n", i)
			t.Num = num + 1
		}
	}()

	time.Sleep(10 * time.Second) // 等待goroutine完成
	fmt.Println(t.Num)
}
