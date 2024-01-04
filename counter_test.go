package kredis

import "time"

func (s *KredisTestSuite) TestCounter() {
	c, err := NewCounter("counter")
	s.NoError(err)

	v, err := c.Increment(1)
	s.NoError(err)
	s.Equal(int64(1), v)

	v, err = c.Increment(2)
	s.NoError(err)
	s.Equal(int64(3), v)

	v, err = c.Decrement(1)
	s.NoError(err)
	s.Equal(int64(2), v)

	v, err = c.Decrement(1)
	s.NoError(err)
	s.Equal(int64(1), v)

	v, err = c.Decrement(2)
	s.NoError(err)
	s.Equal(int64(-1), v)

	s.NoError(c.Reset())
	s.Empty(c.Value())
}

func (s *KredisTestSuite) TestCounterWithExpiry() {
	c, _ := NewCounter("counter", WithExpiry("1ms"))
	c.Increment(1)
	c.Increment(2)
	c.Increment(3)

	time.Sleep(5 * time.Millisecond)

	s.Empty(c.Value())
}

func (s *KredisTestSuite) TestCounterWithDefault() {
	c, _ := NewCounterWithDefault("counter", 23)

	s.Equal(int64(23), c.Value())
}

func (s *KredisTestSuite) TestCounterWithDefaultAndExpiry() {
	c, _ := NewCounterWithDefault("counter", 23, WithExpiry("1ms"))

	c.Increment(1)
	c.Increment(2)
	c.Increment(3)

	time.Sleep(5 * time.Millisecond)

	s.Empty(c.Value())
}
