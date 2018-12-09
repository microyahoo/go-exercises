package main

import "fmt"
import "time"

func main() {
	fmt.Println("Hello, playground")
	// ticker := time.NewTicker(2 * time.Second)
	ticker := time.NewTicker(200 * time.Millisecond)
	t := time.Now()
	test := make(chan string)
	close(test) //UNCOMMENT FOR FAILURE

	i := 0
	for {
		select {
		case _, ok := <-test:
			fmt.Println(i, "-------: ", time.Since(t))
			t = time.Now()
			i++
			if i == 10 {
				return
			}
			if !ok {
				// test = nil  //A nil channel is never ready for communication. So each time you run into a closed channel, you can nil that channel ensuring it is never selected again.
			}
		case <-ticker.C:
			fmt.Println(i, ": ", time.Since(t))
			t = time.Now()
			i++
			if i == 10 {
				return
			}
		}
	}
}
