package kredis

import (
	"context"
	"fmt"
	"testing"

	redistesthooks "github.com/mjc-gh/redis-test-hook"
	"github.com/stretchr/testify/suite"
)

type KredisTestSuite struct {
	suite.Suite
	captureHook *redistesthooks.Hook
}

type testLogger struct{ stdLogger }

var testWarnings []string

func (tl testLogger) Warn(fnName string, err error) {
	testWarnings = append(testWarnings, fmt.Sprintf("%s %s", fnName, err.Error()))
}

func (suite *KredisTestSuite) SetupTest() {
	// TODO use a unique namespace for each test (thus potentially enabling
	// parallel tests)
	SetConfiguration("shared", "ns", "redis://localhost:6379/2")
	SetConfiguration("capture", "ns", "redis://localhost:6379/2")
	SetConfiguration("badconn", "", "redis://localhost:1234/0")

	EnableDebugLogging()

	suite.captureHook = redistesthooks.New()

	client, _, _ := getConnection("capture")
	client.AddHook(suite.captureHook)

	testWarnings = []string{}
	SetDebugLogger(&testLogger{})
}

func (suite *KredisTestSuite) TearDownTest() {
	ctx := context.Background()
	c, _, _ := getConnection("shared")

	// Delete all keys in namespace to reset test state
	keys, _ := c.Keys(ctx, "ns:*").Result()

	for _, key := range keys {
		c.Del(ctx, key)
	}

	suite.captureHook.Reset()

	// Reset connections
	delete(connections, "shared")
	delete(connections, "badconn")
	delete(connections, "capture")
}

// listen for 'go test' command --> run test methods
func TestKredisTestSuit(t *testing.T) {
	suite.Run(t, new(KredisTestSuite))
}

func (s *KredisTestSuite) TestKredisJSON() {
	kj := NewKredisJSON(`{"a":1}`)

	s.Equal(`{"a":1}`, kj.String())

	var data interface{}
	err := kj.Unmarshal(&data)
	s.NoError(err)

	obj, ok := data.(map[string]interface{})
	s.True(ok)
	s.Equal(1.0, obj["a"])
}
