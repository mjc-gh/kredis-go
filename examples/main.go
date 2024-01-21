package main

import (
	"fmt"
	"time"

	"github.com/mjc-gh/kredis-go"
)

func main() {
	kredis.SetConfiguration("shared", "ns", "redis://localhost:6379/2")
	kredis.EnableDebugLogging()

	slot, err := kredis.NewSlot("slot", 3)
	if err != nil {
		panic(err)
	}

	slot.Reserve()
	fmt.Println(slot.Taken())

	// some time later when Redis is down...
	time.Sleep(10 * time.Second)

	fmt.Println(slot.Taken())
	n, err := slot.TakenResult() // 0,
	fmt.Println(n, err)

	fmt.Println(err)
}
