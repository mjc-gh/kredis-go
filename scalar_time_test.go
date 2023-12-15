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

	k, e := NewTime("t", Options{DefaultValue: yesterday})

	s.NoError(e)
	s.Equal(yesterday, k.Value())
}
