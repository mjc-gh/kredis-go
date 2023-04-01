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
	s.Equal(now.UTC(), k.Value())
}

func (s *KredisTestSuite) TestNewTimeWithDefaultValue() {
	yesterday := time.Now().Add(-24 * time.Hour)

	k, e := NewTime("t", Options{DefaultValue: yesterday})

	s.NoError(e)
	s.Equal(yesterday, k.Value())
}
