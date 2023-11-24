package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	if os.Getenv("LISTEN_PID") == strconv.Itoa(os.Getpid()) {
		fmt.Println(os.Getenv("LISTEN_PID"))
		fmt.Println(os.Getenv("LISTEN_FDNAMES"))
		fmt.Println(os.Getenv("LISTEN_FDS"))
		// systemd run
		f := os.NewFile(3, "from systemd")
		l, err := net.FileListener(f)
		if err != nil {
			log.Fatal(err)
		}
		http.Serve(l, nil)
	} else {
		// manual run
		log.Fatal(http.ListenAndServe(":8080", nil))
	}
}
