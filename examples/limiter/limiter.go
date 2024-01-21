package main

import (
	"github.com/mjc-gh/kredis-go"
)

func main() {
	kredis.SetConfiguration("shared", "ns", "redis://localhost:6379/2")
	kredis.EnableDebugLogging()

	limiter, _ := kredis.NewLimiter("limiter", 5)

	for i := 0; i < 4; i++ {
		limiter.Poke()
	}

	limiter.IsExceeded()
	limiter.Reset()
}
