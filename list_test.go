package kredis

func (s *KredisTestSuite) TestStringList() {
	elems := make([]string, 5)

	l, e := NewStringList("list", Options{})
	s.NoError(e)

	n, err := l.Elements(elems)
	s.NoError(e)
	s.Equal(0, n)

	llen, err := l.Append("a", "b", "c")
	s.NoError(err)
	s.Equal(int64(3), llen)

	llen, err = l.Prepend("x", "y")
	s.NoError(err)
	s.Equal(int64(5), llen)

	n, err = l.Elements(elems)
	s.NoError(err)
	s.Equal(5, n)
	s.Equal([]string{"y", "x", "a", "b", "c"}, elems)

	elems = make([]string, 3)
	n, err = l.Elements(elems)
	s.NoError(err)
	s.Equal(3, n)
	s.Equal([]string{"y", "x", "a"}, elems)
}

func (s *KredisTestSuite) TestListBadConnection() {
	var elems = []string{}

	cfg := "badconn"

	l, _ := NewStringList("list", Options{Config: &cfg})
	n, err := l.Elements(elems)

	s.Error(err)
	s.Equal(0, n)

	llen, err := l.Append("x")

	s.Error(err)
	s.Empty(llen)
}
