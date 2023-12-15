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
	k, e := NewJSON("foo", Options{DefaultValue: *kjson})

	s.NoError(e)
	s.Equal(*kjson, k.Value())
}

func TestScalarJSONUnmarshal(t *testing.T) {
	var data interface{}

	kj := NewKredisJSON("[1]")
	err := kj.Unmarshal(&data)

	nums := data.([]interface{})

	assert.NoError(t, err)
	assert.Contains(t, nums, 1.0)
}
