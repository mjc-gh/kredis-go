package main

import "github.com/mjc-gh/kredis-go"

func main() {
	kredis.SetConfiguration("shared", "ns", "redis://localhost:6379/2")
	kredis.EnableDebugLogging()

	slot, _ := kredis.NewSlot("slot", 3)
	slot.Reserve()
	slot.IsAvailable()
	slot.Taken()
	slot.Reserve()
	slot.Reserve()

	slot.Reserve()
	slot.IsAvailable()
	slot.Taken()
}
