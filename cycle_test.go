package kredis

func (s *KredisTestSuite) TestCycle() {
	cycle, err := NewCycle("cycle", []string{"ready", "set", "go"})
	s.NoError(err)
	s.Equal("ready", cycle.Value())

	s.NoError(cycle.Next())
	s.Equal("set", cycle.Value())

	s.NoError(cycle.Next())
	s.Equal("go", cycle.Value())

	s.NoError(cycle.Next())
	s.Equal("ready", cycle.Value())
	s.Empty(cycle.Index())
}
