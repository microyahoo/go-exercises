package main

import "runtime"
import "log"

func main() {
	test()
}

func test() {
	test2()
}

func test2() {
	pc, file, line, ok := runtime.Caller(2)
	log.Println(pc)
	log.Println(file)
	log.Println(line)
	log.Println(ok)
	f := runtime.FuncForPC(pc)
	log.Println(f.Name())

	pc, file, line, ok = runtime.Caller(0)
	log.Println(pc)
	log.Println(file)
	log.Println(line)
	log.Println(ok)
	f = runtime.FuncForPC(pc)
	log.Println(f.Name())
	pc, file, line, ok = runtime.Caller(1)
	log.Println(pc)
	log.Println(file)
	log.Println(line)
	log.Println(ok)
	f = runtime.FuncForPC(pc)
	log.Println(f.Name())
}

// 2019/10/29 23:52:21 17384715
// 2019/10/29 23:52:21 /Users/xsky/go/src/github.com/microyahoo/go-exercises/test_runtime3.go
// 2019/10/29 23:52:21 7
// 2019/10/29 23:52:21 true
// 2019/10/29 23:52:21 main.main

// 2019/10/29 23:52:21 17385279
// 2019/10/29 23:52:21 /Users/xsky/go/src/github.com/microyahoo/go-exercises/test_runtime3.go
// 2019/10/29 23:52:21 23
// 2019/10/29 23:52:21 true
// 2019/10/29 23:52:21 main.test2

// 2019/10/29 23:52:21 17384720
// 2019/10/29 23:52:21 /Users/xsky/go/src/github.com/microyahoo/go-exercises/test_runtime3.go
// 2019/10/29 23:52:21 11
// 2019/10/29 23:52:21 true
// 2019/10/29 23:52:21 main.test
