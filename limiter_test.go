package kredis

func (s *KredisTestSuite) TestLimiter() {
	limiter, _ := NewLimiter("limiter", 3)

	s.False(limiter.IsExceeded())

	limiter.Poke()
	limiter.Poke()
	limiter.Poke()

	s.True(limiter.IsExceeded())
}

func (s *KredisTestSuite) TestLimiterWithBadConnection() {
	limiter, _ := NewLimiter("limiter", 3, WithConfigName("badconn"))
	limiter.Poke()

	s.False(limiter.IsExceeded())

	limiter, _ = NewLimiterWithDefault("limiter", 3, 2, WithConfigName("badconn"))
	limiter.Poke()

	s.False(limiter.IsExceeded())
}
