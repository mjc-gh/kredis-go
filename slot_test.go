package kredis

func (s *KredisTestSuite) TestSlot() {
	slot, err := NewSlot("slot", 3)
	s.NoError(err)

	s.True(slot.IsAvailable())
	s.True(slot.Reserve())
	s.True(slot.Reserve())
	s.True(slot.Reserve())

	s.Equal(int64(3), slot.Taken())
	s.False(slot.IsAvailable())

	s.True(slot.Release())
	s.Equal(int64(2), slot.Taken())
	s.True(slot.IsAvailable())

	s.NoError(slot.Reset())
	s.Equal(int64(0), slot.Taken())
	s.True(slot.IsAvailable())
	s.False(slot.Release())
}

func (s KredisTestSuite) TestSlotExpiry() {
	slot, err := NewSlot("slot", 1, WithConfigName("capture"), WithExpiry("5s"))
	s.NoError(err)

	slot.Reserve()
	s.False(slot.IsAvailable())

	s.Equal("EXPIRE ns:slot 5", s.captureHook.Captures[2].String())
}

func (s *KredisTestSuite) TestSlotWithReserveCallback() {
	var called int

	callback := func() {
		called += 1
	}

	slot, err := NewSlot("slot", 1)
	s.NoError(err)

	s.True(slot.Reserve(callback))
	s.Equal(1, called)
	s.Equal(int64(0), slot.Taken())

	slot.Reserve()

	s.False(slot.Reserve(callback))
	s.Equal(1, called)
	s.Equal(int64(0), slot.Taken())
}
