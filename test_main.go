package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/microyahoo/go-exercises/closure"
	"github.com/microyahoo/go-exercises/http/database"
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

	fmt.Println("-------begin to test http database-----")
	db := database.Database{"shoes": 50, "socks": 5}
	mux := http.NewServeMux()
	mux.Handle("/list", http.HandlerFunc(db.List))
	mux.HandleFunc("/price", db.Price)
	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}
