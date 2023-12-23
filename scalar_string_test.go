package kredis

func (s *KredisTestSuite) TestNewString() {
	k, e := NewString("foo", Options{})

	s.NoError(e)
	s.Empty(k.Value())

	e = k.SetValue("bar")

	s.NoError(e)
	s.Equal("bar", k.Value())
}

func (s *KredisTestSuite) TestNewStringWithDefaultValue() {
	k, e := NewStringWithDefault("foo", Options{}, "bar")

	s.NoError(e)
	s.Equal("bar", k.Value())

	k, e = NewStringWithDefault("bar", Options{}, "baz")
	s.NoError(e)
	s.Equal("baz", k.Value())

	k, e = NewStringWithDefault("bar", Options{}, "quix")
	s.NoError(e)
	s.Equal("baz", k.Value())
}
