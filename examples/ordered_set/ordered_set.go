package main

import (
	"fmt"

	"github.com/mjc-gh/kredis-go"
)

func main() {
	kredis.SetConfiguration("shared", "ns", "redis://localhost:6379/2")
	kredis.EnableDebugLogging()

	oset, _ := kredis.NewStringOrderedSet("ranks", 4)
	oset.Append("a", "b", "c")
	oset.Prepend("d", "e")
	members, _ := oset.Members()

	fmt.Printf("%d %v\n", oset.Size(), members)

	fmt.Println(oset.Includes("d"))
	fmt.Println(oset.Includes("c"))
}
