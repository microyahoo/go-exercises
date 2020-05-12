package main

/*
extern int* getGoPtr();

static void Main() {
	int* p = getGoPtr();
	*p = 43;
}
*/
import "C"

func main() {
	C.Main()
}

//export getGoPtr
func getGoPtr() *C.int {
	return new(C.int)
}
