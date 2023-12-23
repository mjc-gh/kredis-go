package kredis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *KredisTestSuite) TestScalarBoolWithUnknownConnection() {
	_, e := NewBool("foo", WithConfigName("unknown"))

	s.Error(e)
}

func (s *KredisTestSuite) TestScalarBool() {
	k, e := NewBool("foo")

	s.NoError(e)
	s.Empty(k.Value())

	e = k.SetValue(false)

	s.NoError(e)
	s.False(k.Value())

	e = k.SetValue(true)

	s.NoError(e)
	s.True(k.Value())
}

func (s *KredisTestSuite) TestScalarBoolWithDefaultValue() {
	k, e := NewBoolWithDefault("foo", true)

	s.NoError(e)
	s.True(k.Value())

	k, e = NewBoolWithDefault("bar", false)
	s.NoError(e)
	s.False(k.Value())

	k, e = NewBoolWithDefault("bar", true)
	s.NoError(e)
	s.False(k.Value())
}

func (s *KredisTestSuite) TestScalarBoolBadConnection() {
	k, _ := NewBool("foo", WithConfigName("badconn"))
	_, e := k.ValueResult()

	s.Error(e)
	s.Empty(k.Value())

	k, e = NewBoolWithDefault("foo", true, WithConfigName("badconn"))
	s.Error(e)
	s.Nil(k)
}

func (s *KredisTestSuite) TestScalarBoolExpiresIn() {
	k, _ := NewBool("foo", WithExpiry("1ms"))

	e := k.SetValue(true)
	s.NoError(e)

	time.Sleep(5 * time.Millisecond)

	v, e := k.ValueResult()
	s.Equal(redis.Nil, e)
	s.Nil(v)
}
