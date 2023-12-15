package kredis

import (
	"time"
)

func (s *KredisTestSuite) TestStringList() {
	elems := make([]string, 5)

	l, e := NewStringList("list", Options{})
	s.NoError(e)

	n, err := l.Elements(elems)
	s.NoError(err)
	s.Zero(n)

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

func (s *KredisTestSuite) TestBoolList() {
	elems := make([]bool, 5)

	l, e := NewBoolList("bool_list", Options{})
	s.NoError(e)

	n, err := l.Elements(elems)
	s.NoError(err)
	s.Zero(n)

	llen, err := l.Append(true, false, true)
	s.NoError(err)
	s.Equal(int64(3), llen)

	llen, err = l.Prepend(false, false)
	s.NoError(err)
	s.Equal(int64(5), llen)

	n, err = l.Elements(elems)
	s.NoError(err)
	s.Equal(5, n)
	s.Equal([]bool{false, false, true, false, true}, elems)

	elems = make([]bool, 2)
	n, err = l.Elements(elems)
	s.NoError(err)
	s.Equal(2, n)
	s.Equal([]bool{false, false}, elems)
}

func (s *KredisTestSuite) TestTimeList() {
	elems := make([]time.Time, 2)

	t1 := time.Now()
	t2 := time.Now()

	l, e := NewTimeList("list", Options{})
	s.NoError(e)

	n, err := l.Elements(elems)
	s.NoError(err)
	s.Zero(n)

	llen, err := l.Append(t1)
	s.NoError(err)
	s.Equal(int64(1), llen)

	llen, err = l.Prepend(t2)
	s.NoError(err)
	s.Equal(int64(2), llen)

	n, err = l.Elements(elems)
	s.NoError(err)
	s.Equal(2, n)

	s.Equal([]time.Time{t1.Round(0), t2.Round(0)}, elems)
}

func (s *KredisTestSuite) TestJSONList() {
	elems := make([]kredisJSON, 3)

	l, e := NewJSONList("json_list", Options{})
	s.NoError(e)

	n, err := l.Elements(elems)
	s.NoError(err)
	s.Zero(n)

	kj_1 := NewKredisJSON(`{"k1":"v1"}`)
	kj_2 := NewKredisJSON(`{"k2":"v2"}`)

	llen, err := l.Append(*kj_1, *kj_2)
	s.NoError(err)
	s.Equal(int64(2), llen)

	kj_3 := NewKredisJSON(`{"k3":"v3"}`)

	llen, err = l.Prepend(*kj_3)
	s.NoError(err)
	s.Equal(int64(3), llen)

	n, err = l.Elements(elems)
	s.NoError(err)
	s.Equal(3, n)

	s.Equal([]kredisJSON{*kj_3, *kj_1, *kj_2}, elems)

	var data interface{}

	err = elems[1].Unmarshal(&data)
	s.NoError(err)

	s.Equal(map[string]interface{}{"k1": "v1"}, data)
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
	s.Zero(n)
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
	s.Zero(n)

	llen, err := l.Append("x")

	s.Error(err)
	s.Empty(llen)
}
