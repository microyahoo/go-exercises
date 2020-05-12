package main

import (
	"fmt"
	"reflect"
	_ "unsafe"
)

type MyInt int

type A struct {
	b B
}
type B struct {
	c C
}

type C int

func main() {
	var x MyInt = 7
	v := reflect.ValueOf(x)
	fmt.Println(v)
	fmt.Printf("v.CanSet() = %v\n", v.CanSet())
	fmt.Printf("v.CanInterface() = %v\n", v.CanInterface())
	fmt.Printf("v.CanAddr() = %v\n", v.CanAddr())
	fmt.Printf("v.Interface() = %v\n", v.Interface())
	fmt.Printf("v.Kind() = %v\n", v.Kind())
	fmt.Printf("v.Type() = %v\n", v.Type())

	fmt.Println("-----------------1------------")
	var z []A
	var y []map[string]int
	m := reflect.ValueOf(z)
	// fmt.Println(m)
	fmt.Printf("m.Kind() = %v\n", m.Kind())
	fmt.Printf("m.Type() = %v\n", m.Type())
	n := reflect.ValueOf(y)
	fmt.Printf("n.Kind() = %v\n", n.Kind())
	fmt.Printf("n.Type() = %v\n", n.Type())
	if n.Kind() == reflect.Slice {
		fmt.Printf("n.Type().Elem() = %v\n", n.Type().Elem())
		fmt.Printf("n.Type().Elem().Elem() = %v\n", n.Type().Elem().Elem())
		fmt.Printf("n.Type() = %v\n", n.Type())
	}
	if m.Kind() == reflect.Slice {
		fmt.Printf("m.Type().Elem() = %v\n", m.Type().Elem())
		fmt.Printf("m.Type() = %v\n", m.Type())
	}

	fmt.Println("------------2-----------------")
	a := reflect.ValueOf(&x)
	fmt.Println(a)
	fmt.Printf("a.CanSet() = %v\n", a.CanSet())
	fmt.Printf("a.CanInterface() = %v\n", a.CanInterface())
	fmt.Printf("a.CanAddr() = %v\n", a.CanAddr())
	fmt.Printf("a.Interface() = %v\n", a.Interface())
	fmt.Printf("a.Kind() = %v\n", a.Kind())
	fmt.Printf("a.Type() = %v\n", a.Type())
	fmt.Printf("a.Elem() = %v\n", a.Elem())
	fmt.Printf("a.Elem().CanSet() = %v\n", a.Elem().CanSet())
	fmt.Printf("a.Elem().Type() = %v\n", a.Elem().Type())

	fmt.Println("---------------1--------------")
	var f float64 = 3.4
	v2 := reflect.ValueOf(f)
	o := v2.Interface().(float64)
	fmt.Println(o)
}
