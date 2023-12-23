package kredis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *KredisTestSuite) TestScalarIntegerWithUnknownConnection() {
	_, e := NewInteger("foo", WithConfigName("unknown"))

	s.Error(e)
}

func (s *KredisTestSuite) TestScalarInteger() {
	k, e := NewInteger("foo")

	s.NoError(e)
	s.Empty(k.Value())

	e = k.SetValue(1234)

	s.NoError(e)
	s.Equal(1234, k.Value())
}

func (s *KredisTestSuite) TestScalarIntegerWithDefaultValue() {
	k, e := NewIntegerWithDefault("foo", 5678)

	s.NoError(e)
	s.Equal(5678, k.Value())
}

func (s *KredisTestSuite) TestScalarIntegerBadConnection() {
	k, _ := NewInteger("foo", WithConfigName("badconn"))
	_, e := k.ValueResult()

	s.Error(e)
	s.Empty(k.Value())

	k, e = NewIntegerWithDefault("foo", -1, WithConfigName("badconn"))
	s.Error(e)
	s.Nil(k)
}

func (s *KredisTestSuite) TestScalarIntegerExpiresIn() {
	k, _ := NewInteger("foo", WithExpiry("1ms"))

	e := k.SetValue(828)
	s.NoError(e)

	time.Sleep(5 * time.Millisecond)

	v, e := k.ValueResult()
	s.Equal(redis.Nil, e)
	s.Nil(v)
}
