package kredis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type KredisTestSuite struct {
	suite.Suite
}

func (suite *KredisTestSuite) SetupTest() {
	SetConfiguration("shared", "ns", "redis://localhost:6379/2")
	SetConfiguration("badconn", "ns", "redis://localhost:1234/0")
}

func (suite *KredisTestSuite) TearDownTest() {
	ctx := context.Background()
	c, _, _ := getConnection("shared")

	// Delete all keys in namespace to reset test state
	keys, _ := c.Keys(ctx, "ns:*").Result()

	for _, key := range keys {
		c.Del(ctx, key)
	}

	// Reset connections
	delete(connections, "shared")
	delete(connections, "badconn")
}

// listen for 'go test' command --> run test methods
func TestKredisTestSuit(t *testing.T) {
	suite.Run(t, new(KredisTestSuite))
}
