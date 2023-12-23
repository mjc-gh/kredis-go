package kredis

import (
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func (s *KredisTestSuite) TestNewJSON() {
	kjson := NewKredisJSON(`{"key":"value"}`)
	k, e := NewJSON("json")

	s.NoError(e)
	s.Empty(k.Value())

	e = k.SetValue(kjson)

	s.NoError(e)
	s.Equal(*kjson, k.Value())
}

func (s *KredisTestSuite) TestNewJSONWithDefaultValue() {
	kjson := NewKredisJSON(`{"key":"default"}`)
	kjson2 := NewKredisJSON(`{"k2":"v2"}`)
	k, e := NewJSONWithDefault("foo", kjson)

	s.NoError(e)
	s.Equal(*kjson, k.Value())

	k, e = NewJSONWithDefault("bar", kjson2)
	s.NoError(e)
	s.Equal(*kjson2, k.Value())

	k, e = NewJSONWithDefault("bar", kjson)
	s.NoError(e)
	s.Equal(*kjson2, k.Value())
}

func (s *KredisTestSuite) TestScalarJSONExpiresIn() {
	k, _ := NewJSON("foo", WithExpiry("1ms"))

	e := k.SetValue(NewKredisJSON(`[1]`))
	s.NoError(e)

	time.Sleep(5 * time.Millisecond)

	v, e := k.ValueResult()
	s.Equal(redis.Nil, e)
	s.Nil(v)
}

func TestScalarJSONUnmarshal(t *testing.T) {
	var data interface{}

	kj := NewKredisJSON("[1]")
	err := kj.Unmarshal(&data)

	nums := data.([]interface{})

	assert.NoError(t, err)
	assert.Contains(t, nums, 1.0)
}
