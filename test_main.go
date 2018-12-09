package main

import (
	"bufio"
	"bytes"
	"fmt"
	// "github.com/astaxie/beego"
	"io"
	_ "log"
	_ "net/http"
	"os"
	"sync"
	_ "time"

	"github.com/microyahoo/go-exercises/closure"
	_ "github.com/microyahoo/go-exercises/http/database"
)

func main() {
	f := closure.Squares()
	fmt.Printf("%T\n", f)
	fmt.Println(f())
	fmt.Println(f())
	fmt.Println(f())
	fmt.Println(f())

	fmt.Println("-----begin to test interface----")
	var r io.Reader
	fmt.Println(r)
	fmt.Println(r == nil)
	fmt.Printf("%#v\n", r)
	fmt.Printf("-io.Reader---1---%T\n\n", r)
	r = os.Stdin
	fmt.Println(r)
	fmt.Println(r == nil)
	fmt.Printf("%#v\n", r)
	fmt.Printf("-os.Stdin---2---%T\n\n", r)
	r = bufio.NewReader(r)
	//fmt.Println(r)
	fmt.Println(r == nil)
	//fmt.Printf("%#v\n", r)
	fmt.Printf("-bufio.NewReader---3---%T\n\n", r)
	r = new(bytes.Buffer)
	fmt.Println(r)
	fmt.Println(r == nil)
	fmt.Printf("%#v\n", r)
	fmt.Printf("-new(bytes.Buffer)---4---%T\n\n", r)

	fmt.Println("-------begin to test interface again-----")
	var r2 io.Reader
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return
	}
	r2 = tty
	fmt.Printf("++==%T\n\n", r2)
	var w io.Writer
	w = r2.(io.Writer)
	fmt.Printf("++==%T\n\n", w)

	fmt.Println("-------begin to test error-----")
	_, err = os.Open("/no/such/file")
	fmt.Println(err)
	fmt.Printf("%#v\n", err)

	fmt.Println("-------begin to test select-----")
	for j := 0; j < 10; j++ {
		ch := make(chan int, 2)
		for i := 0; i < 10; i++ {
			select {
			case x := <-ch:
				fmt.Println(x)
			case ch <- i:
			}
		}
		fmt.Println("-------------++++++++++++++-----------")
	}

	// fmt.Println("-------begin to test time.NewTicker-----")
	// ticker := time.NewTicker(time.Second)
	// defer ticker.Stop()
	// done := make(chan bool)
	// go func() {
	// 	time.Sleep(10 * time.Second)
	// 	done <- true
	// }()
	// for {
	// 	select {
	// 	case <-done:
	// 		fmt.Println("Done!")
	// 		return
	// 	case t := <-ticker.C:
	// 		fmt.Println("Current time: ", t)
	// 	}
	// }

	fmt.Println("-------begin to test concurrency-----")
	var (
		mu sync.Mutex
	)
	mu.Lock()
	mu.Unlock()

	// fmt.Println("-------begin to test http database-----")
	// db := database.Database{"shoes": 50, "socks": 5}
	// mux := http.NewServeMux()
	// mux.Handle("/list", http.HandlerFunc(db.List))
	// mux.HandleFunc("/price", db.Price)
	// log.Fatal(http.ListenAndServe("localhost:8080", mux))

	fmt.Println("-------begin to test beego -----")
	// beego.Run()

	type values map[string][]string
	v := values{}
	v["x"] = []string{"1", "2"}
	fmt.Println(v)
	for key, value := range v {
		fmt.Printf("%s = %v", key, value)
	}

	fmt.Println("-------begin to test make-----")
	ar := make([]int, 0)
	fmt.Println(ar)
	fmt.Printf("%v\n", ar)
	ar = append(ar, 1)
	ar = append(ar, 2)
	fmt.Println(ar)
	fmt.Printf("%v\n", ar)
	for i, a := range ar {
		fmt.Println(i)
		fmt.Println(a)
	}
	fmt.Println(ar)
	for a := range ar {
		fmt.Println(a)
	}
	cls := new(Cls)
	cls.Array = make([]string, 0)
	cls.Array = append(cls.Array, "x")
	cls.Array = append(cls.Array, "y")
	fmt.Printf("%v", cls)

	fmt.Println("\n--------------compare struct------------")
	c1 := &Code{"x", "y"}
	c2 := &Code{"x", "y"}
	fmt.Println(c1 == c2)
	fmt.Println("\n--------------type assertion------------")
	do(21)
	do("hello")
	do(true)
	x := 22
	do(&x)
	do(23.12)
	do(cls)

}

func do(i interface{}) {
	switch v := i.(type) {
	// case int:
	// case float64:
	// case float32, float64:
	// 	fmt.Printf("Twice %v is %v\n", v, v*2)
	case int, uint:
		fmt.Printf("Twice the int %d\n", v)
	case string:
		fmt.Printf("%q is %v bytes long\n", v, len(v))
	case *int:
		fmt.Println("*int type")
	case *Cls:
		fmt.Println(v.Array)
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}
}

type Cls struct {
	Array []string
}

type Code struct {
	a string
	b string
}
