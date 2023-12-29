package kredis

func (s *KredisTestSuite) TestUniqueList() {
	l, e := NewIntegerUniqueList("uniq_ints", 3)
	s.NoError(e)

	_, e = l.Append(1, 2)
	s.NoError(e)

	_, e = l.Prepend(2, 4, 6)
	s.NoError(e)

	n, e := l.Length()
	s.NoError(e)
	s.Equal(int64(3), n)

	elements := make([]int, 3)
	n, e = l.Elements(elements)
	s.NoError(e)
	s.Equal(int64(3), n)
	s.Equal([]int{4, 2, 1}, elements)
}
