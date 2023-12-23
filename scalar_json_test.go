package kredis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func (s *KredisTestSuite) TestNewJSON() {
	kjson := NewKredisJSON(`{"key":"value"}`)
	k, e := NewJSON("json", Options{})

	s.NoError(e)
	s.Empty(k.Value())

	e = k.SetValue(kjson)

	s.NoError(e)
	s.Equal(*kjson, k.Value())
}

func (s *KredisTestSuite) TestNewJSONWithDefaultValue() {
	kjson := NewKredisJSON(`{"key":"default"}`)
	kjson2 := NewKredisJSON(`{"k2":"v2"}`)
	k, e := NewJSONWithDefault("foo", Options{}, kjson)

	s.NoError(e)
	s.Equal(*kjson, k.Value())

	k, e = NewJSONWithDefault("bar", Options{}, kjson2)
	s.NoError(e)
	s.Equal(*kjson2, k.Value())

	k, e = NewJSONWithDefault("bar", Options{}, kjson)
	s.NoError(e)
	s.Equal(*kjson2, k.Value())
}

func TestScalarJSONUnmarshal(t *testing.T) {
	var data interface{}

	kj := NewKredisJSON("[1]")
	err := kj.Unmarshal(&data)

	nums := data.([]interface{})

	assert.NoError(t, err)
	assert.Contains(t, nums, 1.0)
}
