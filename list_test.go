package kredis

import "time"

func (s *KredisTestSuite) TestStringList() {
	elems := make([]string, 5)

	l, e := NewStringList("list", Options{})
	s.NoError(e)

	n, err := l.Elements(elems)
	s.NoError(err)
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

func (s *KredisTestSuite) TestTimeList() {
	elems := make([]time.Time, 1)

	t1 := time.Now()

	l, e := NewTimeList("list", Options{})
	s.NoError(e)

	n, err := l.Elements(elems)
	s.NoError(err)
	s.Equal(0, n)

	llen, err := l.Append(t1)
	s.NoError(err)
	s.Equal(int64(1), llen)

	n, err = l.Elements(elems)
	s.NoError(err)
	s.Equal([]time.Time{t1.UTC()}, elems)
}

func (s *KredisTestSuite) TestListClear() {
	elems := make([]string, 5)

	l, _ := NewStringList("list", Options{})
	llen, err := l.Append("a", "b", "c")
	s.NoError(err)
	s.Equal(int64(3), llen)

	err = l.Clear()
	s.NoError(err)

	n, _ := l.Elements(elems)
	s.Equal(0, n)
}

func (s *KredisTestSuite) TestListRemove() {
	elems := make([]string, 3)

	l, _ := NewStringList("list", Options{})
	llen, err := l.Append("a", "b", "c", "d", "e")
	s.NoError(err)
	s.Equal(int64(5), llen)

	err = l.Remove("b", "d")
	s.NoError(err)

	n, _ := l.Elements(elems)
	s.Equal(3, n)
	s.Equal([]string{"a", "c", "e"}, elems)
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
