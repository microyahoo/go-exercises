package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

// Index handler
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// Hello handler
func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
	fmt.Fprintf(w, "hello, parameters %v!\n", ps)
}

// World handler
func World(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "hello world!!!\n")
}

// WildcardWorld handler
func WildcardWorld(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello wildcard world!!!\n")
}

func wildcardApi(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello wildcard api!!!\n")
}

func api(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello api!!!\n")
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)
	router.GET("/hell/world", World)
	router.GET("/helx/*world", WildcardWorld)
	router.GET("/xxx/*path", wildcardApi)
	router.GET("/xxx/api", api)

	log.Fatal(http.ListenAndServe(":8080", router))
}
