package main

import (
	"fmt"
	"time"

	"github.com/mjc-gh/kredis-go"
)

func main() {
	kredis.SetConfiguration("shared", "ns", "redis://localhost:6379/2")
	kredis.EnableDebugLogging()

	t := time.Date(2021, 8, 28, 23, 0, 0, 0, time.UTC)
	times := []time.Time{t, t.Add(1 * time.Hour), t.Add(2 * time.Hour), t.Add(3 * time.Hour)}

	set, _ := kredis.NewTimeSet("times")
	set.Add(times...)
	members, _ := set.Members()

	fmt.Println(set.Includes(t))
	fmt.Println(set.Includes(t.Add(4 * time.Hour)))
	fmt.Printf("%d %v\n", set.Size(), members)

	sample := make([]time.Time, 2)
	n, _ := set.Sample(sample)

	fmt.Printf("n = %d %v\n", n, sample)
}
