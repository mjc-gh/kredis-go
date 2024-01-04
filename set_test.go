package kredis

import "time"

func (s *KredisTestSuite) TestBoolSet() {
	set, e := NewBoolSet("bools")
	s.NoError(e)

	members, e := set.Members()
	s.NoError(e)
	s.Empty(members)

	n, e := set.Add(true, false)
	s.NoError(e)
	s.Equal(int64(2), n)

	members, e = set.Members()
	s.NoError(e)
	s.Contains(members, true)
	s.Contains(members, false)

	n, e = set.Remove(false)
	s.NoError(e)
	s.Equal(int64(1), n)

	members, e = set.Members()
	s.NoError(e)
	s.Len(members, 1)

	members = make([]bool, 1)
	m, e := set.Sample(members)
	s.NoError(e)
	s.Equal(int64(1), m)
	s.Equal([]bool{true}, members)

	t, ok := set.Take()
	s.True(ok)
	s.True(t)

	t, ok = set.Take()
	s.False(ok)
}

func (s *KredisTestSuite) TestIntegerSet() {
	set, e := NewIntegerSet("ints")
	s.NoError(e)

	n, e := set.Add(1, 5, 2)
	s.NoError(e)
	s.Equal(int64(3), n)

	n, e = set.Replace(4, 3)
	s.NoError(e)
	s.Equal(int64(2), n)

	s.True(set.Includes(4))
	s.False(set.Includes(5))
	s.Equal(int64(2), set.Size())

	s.NoError(set.Clear())
	s.False(set.Includes(4))
	s.Equal(int64(0), set.Size())

	set.Add(1)
	t, ok := set.Take()
	s.True(ok)
	s.Equal(int(1), t)

	t, ok = set.Take()
	s.False(ok)
}

func (s *KredisTestSuite) TestFloatSet() {
	set, e := NewFloatSet("floats")
	s.NoError(e)

	n, e := set.Add(1.1, 5.2, 2.5)
	s.NoError(e)
	s.Equal(int64(3), n)

	n, e = set.Replace(4.4, 3.7)
	s.NoError(e)
	s.Equal(int64(2), n)

	s.True(set.Includes(4.4))
	s.False(set.Includes(5.2))
	s.Equal(int64(2), set.Size())

	s.NoError(set.Clear())
	s.False(set.Includes(4))
	s.Equal(int64(0), set.Size())

	set.Add(1.1)
	t, ok := set.Take()
	s.True(ok)
	s.Equal(1.1, t)

	t, ok = set.Take()
	s.False(ok)
}

func (s *KredisTestSuite) TestStringSet() {
	set, e := NewStringSet("strings")
	s.NoError(e)

	n, e := set.Add("a", "b", "c")
	s.NoError(e)
	s.Equal(int64(3), n)

	members := make([]string, 5)
	m, e := set.Sample(members)
	s.NoError(e)
	s.Equal(int64(3), m)
	s.Contains(members, "a")
	s.Contains(members, "b")
	s.Contains(members, "c")

	n, e = set.Replace("d", "e")
	s.NoError(e)
	s.Equal(int64(2), n)

	s.True(set.Includes("d"))
	s.False(set.Includes("a"))
	s.Equal(int64(2), set.Size())

	s.NoError(set.Clear())
	s.False(set.Includes("d"))
	s.Equal(int64(0), set.Size())

	set.Add("x")
	t, ok := set.Take()
	s.True(ok)
	s.Equal("x", t)

	t, ok = set.Take()
	s.False(ok)
}

func (s *KredisTestSuite) TesTimeSet() {
	set, e := NewTimeSet("times")
	s.NoError(e)

	n, e := set.Add(time.Now(), time.Time{})
	s.NoError(e)
	s.Equal(int64(2), n)
}

func (s *KredisTestSuite) TestTimeSetWithDefault() {
	t1 := time.Date(2021, 8, 15, 14, 30, 0, 0, time.UTC)
	t2 := time.Date(2021, 8, 15, 15, 30, 0, 0, time.UTC)

	set, e := NewTimeSetWithDefault("times", []time.Time{t1, t2})
	s.NoError(e)
	s.True(set.Includes(t1))
	s.True(set.Includes(t2))
	s.False(set.Includes(time.Now()))
}
