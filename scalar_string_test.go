package kredis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *KredisTestSuite) TestNewString() {
	k, e := NewString("foo")

	s.NoError(e)
	s.Empty(k.Value())

	e = k.SetValue("bar")

	s.NoError(e)
	s.Equal("bar", k.Value())
}

func (s *KredisTestSuite) TestNewStringWithDefaultValue() {
	k, e := NewStringWithDefault("foo", "bar")

	s.NoError(e)
	s.Equal("bar", k.Value())

	k, e = NewStringWithDefault("bar", "baz")
	s.NoError(e)
	s.Equal("baz", k.Value())

	k, e = NewStringWithDefault("bar", "quix")
	s.NoError(e)
	s.Equal("baz", k.Value())
}

func (s *KredisTestSuite) TestScalarStringExpiresIn() {
	k, _ := NewString("foo", WithExpiry("1ms"))

	e := k.SetValue("bar")
	s.NoError(e)

	time.Sleep(5 * time.Millisecond)

	v, e := k.ValueResult()
	s.Equal(redis.Nil, e)
	s.Nil(v)
}
