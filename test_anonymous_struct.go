package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	nesCarts := []string{"Battletoads", "Mega Man 1", "Clash at Demonhead"}
	numberOfCarts := len(nesCarts)

	anonymousStruct := struct {
		NESCarts      []string
		numberOfCarts int
	}{
		nesCarts,
		numberOfCarts,
	}
	data, err := json.Marshal(anonymousStruct)
	fmt.Printf("data=%v, err=%v\n", data, err)
	fmt.Printf("data=%v\n", string(data))
}
