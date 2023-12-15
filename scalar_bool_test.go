package kredis

func (s *KredisTestSuite) TestScalarBoolWithUnknownConnection() {
	config := "unknown"

	_, e := NewBool("foo", Options{Config: &config})

	s.Error(e)
}

func (s *KredisTestSuite) TestScalarBool() {
	k, e := NewBool("foo", Options{})

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
	k, e := NewBoolWithDefault("foo", Options{}, true)

	s.NoError(e)
	s.True(k.Value())

	k, e = NewBoolWithDefault("bar", Options{}, false)
	s.NoError(e)
	s.False(k.Value())

	k, e = NewBoolWithDefault("bar", Options{}, true)
	s.NoError(e)
	s.False(k.Value())
}

func (s *KredisTestSuite) TestScalarBoolBadConnection() {
	cfg := "badconn"

	k, e := NewBoolWithDefault("foo", Options{Config: &cfg}, true)
	s.Error(e)
	s.Nil(k)
}
