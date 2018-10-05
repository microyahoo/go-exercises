package main

import (
    "fmt"
)

func main() {
    f := squares()
    fmt.Printf("%T\n", f)
    fmt.Println(f())
    fmt.Println(f())
    fmt.Println(f())
    fmt.Println(f())
}

func squares() func() int {
    var x int
    fmt.Printf("**%d\n", x)
    return func() int {
        x++
        return x * x
    }
}
