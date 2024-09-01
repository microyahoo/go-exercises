package main

import (
	"context"
	"fmt"

	gopsutilhost "github.com/shirou/gopsutil/v3/host"
)

func main() {
	temperatures, _ := gopsutilhost.SensorsTemperaturesWithContext(context.Background())
	for _, temp := range temperatures {
		fmt.Println(temp)
	}
}
