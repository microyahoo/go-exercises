package main

//blog: xiaorui.cc
import (
	"github.com/robfig/cron"
	"log"
	"time"
)

func main() {
	i := 0
	c := cron.New()
	// spec := "0 */1 * * * *"
	// spec := "@every 1s"
	spec := "@every 30s"
	// spec := "@daily"
	err := c.AddFunc(spec, func() {
		i++
		log.Println("start", i)
	})
	log.Println(err)
	c.Start()
	// log.Println("hello")
	for j := 0; j < 100; j++ {
		log.Printf("world: %d\n", j)
		time.Sleep(3 * time.Second)
	}
	select {} //阻塞主线程不退出

}
