package main

import (
	_ "crypto/sha256" // register hash algorithm
	"fmt"

	"github.com/opencontainers/go-digest"
)

func main() {
	dgst := digest.FromString("sha256:1a4b70b3af7b490abe4edc67b3292fb5698b10f94b552a0681cde5d95fed599b")
	fmt.Println(dgst)
	fmt.Println(dgst.Encoded())
	dgst = digest.FromString("sha256:1f3ed09f")
	fmt.Println(dgst)
	fmt.Println(dgst.Encoded())
}
