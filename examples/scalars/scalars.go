package main

import (
	"fmt"
	"time"

	"github.com/mjc-gh/kredis-go"
)

func main() {
	kredis.SetConfiguration("shared", "ns", "redis://localhost:6379/2")
	kredis.EnableDebugLogging()

	k, _ := kredis.NewTime("sessionStart", kredis.WithExpiry("1ms"))
	_ = k.SetValue(time.Now()) // SET sessionStart

	time.Sleep(1 * time.Second)

	fmt.Println(k.Value())
	fmt.Println(k.TTL())
}
