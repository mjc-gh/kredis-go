package main

import (
	"github.com/mjc-gh/kredis-go"
)

func main() {
	kredis.SetConfiguration("shared", "ns", "redis://localhost:6379/2")
	kredis.SetCommandLogging(true)

	cntr, _ := kredis.NewCounter("counter")
	cntr.Increment(1)
	cntr.Increment(2)
	cntr.Decrement(3)
	cntr.Value()
}
