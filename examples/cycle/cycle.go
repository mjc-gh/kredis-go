package main

import (
	"fmt"

	"github.com/mjc-gh/kredis-go"
)

func main() {
	kredis.SetConfiguration("shared", "ns", "redis://localhost:6379/2")
	kredis.SetCommandLogging(true)

	cycle, _ := kredis.NewCycle("cycle", []string{"ready", "set", "go"})
	cycle.Index()
	cycle.Next()
	cycle.Next()

	fmt.Println(cycle.Index())
	fmt.Println(cycle.Value())
}
