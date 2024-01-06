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
	s.Equal(int64(5), n)
	s.Equal([]string{"y", "x", "a", "b", "c"}, elems)

	elems = make([]string, 3)
	n, err = l.Elements(elems)
	s.NoError(err)
	s.Equal(int64(3), n)
	s.Equal([]string{"y", "x", "a"}, elems)

	last, ok := l.Last()
	s.True(ok)
	s.Equal("c", last)

	last2 := make([]string, 2)
	n, e = l.LastN(last2)
	s.NoError(e)
	s.Equal(int64(2), n)
	s.Equal([]string{"b", "c"}, last2)

	s.NoError(l.Remove("x", "a"))

	n, err = l.Elements(elems)
	s.NoError(err)
	s.Equal(int64(3), n)
	s.Equal([]string{"y", "b", "c"}, elems)
}

func (s *KredisTestSuite) TestStringListWithDefault() {
	l, _ := NewStringListWithDefault("list_default", []string{"a"})

	NewStringListWithDefault("list_default", []string{"x", "y", "z"})

	elems := make([]string, 3)
	n, e := l.Elements(elems)
	s.NoError(e)
	s.Equal(int64(1), n)
	s.Equal([]string{"a"}, elems[0:1])
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
	s.Equal(int64(5), n)
	s.Equal([]int{9, 8, 1, 2, 3}, elems)

	elems = make([]int, 3)
	n, err = l.Elements(elems)
	s.NoError(err)
	s.Equal(int64(3), n)
	s.Equal([]int{9, 8, 1}, elems)

	s.NoError(l.Clear())

	last, ok := l.Last()
	s.False(ok)
	s.Empty(last)
}

func (s *KredisTestSuite) TestIntegerListWithDefault() {
	l, _ := NewIntegerListWithDefault("list_default", []int{5, 7, 9})

	NewIntegerListWithDefault("list_default", []int{1, 2, 3, 4, 5})

	elems := make([]int, 3)
	n, e := l.Elements(elems)
	s.NoError(e)
	s.Equal(int64(3), n)
	s.Equal([]int{5, 7, 9}, elems)
}

func (s *KredisTestSuite) TestFloatList() {
	elems := make([]float64, 5)

	l, e := NewFloatList("list")
	s.NoError(e)

	n, err := l.Elements(elems)
	s.NoError(err)
	s.Zero(n)

	llen, err := l.Append(1.1, 2.4, 3.17)
	s.NoError(err)
	s.Equal(int64(3), llen)

	llen, err = l.Prepend(8.001, 9.000001)
	s.NoError(err)
	s.Equal(int64(5), llen)

	n, err = l.Elements(elems)
	s.NoError(err)
	s.Equal(int64(5), n)
	s.Equal([]float64{9.000001, 8.001, 1.1, 2.4, 3.17}, elems)

	elems = make([]float64, 3)
	n, err = l.Elements(elems)
	s.NoError(err)
	s.Equal(int64(3), n)
	s.Equal([]float64{9.000001, 8.001, 1.1}, elems)

	s.NoError(l.Clear())
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
	s.Equal(int64(5), n)
	s.Equal([]bool{false, false, true, false, true}, elems)

	elems = make([]bool, 2)
	n, err = l.Elements(elems)
	s.NoError(err)
	s.Equal(int64(2), n)
	s.Equal([]bool{false, false}, elems)
}

func (s *KredisTestSuite) TestBoolListWithDefault() {
	l, _ := NewBoolListWithDefault("bool_list_default", []bool{true, false, true})

	NewBoolListWithDefault("bool_list_default", []bool{false, false, false, false})

	elems := make([]bool, 3)
	n, e := l.Elements(elems)
	s.NoError(e)
	s.Equal(int64(3), n)
	s.Equal([]bool{true, false, true}, elems)
}

// TODO fix flaky tests with time values
func (s *KredisTestSuite) TestTimeList() {
	elems := make([]time.Time, 2)

	t1 := time.Now()
	t2 := time.Date(2021, 8, 28, 23, 0, 0, 0, time.Local)

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
	s.Equal(int64(2), n)

	s.Equal(t2.Round(0), elems[0].Local())
	s.Equal(t1.Round(0), elems[1].Local())

	s.NoError(l.Remove(t2))

	n, e = l.Length()
	s.NoError(e)
	s.Equal(int64(1), n)

	e = l.Clear()
	s.NoError(e)

	n, _ = l.Elements(elems)
	s.Zero(n)
}

func (s *KredisTestSuite) TestTimeListWithDefault() {
	t1 := time.Date(2021, time.Month(2), 21, 1, 10, 30, 0, time.UTC)
	t2 := time.Date(2022, time.Month(2), 21, 12, 10, 30, 0, time.UTC)

	l, _ := NewTimeListWithDefault("list_default", []time.Time{t1, t2})

	NewTimeListWithDefault("list_default", []time.Time{time.Now()})

	elems := make([]time.Time, 3)
	n, e := l.Elements(elems)
	s.NoError(e)
	s.Equal(int64(2), n)
	s.Equal([]time.Time{t1, t2}, elems[0:2])
}

func (s *KredisTestSuite) TestJSONList() {
	elems := make([]KredisJSON, 3)

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
	s.Equal(int64(3), n)

	s.Equal([]KredisJSON{*kj_3, *kj_1, *kj_2}, elems)

	var data interface{}

	err = elems[1].Unmarshal(&data)
	s.NoError(err)

	s.Equal(map[string]interface{}{"k1": "v1"}, data)
}

func (s *KredisTestSuite) TestJSONListWithDefault() {
	kj_1 := NewKredisJSON(`{"k1":"v1"}`)
	kj_2 := NewKredisJSON(`{"k2":"v2"}`)

	l, _ := NewJSONListWithDefault("list_default", []KredisJSON{*kj_1, *kj_2})

	NewJSONListWithDefault("list_default", []KredisJSON{*NewKredisJSON(`{"abc":"xyz"}`)})

	elems := make([]KredisJSON, 3)
	n, e := l.Elements(elems)
	s.NoError(e)
	s.Equal(int64(2), n)
	s.Equal([]KredisJSON{*kj_1, *kj_2}, elems[0:2])
}

func (s *KredisTestSuite) TestElementsWithRangeOptions() {
	l, _ := NewIntegerList("list")
	_, e := l.Append(1, 2, 3, 4, 5, 6)
	s.NoError(e)

	elems := make([]int, 6)
	n, e := l.Elements(elems)
	s.NoError(e)
	s.Equal(int64(6), n)
	s.Equal([]int{1, 2, 3, 4, 5, 6}, elems)

	elems = make([]int, 2)
	n, e = l.Elements(elems, WithRangeStart(3))
	s.NoError(e)
	s.Equal(int64(2), n)
	s.Equal([]int{4, 5}, elems)
}

func (s *KredisTestSuite) TestListLength() {
	l, _ := NewStringList("list")
	n, e := l.Length()
	s.NoError(e)
	s.Equal(int64(0), n)

	_, e = l.Append("a", "b", "c", "d", "e")
	s.NoError(e)

	n, e = l.Length()
	s.NoError(e)
	s.Equal(int64(5), n)
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
