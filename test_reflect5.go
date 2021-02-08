package main

import (
	"fmt"
	"reflect"
)

type Foo struct {
	Number        int32 `tag:"number"`
	Ptr           *Bar
	privateNumber int32
}

type Bar struct {
	privateNumber int32
	Number        int32
}

type Baz struct {
	Bar
}

type anInterface interface{}

func nilHasInvalidKindType() {
	meta := reflect.ValueOf(nil)
	fmt.Println("nilHasInvalidKindType:", meta.Kind() == reflect.Invalid) // true
}

func nilPtrValueOf() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in nilPtrValueOf:", r)
		}
	}()
	var f *Foo
	meta := reflect.ValueOf(f)
	ref := meta.Elem() // no panic
	fmt.Println("ref: ", ref)
	fmt.Println(
		"nilPtrValueOf:",
		meta.Kind() == reflect.Ptr,    // true
		ref.Kind() != reflect.Struct,  // true
		ref.Kind() == reflect.Invalid, // true
		ref.IsValid() != true)         // true
	// ref.FieldByName("Number") // panic because zero value

	f = &Foo{}
	meta = reflect.ValueOf(f)
	ref = meta.Elem() // no panic
	fmt.Println("ref: ", ref)
	fmt.Println(
		"nilPtrValueOf:",
		meta.Kind() == reflect.Ptr,    // true
		ref.Kind() != reflect.Struct,  // false
		ref.Kind() == reflect.Invalid, // false
		ref.IsValid() != true)         // false
	ref.FieldByName("Number") // not panic
}

func nilPtrTypeOf() {
	var f *Foo
	meta := reflect.TypeOf(f)
	ref := meta.Elem()
	// No panic! this is not suprising since it's supposed to be type info
	// but kinda confusing
	numField, _ := ref.FieldByName("Number") // no panic!
	numTag, _ := numField.Tag.Lookup("tag")
	fmt.Println(
		"nilPtrTypeOf:",
		// again, not suprising but kinda confusing
		// (compare with nilPtrValueOf)
		ref.Kind() != reflect.Invalid, // true
		ref.Kind() == reflect.Struct,  // true
		numTag == "number")            // true

}

func nilIsNotNil() {
	var f *Foo
	fmt.Println("nilIsNotNil:", f == nil) // true
	var ptr anInterface = f
	fmt.Println("nilIsNotNil:", ptr != nil) // true
	// let's pass a ptr to the ptr
	meta := reflect.ValueOf(&ptr)
	ref := meta.Elem()
	fmt.Println(
		"nilIsNotNil:",
		meta.Kind() == reflect.Ptr,          // true
		ref.Kind() != reflect.Struct,        // true
		ref.Elem().Kind() != reflect.Struct, // true
		ref.Elem(),
		ref.Kind() == reflect.Interface) // true
	// let's pass ptr by value
	meta = reflect.ValueOf(ptr)
	ref = meta.Elem()
	fmt.Println(
		"nilIsNotNil:",
		meta.Kind() == reflect.Ptr,    // true
		ref.Kind() != reflect.Struct,  // true
		ref.Kind() == reflect.Invalid) // true
}

func canSet() {
	var bar Bar
	var baz Baz
	// can't set if pass by value
	meta := reflect.ValueOf(bar)
	publicField := meta.FieldByName("Number")
	fmt.Println(
		"canSet (pass by value):",
		publicField.CanSet() != true) // true

	// must pass a reference
	meta = reflect.ValueOf(&bar)
	privateField := meta.Elem().FieldByName("privateNumber")
	publicField = meta.Elem().FieldByName("Number")
	fmt.Println(
		"canSet (pass by reference):",
		publicField.CanSet() == true, // true
		// but cannot set private fields
		privateField.CanSet() != true) // true

	// works the same on embedded fields
	meta = reflect.ValueOf(&baz)
	privateField = meta.Elem().FieldByName("privateNumber")
	publicField = meta.Elem().FieldByName("Number")
	fmt.Println(
		"canSet (embedded fields):",
		privateField.CanSet() != true, // true
		publicField.CanSet() == true)  // true
}

func canInterface() {
	num := 6
	meta := reflect.ValueOf(num)
	fmt.Println("canInterface:", meta.CanInterface() == true)

	meta = reflect.ValueOf(&num)
	fmt.Println("canInterface:", meta.CanInterface() == true)

	foo := Foo{}
	meta = reflect.ValueOf(&foo)
	fmt.Println("canInterface:", meta.CanInterface() == true)
	meta = meta.Elem()
	fmt.Println("canInterface:", meta.CanInterface() == true)
	publicField := meta.FieldByName("Number")
	fmt.Println("canInterface:", publicField.CanInterface() == true)
	privateField := meta.FieldByName("privateNumber")
	fmt.Println("canInterface:", privateField.CanInterface() != true)

	var fooPtr *Foo
	var ptr anInterface = fooPtr
	meta = reflect.ValueOf(ptr)
	fmt.Println("canInterface:", meta.CanInterface() == true)

	meta = reflect.ValueOf(&foo)
	meta = meta.Elem() // ptr to actual value
	publicField = meta.FieldByName("Number")
	ptrToField := publicField.Addr()
	fmt.Println("canInterface:", ptrToField.CanInterface() == true)
	privateField = meta.FieldByName("privateNumber")
	ptrToField = privateField.Addr()
	fmt.Println("canInterface:", ptrToField.CanInterface() != true)
}

func main() {
	nilHasInvalidKindType()
	nilPtrValueOf()
	nilPtrTypeOf()
	nilIsNotNil()
	canSet()
	canInterface()
}
