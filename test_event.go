package main

import (
	"fmt"
	"time"

	e "github.com/microyahoo/go-exercises/events"
)

const HELLO_WORLD = "helloWorld"

func main() {
	dispatcher := e.NewEventDispatcher()
	listener := e.NewEventListener(myEventListener)
	dispatcher.AddEventListener(HELLO_WORLD, listener)

	time.Sleep(time.Second * 2)
	//dispatcher.RemoveEventListener(HELLO_WORLD, listener)

	dispatcher.DispatchEvent(e.NewEvent(HELLO_WORLD, nil))
}

func myEventListener(event e.Event) {
	fmt.Println(event.Type, event.Object, event.Target)
}
