package main

import (
	"fmt"

	"github.com/mjc-gh/kredis-go"
)

func main() {
	kredis.SetConfiguration("shared", "ns", "redis://localhost:6379/2")
	kredis.SetCommandLogging(true)

	l, _ := kredis.NewIntegerList("key")
	l.Append(1, 2, 3)
	l.Prepend(8, 9)
	llen, _ := l.Length()

	fmt.Printf("list LLEN = %d\n", llen)

	elements := make([]int, 3)
	n, _ := l.Elements(elements)

	fmt.Printf("n = %d (%v)\n", n, elements)

	lastTwo := make([]int, 2)
	n, _ = l.Elements(lastTwo, kredis.WithRangeStart(3))

	fmt.Printf("n = %d (%v)\n", n, lastTwo)

	l.Clear()
}
