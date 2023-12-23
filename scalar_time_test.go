package kredis

import (
	"time"
)

func (s *KredisTestSuite) TestNewTime() {
	k, e := NewTime("t", Options{})

	s.NoError(e)
	s.Empty(k.Value())

	now := time.Now()

	e = k.SetValue(now)
	s.NoError(e)

	// The canonical way to strip a monotonic clock reading is to use t =
	// t.Round(0).
	s.Equal(now.Round(0), k.Value())
}

func (s *KredisTestSuite) TestNewTimeWithDefaultValue() {
	yesterday := time.Now().Add(-24 * time.Hour)
	now := time.Now()

	k, e := NewTimeWithDefault("t1", Options{}, yesterday)

	s.NoError(e)
	s.Equal(yesterday.Round(0), k.Value())

	k, e = NewTimeWithDefault("t2", Options{}, now)
	s.NoError(e)
	s.Equal(now.Round(0), k.Value())

	k, e = NewTimeWithDefault("t2", Options{}, yesterday)
	s.NoError(e)
	s.Equal(now.Round(0), k.Value())
}
