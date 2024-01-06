package main

import (
	"fmt"
	"time"

	"github.com/mjc-gh/kredis-go"
)

func main() {
	kredis.SetConfiguration("shared", "ns", "redis://localhost:6379/2")
	kredis.SetCommandLogging(true)

	flag, _ := kredis.NewFlag("flag")
	fmt.Println(flag.IsMarked())
	flag.Mark()
	fmt.Println(flag.IsMarked())
	flag.Remove()

	flag.Mark(kredis.WithFlagExpiry("1s"))
	fmt.Println(flag.IsMarked())

	time.Sleep(2 * time.Second)

	fmt.Println(flag.IsMarked())

	flag.Mark(kredis.WithFlagExpiry("1s"))
	flag.SoftMark(kredis.WithFlagExpiry("10s"))
	fmt.Println(flag.IsMarked())

	time.Sleep(2 * time.Second)

	fmt.Println(flag.IsMarked())
}
