package main

/*
extern char* NewGoString(char*);
extern void FreeGoString(char*);
extern void PrintGoString(char*);

static void printString(char* s) {
	char* gs = NewGoString(s);
	PrintGoString(gs);
	FreeGoString(gs);
}
*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
)

type ObjectId int32

var refs struct {
	sync.Mutex
	objs map[ObjectId]interface{}
	next ObjectId
}

func NewObjectId(obj interface{}) ObjectId {
	refs.Lock()
	defer refs.Unlock()

	id := refs.next
	refs.next++
	refs.objs[id] = obj

	return id
}

func (id ObjectId) IsNil() bool {
	return id == 0
}

func (id ObjectId) Get() interface{} {
	refs.Lock()
	defer refs.Unlock()
	return refs.objs[id]
}

func (id *ObjectId) Free() interface{} {
	refs.Lock()
	defer refs.Unlock()

	obj := refs.objs[*id]
	delete(refs.objs, *id)
	*id = 0

	return obj
}

//export NewGoString
func NewGoString(s *C.char) *C.char {
	gs := C.GoString(s)
	id := NewObjectId(gs)
	fmt.Println("NewGoString", id)
	x := (*C.char)(unsafe.Pointer(uintptr(id)))
	fmt.Println(x)
	return x
}

//export FreeGoString
func FreeGoString(s *C.char) {
	id := ObjectId(uintptr(unsafe.Pointer(s)))
	fmt.Println("FreeGoString", id)
	id.Free()
}

//export PrintGoString
func PrintGoString(s *C.char) {
	id := ObjectId(uintptr(unsafe.Pointer(s)))
	fmt.Println("PrintGoString", id)
	gs := id.Get().(string)
	fmt.Println(gs)
}

func main() {
	C.printString(C.CString("hello world"))
}

func init() {
	refs.Lock()
	defer refs.Unlock()

	refs.objs = make(map[ObjectId]interface{})
	refs.next = 1000
}
