package kredis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIntegerWithUnknownConnection(t *testing.T) {
	setup()

	config := "unknown"

	_, e := NewInteger("foo", Options{Config: &config})

	assert.Error(t, e)
}

func TestNewInteger(t *testing.T) {
	setup()
	defer teardown()

	k, e := NewInteger("foo", Options{})

	assert.NoError(t, e)
	assert.Empty(t, k.Value())

	e = k.SetValue(1234)

	assert.NoError(t, e)
	assert.Equal(t, 1234, k.Value())
}

func TestNewIntegerWithDefaultValue(t *testing.T) {
	setup()
	defer teardown()

	k, e := NewInteger("foo", Options{DefaultValue: 5678})

	assert.NoError(t, e)
	assert.Equal(t, 5678, k.Value())
}

func TestNewString(t *testing.T) {
	setup()
	defer teardown()

	k, e := NewString("foo", Options{})

	assert.NoError(t, e)
	assert.Empty(t, k.Value())

	e = k.SetValue("bar")

	assert.NoError(t, e)
	assert.Equal(t, "bar", k.Value())
}
