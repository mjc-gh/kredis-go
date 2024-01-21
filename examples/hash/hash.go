package main

import (
	"fmt"

	"github.com/mjc-gh/kredis-go"
)

func main() {
	kredis.SetConfiguration("shared", "ns", "redis://localhost:6379/2")
	kredis.EnableDebugLogging()

	dogs := map[string]kredis.KredisJSON{
		"ollie": *kredis.NewKredisJSON(`{"weight":9.72}`),
		"leo":   *kredis.NewKredisJSON(`{"weight":23.33}`),
	}

	hash, _ := kredis.NewJSONHash("pets")
	hash.Update(dogs)

	val, ok := hash.Get("ollie")
	fmt.Printf("%v %v\n", ok, val)

	keys, _ := hash.Keys()
	vals, _ := hash.Values()
	entries, _ := hash.Entries()

	fmt.Printf("%v %v %v\n", keys, vals, entries)

	hash.Clear()
}
