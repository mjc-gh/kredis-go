package kredis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *KredisTestSuite) TestNewTime() {
	k, e := NewTime("t")

	s.NoError(e)
	s.Empty(k.Value())

	now := time.Now()

	e = k.SetValue(now)
	s.NoError(e)

	// The canonical way to strip a monotonic clock reading is to use t =
	// t.Round(0).
	s.Equal(now.Round(0), k.Value().Local())
}

func (s *KredisTestSuite) TestNewTimeWithDefaultValue() {
	yesterday := time.Now().Add(-24 * time.Hour)
	now := time.Now()

	k, e := NewTimeWithDefault("t1", yesterday)

	s.NoError(e)
	s.Equal(yesterday.Round(0), k.Value().Local())

	k, e = NewTimeWithDefault("t2", now)
	s.NoError(e)
	s.Equal(now.Round(0), k.Value().Local())

	k, e = NewTimeWithDefault("t2", yesterday)
	s.NoError(e)
	s.Equal(now.Round(0), k.Value().Local())
}

func (s *KredisTestSuite) TestScalarTimeExpiresIn() {
	k, _ := NewTime("foo", WithExpiry("1ms"))

	e := k.SetValue(time.Now())
	s.NoError(e)

	time.Sleep(5 * time.Millisecond)

	v, e := k.ValueResult()
	s.Equal(redis.Nil, e)
	s.Nil(v)
}
