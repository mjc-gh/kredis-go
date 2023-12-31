package kredis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProxyOptionsConfigDefault(t *testing.T) {
	po := ProxyOptions{}

	assert.Equal(t, "shared", po.Config())
}

func TestProxyOptionsWithConfigName(t *testing.T) {
	po := ProxyOptions{}
	WithConfigName("named")(&po)

	assert.Equal(t, "named", po.Config())
}
