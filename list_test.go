package kredis

import (
	"time"
)

func (s *KredisTestSuite) TestStringList() {
	elems := make([]string, 5)

	l, e := NewStringList("list")
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

func (s *KredisTestSuite) TestIntegerList() {
	elems := make([]int, 5)

	l, e := NewIntegerList("list")
	s.NoError(e)

	n, err := l.Elements(elems)
	s.NoError(err)
	s.Zero(n)

	llen, err := l.Append(1, 2, 3)
	s.NoError(err)
	s.Equal(int64(3), llen)

	llen, err = l.Prepend(8, 9)
	s.NoError(err)
	s.Equal(int64(5), llen)

	n, err = l.Elements(elems)
	s.NoError(err)
	s.Equal(5, n)
	s.Equal([]int{9, 8, 1, 2, 3}, elems)

	elems = make([]int, 3)
	n, err = l.Elements(elems)
	s.NoError(err)
	s.Equal(3, n)
	s.Equal([]int{9, 8, 1}, elems)
}

func (s *KredisTestSuite) TestBoolList() {
	elems := make([]bool, 5)

	l, e := NewBoolList("bool_list")
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

func (s *KredisTestSuite) TestBoolListWithDefault() {
	elems := make([]bool, 3)

	l, e := NewBoolListWithDefault("bool_list_default", []bool{true, false, true})
	s.NoError(e)

	n, e := l.Elements(elems)
	s.NoError(e)
	s.Equal(3, n)
	s.Equal([]bool{true, false, true}, elems)

	l2, e := NewBoolListWithDefault("bool_list_default", []bool{false, false, false, false})
	s.NoError(e)

	elems = make([]bool, 3)
	n, e = l2.Elements(elems)
	s.NoError(e)
	s.Equal(3, n)
	s.Equal([]bool{true, false, true}, elems)
}

// TODO fix flaky tests with time values
func (s *KredisTestSuite) TestTimeList() {
	elems := make([]time.Time, 2)

	t1 := time.Now()
	t2 := time.Now()

	l, e := NewTimeList("list")
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

	l, e := NewJSONList("json_list")
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

func (s *KredisTestSuite) TestElementsWithOptions() {
	l, _ := NewIntegerList("list")
	_, e := l.Append(1, 2, 3, 4, 5, 6)
	s.NoError(e)

	elems := make([]int, 6)
	n, e := l.Elements(elems)
	s.NoError(e)
	s.Equal(6, n)
	s.Equal([]int{1, 2, 3, 4, 5, 6}, elems)

	elems = make([]int, 2)
	n, e = l.Elements(elems, WithRangeStart(3))
	s.NoError(e)
	s.Equal(2, n)
	s.Equal([]int{4, 5}, elems)
}

func (s *KredisTestSuite) TestListClear() {
	elems := make([]string, 5)

	l, _ := NewStringList("list")
	llen, err := l.Append("a", "b", "c")
	s.NoError(err)
	s.Equal(int64(3), llen)

	err = l.Clear()
	s.NoError(err)

	n, _ := l.Elements(elems)
	s.Zero(n)
}

func (s *KredisTestSuite) TestListAppendAndPrependWithEmptyElements() {
	l, _ := NewStringList("list")

	n, e := l.Append()
	s.NoError(e)
	s.Empty(n)

	n, e = l.Prepend()
	s.NoError(e)
	s.Empty(n)
}

func (s *KredisTestSuite) TestListRemove() {
	elems := make([]string, 3)

	l, _ := NewStringList("list")
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
	elems := make([]string, 1)

	l, _ := NewStringList("list", WithConfigName("badconn"))
	n, err := l.Elements(elems)

	s.Error(err)
	s.Zero(n)

	llen, err := l.Append("x")

	s.Error(err)
	s.Empty(llen)
}
