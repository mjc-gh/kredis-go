package kredis

func (s *KredisTestSuite) TestBoolHash() {
	hash, e := NewBoolHash("bools")
	s.NoError(e)

	s.NoError(hash.Set("a", true))
	s.NoError(hash.Set("b", true))
	s.NoError(hash.Set("c", false))

	b, ok := hash.Get("a")
	s.True(ok)
	s.True(b)

	_, ok = hash.Get("x")
	s.False(ok)

	entries, e := hash.Entries()
	s.NoError(e)
	s.Equal(map[string]bool{
		"a": true, "b": true, "c": false,
	}, entries)

	valuesAt, e := hash.ValuesAt("c", "a")
	s.NoError(e)
	s.Equal([]bool{false, true}, valuesAt)

	keys, e := hash.Keys()
	s.NoError(e)
	s.Equal([]string{"a", "b", "c"}, keys)

	values, e := hash.Values()
	s.NoError(e)
	s.Equal([]bool{true, true, false}, values)

	n, e := hash.Delete("a", "b", "x")
	s.NoError(e)
	s.Equal(int64(2), n)

	n, e = hash.Update(map[string]bool{"x": true, "y": false})
	s.NoError(e)
	s.Equal(int64(2), n)

	entries, e = hash.Entries()
	s.NoError(e)
	s.Equal(map[string]bool{
		"c": false, "x": true, "y": false,
	}, entries)

	s.NoError(hash.Clear())
}
