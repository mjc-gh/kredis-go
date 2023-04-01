package kredis

import (
	"context"
	"fmt"
)

func setup() {
	SetConfiguration("shared", "ns", "redis://localhost:6379/2")
}

func teardown() {
	ctx := context.Background()
	c, _, _ := getConnection("shared")

	// TODO insert namespace here when impl
	keys, _ := c.Keys(ctx, "*").Result()

	fmt.Println("teardown")

	for _, key := range keys {
		c.Del(ctx, key)
	}
}
