package main

import "fmt"
import "time"
import "log"
import "net/http"
import "io/ioutil"
import _ "sync"

type Memo struct {
	requests chan request
}

type Func func(key string) (interface{}, error)

type entry struct {
	res   result
	ready chan struct{}
}

type result struct {
	value interface{}
	err   error
}

type request struct {
	key      string
	response chan<- result
}

func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f)
	return memo
}

func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	for req := range memo.requests {
		e := cache[req.key]
		if e == nil {
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key)
		}
		go e.deliver(req.response)
	}
}

func (e *entry) deliver(response chan<- result) {
	<-e.ready
	response <- e.res
}

func (e *entry) call(f Func, key string) {
	e.res.value, e.res.err = f(key)
	//NOTE broadcast
	close(e.ready)
}

func (memo Memo) Get(key string) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response}
	res := <-response
	return res.value, res.err
}

func (memo *Memo) Close() {
	close(memo.requests)
}

func httpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func incomingURLs() <-chan string {
	ch := make(chan string)
	go func() {
		for _, url := range []string{
			"https://baidu.com",
			"https://godoc.org",
			"https://github.com",
			"http://gopl.io",
			"https://baidu.com",
			"https://godoc.org",
			"https://github.com",
			"http://gopl.io",
		} {
			ch <- url
		}
		close(ch)
	}()
	return ch
}

func main() {
	// var n sync.WaitGroup
	m := New(httpGetBody)
	for url := range incomingURLs() {
		start := time.Now()
		value, err := m.Get(url)
		if err != nil {
			log.Print(err)
		}
		fmt.Printf("%s, %s, %d bytes\n", url, time.Since(start), len(value.([]byte)))
	}
}
