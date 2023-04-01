package kredis

func (s *KredisTestSuite) TestGetConfigurationUsesCacheMap() {
	_, ok := connections["shared"]
	s.False(ok)

	c, _, e := getConnection("shared")
	s.NoError(e)

	c2, _, e := getConnection("shared")
	s.Same(c, c2)
}
