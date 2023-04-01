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
	k, e := NewString("foo", Options{DefaultValue: "bar"})

	s.NoError(e)
	s.Equal("bar", k.Value())
}
