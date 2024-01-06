package kredis

func (s *KredisTestSuite) TestEnum() {
	enum, err := NewEnum("e", "ready", []string{"ready", "set", "go"})
	s.NoError(err)

	s.Equal("ready", enum.Value())
	s.True(enum.Is("ready"))

	s.NoError(enum.SetValue("set"))
	s.Equal("set", enum.Value())
	s.True(enum.Is("set"))

	s.NoError(enum.SetValue("go"))
	s.Equal("go", enum.Value())
	s.True(enum.Is("go"))

	err = enum.SetValue("badval")
	s.Error(err)
	s.Equal(EnumInvalidValue, err)

	s.NoError(enum.Reset())
	s.Equal("ready", enum.Value())
	s.True(enum.Is("ready"))
}

func (s *KredisTestSuite) TestEnumWithEmptyValues() {
	_, err := NewEnum("key", "ready", []string{})
	s.Error(err)
	s.Equal(EnumEmptyValues, err)
}

func (s *KredisTestSuite) TestEnumWithExpiry() {
	_, err := NewEnum("key", "go", []string{"go"}, WithExpiry("1ms"))
	s.Error(err)
	s.Equal(EnumExpiryNotSupported, err)
}
