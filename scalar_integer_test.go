package kredis

func (s *KredisTestSuite) TestNewIntegerWithUnknownConnection() {
	config := "unknown"

	_, e := NewInteger("foo", Options{Config: &config})

	s.Error(e)
}

func (s *KredisTestSuite) TestNewInteger() {
	k, e := NewInteger("foo", Options{})

	s.NoError(e)
	s.Empty(k.Value())

	e = k.SetValue(1234)

	s.NoError(e)
	s.Equal(1234, k.Value())
}

func (s *KredisTestSuite) TestNewIntegerWithDefaultValue() {
	k, e := NewInteger("foo", Options{DefaultValue: 5678})

	s.NoError(e)
	s.Equal(5678, k.Value())
}

func (s *KredisTestSuite) TestBadConnection() {
	cfg := "badconn"

	k, _ := NewInteger("foo", Options{Config: &cfg})
	_, e := k.ValueWithErr()

	s.Error(e)
	s.Empty(k.Value())
}
