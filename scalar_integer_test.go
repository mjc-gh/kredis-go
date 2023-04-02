package kredis

func (s *KredisTestSuite) TestScalarIntegerWithUnknownConnection() {
	config := "unknown"

	_, e := NewInteger("foo", Options{Config: &config})

	s.Error(e)
}

func (s *KredisTestSuite) TestScalarInteger() {
	k, e := NewInteger("foo", Options{})

	s.NoError(e)
	s.Empty(k.Value())

	e = k.SetValue(1234)

	s.NoError(e)
	s.Equal(1234, k.Value())
}

func (s *KredisTestSuite) TestScalarIntegerWithDefaultValue() {
	k, e := NewInteger("foo", Options{DefaultValue: 5678})

	s.NoError(e)
	s.Equal(5678, k.Value())
}

func (s *KredisTestSuite) TestScalarIntegerBadConnection() {
	cfg := "badconn"

	k, _ := NewInteger("foo", Options{Config: &cfg})
	_, e := k.ValueResult()

	s.Error(e)
	s.Empty(k.Value())
}
