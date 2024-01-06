package main

import (
	"fmt"

	"github.com/mjc-gh/kredis-go"
)

func main() {
	kredis.SetConfiguration("shared", "ns", "redis://localhost:6379/2")
	kredis.SetCommandLogging(true)

	uniq, _ := kredis.NewFloatUniqueList("uniq", 5)
	uniq.Append(3.14, 2.718)
	uniq.Prepend(1.1)
	llen, _ := uniq.Length()

	fmt.Printf("list LLEN = %d\n", llen)

	elements := make([]float64, 3)
	n, _ := uniq.Elements(elements)

	fmt.Printf("n = %d (%v)\n", n, elements)

	uniq.Clear()
}
