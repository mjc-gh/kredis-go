package kredis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionsGetConfigDefault(t *testing.T) {
	o := Options{}

	assert.Equal(t, "shared", o.GetConfig())
}
