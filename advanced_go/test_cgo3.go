package main

/*
#include<errno.h>

static int div(int a, int b) {
	if (b == 0) {
		errno = EINVAL;
		return 0;
	}
	return a/b;
}

static void  noreturn() {}

int sum(int a, int b) {
	return a+b;
}

*/
import "C"
import (
	"fmt"
)

func main() {
	v0, err0 := C.div(2, 1)
	fmt.Println(v0, err0)

	v1, err1 := C.div(2, 0)
	fmt.Println(v1, err1)

	v2, err2 := C.noreturn()
	fmt.Println(v2, err2)
	fmt.Printf("%#v\n", v2)

	fmt.Println(C.sum(1, 2))
}
