package main

import (
	"fmt"
)

func main() {
	switch num := 70; {
	case num < 50:
		fmt.Printf("%d is lesser than 50\n", num)
		fallthrough
	case num < 100:
		fmt.Printf("%d is lesser than 100\n", num)
		fallthrough
	case num < 200:
		fmt.Printf("%d is lesser than 200\n", num)
		fallthrough
	case num < 300:
		fmt.Printf("%d is lesser than 300", num)
	}
	fmt.Println("\n", checkOsdAlertSent("full"))
	bits := 4
	fmt.Println(1<<bits - 1)

	for i := 0; i < 5; i++ {
		fmt.Println(i, isEntry(i))
	}
	for i := 0; i < 5; i++ {
		fmt.Println(i, isEntry2(i))
	}

	var x float64
	x = 349473.479916
	fmt.Println(int64(x / (1024 * 1024)))
	x = 970558111.491300
	fmt.Println(int64(x / (1024 * 1024)))

}

func checkOsdAlertSent(alertValue string) bool {
	switch alertValue {
	case "full":
	case "backfillFull":
		return true
	case "nearFull":
		return true
	}
	return false
}

func isEntry(i int) bool {
	switch {
	case i == 0:
		return true
	case i == 1:
		return true
	case i == 2:
		return true
	default:
		return false
	}
}

func isEntry2(i int) bool {
	switch {
	case i == 0, i == 1, i == 2:
		return true
	default:
		return false
	}
}
