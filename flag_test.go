package kredis

import "time"

func (s *KredisTestSuite) TestFlag() {
	flag, err := NewFlag("flag")
	s.NoError(err)
	s.False(flag.IsMarked())

	s.NoError(flag.Mark())
	s.True(flag.IsMarked())

	s.NoError(flag.Remove())
	s.False(flag.IsMarked())
}

func (s *KredisTestSuite) TestMarkedFlag() {
	flag, err := NewMarkedFlag("flag")
	s.NoError(err)
	s.True(flag.IsMarked())
}

// TODO refactor test to check redis cmds with some sort of test env
// ProcessHook
func (s *KredisTestSuite) TestFlagWithMarkOptions() {
	s.T().Skip()

	flag, _ := NewFlag("flag_ex")

	s.NoError(flag.Mark(WithFlagExpiry("2ms")))
	s.True(flag.IsMarked())

	time.Sleep(1 * time.Millisecond)

	s.NoError(flag.Mark(WithFlagExpiry("2ms")))
	s.True(flag.IsMarked())

	time.Sleep(2 * time.Millisecond)

	s.False(flag.IsMarked())

	s.NoError(flag.Mark(WithFlagExpiry("2ms")))
	s.True(flag.IsMarked())

	time.Sleep(1 * time.Millisecond)

	s.NoError(flag.Mark(WithFlagExpiry("5ms")))
	s.True(flag.IsMarked())

	time.Sleep(2 * time.Millisecond)

	s.True(flag.IsMarked()) // still marked because of forced SET cmd
}
