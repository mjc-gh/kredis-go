package kredis

func (s *KredisTestSuite) TestIntegerUniqueList() {
	l, e := NewIntegerUniqueList("uniq_ints", 3)
	s.NoError(e)

	n, e := l.Append(1, 2)
	s.NoError(e)
	s.Equal(int64(2), n)

	n, e = l.Prepend(2, 4, 6)
	s.NoError(e)
	s.Equal(int64(4), n)

	n, e = l.Length()
	s.NoError(e)
	s.Equal(int64(3), n)

	elements := make([]int, 3)
	n, e = l.Elements(elements)
	s.NoError(e)
	s.Equal(int64(3), n)
	s.Equal([]int{4, 2, 1}, elements)

	s.NoError(l.Remove(2, 1))

	elements = make([]int, 1)
	n, e = l.Elements(elements)
	s.NoError(e)
	s.Equal(int64(1), n)
	s.Equal([]int{4}, elements)
}

func (s *KredisTestSuite) TestFloatUniqueList() {
	l, e := NewFloatUniqueList("uniq_float64s", 3)
	s.NoError(e)

	n, e := l.Append(1.2, 2.3)
	s.NoError(e)
	s.Equal(int64(2), n)

	n, e = l.Prepend(2.3, 4.5, 6.7)
	s.NoError(e)
	s.Equal(int64(4), n)

	n, e = l.Length()
	s.NoError(e)
	s.Equal(int64(3), n)

	elements := make([]float64, 3)
	n, e = l.Elements(elements)
	s.NoError(e)
	s.Equal(int64(3), n)
	s.Equal([]float64{4.5, 2.3, 1.2}, elements)

	s.NoError(l.Remove(2, 1))

	elements = make([]float64, 1)
	n, e = l.Elements(elements)
	s.NoError(e)
	s.Equal(int64(1), n)
	s.Equal([]float64{4.5}, elements)
}

func (s *KredisTestSuite) TestStringUniqueList() {
	l, e := NewStringUniqueList("uniq_strs", 5)
	s.NoError(e)

	elements := make([]string, 5)
	n, e := l.Elements(elements)
	s.NoError(e)
	s.Equal(int64(0), n)

	n, e = l.Append("a", "b")
	s.NoError(e)
	s.Equal(int64(2), n)

	n, e = l.Prepend("x", "y", "z")
	s.NoError(e)
	s.Equal(int64(5), n)

	n, e = l.Length()
	s.NoError(e)
	s.Equal(int64(5), n)

	elements = make([]string, 5)
	n, e = l.Elements(elements)
	s.NoError(e)
	s.Equal(int64(5), n)
	s.Equal([]string{"z", "y", "x", "a", "b"}, elements)

	s.NoError(l.Remove("y", "z"))

	elements = make([]string, 3)
	n, e = l.Elements(elements)
	s.NoError(e)
	s.Equal(int64(3), n)
	s.Equal([]string{"x", "a", "b"}, elements)

	last, ok := l.Last()
	s.True(ok)
	s.Equal("b", last)

	last2 := make([]string, 2)
	n, e = l.LastN(last2)
	s.NoError(e)
	s.Equal(int64(2), n)
	s.Equal([]string{"a", "b"}, last2)
}
