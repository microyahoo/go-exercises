package main

import (
	"fmt"
	"os"
	"text/template"
)

func main() {
	s1, _ := template.ParseFiles("header.tmpl", "content.tmpl", "footer.tmpl")
	s1.ExecuteTemplate(os.Stdout, "header", nil)
	fmt.Println("--------------------1--------------------")
	s1.ExecuteTemplate(os.Stdout, "content", nil)
	fmt.Println("--------------------2--------------------")
	s1.ExecuteTemplate(os.Stdout, "footer", nil)
	fmt.Println("--------------------3--------------------")
	fmt.Println("hello")
	s1.Execute(os.Stdout, nil)
}
