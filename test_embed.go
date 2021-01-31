package main

import (
	"embed"
	"fmt"
	"net/http"
)

func main() {
	//go:embed content.tmpl
	//go:embed *.go
	var content embed.FS

	fmt.Println("vim-go")
	http.ListenAndServe(":8080", http.FileServer(http.FS(content)))
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(content))))
}

// go1.16beta1 run test_embed.go
