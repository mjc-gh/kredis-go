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

func (s *KredisTestSuite) TestFlagWithMarkOptions() {
	flag, _ := NewFlag("flag_ex")

	s.NoError(flag.Mark(WithFlagMarkExpiry("2ms")))
	s.True(flag.IsMarked())

	time.Sleep(1 * time.Millisecond)

	s.NoError(flag.Mark(WithFlagMarkExpiry("2ms")))
	s.True(flag.IsMarked())

	time.Sleep(2 * time.Millisecond)

	s.False(flag.IsMarked())

	s.NoError(flag.Mark(WithFlagMarkExpiry("2ms")))
	s.True(flag.IsMarked())

	time.Sleep(1 * time.Millisecond)

	s.NoError(flag.Mark(WithFlagMarkExpiry("5ms"), WithFlagMarkForced()))
	s.True(flag.IsMarked())

	time.Sleep(2 * time.Millisecond)

	s.True(flag.IsMarked()) // still marked because of forced SET cmd
}
