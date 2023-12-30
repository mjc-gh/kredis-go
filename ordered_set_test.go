package kredis

func (s *KredisTestSuite) TestOrderedSet() {
	oset, e := NewStringOrderedSet("ints", 3)
	s.NoError(e)

	add, rm, e := oset.Append("ready", "set", "go")
	s.NoError(e)
	s.Equal(int64(3), add)
	s.Equal(int64(0), rm)
	s.True(oset.Includes("go"))
	s.False(oset.Includes("not"))

	add, rm, e = oset.Prepend("not", "were")
	s.NoError(e)
	s.Equal(int64(2), add)
	s.Equal(int64(2), rm)

	members, e := oset.Members()
	s.NoError(e)
	s.Equal([]string{"were", "not", "ready"}, members)

	n, e := oset.Remove("not")
	s.NoError(e)
	s.Equal(int64(1), n)

	s.NoError(oset.Clear())
	s.Equal(int64(0), oset.Size())
}
