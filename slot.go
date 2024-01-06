package kredis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Slot struct {
	Proxy
	available int64
}

func NewSlot(key string, available int64, opts ...ProxyOption) (*Slot, error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &Slot{Proxy: *proxy, available: available}, nil
}

type SlotCallback func()

func (s *Slot) Reserve(callbacks ...SlotCallback) (reserved bool) {
	// no callback
	if len(callbacks) == 0 {
		if s.IsAvailable() {
			s.incr()
			reserved = true
		}

		return
	}

	// callback given
	if s.Reserve() {
		reserved = true

		for _, callback := range callbacks {
			callback()
		}
	}

	// alawys call Release when given a callback, even if the callback is not
	// invoked and there is nothing to Reserve
	s.Release()

	return
}

func (s *Slot) Release() bool {
	if s.Taken() > 0 {
		s.decr()
		return true
	}

	return false
}

func (s *Slot) IsAvailable() bool {
	taken, err := s.TakenResult()
	if err != nil && err != redis.Nil {
		return false // failsafe
	}

	return taken < int64(s.available)
}

func (s *Slot) Reset() error {
	_, err := s.client.Del(s.ctx, s.key).Result()
	if err != nil {
		return err
	}

	return nil
}

func (s *Slot) Taken() int64 {
	t, err := s.TakenResult()
	if err != nil {
		return 0
	}

	return t
}

func (s *Slot) TakenResult() (t int64, err error) {
	t, err = s.client.Get(s.ctx, s.key).Int64()
	return
}

func (s *Slot) incr() {
	_, err := s.client.Incr(s.ctx, s.key).Result()
	if err != nil {
		// TODO debug logging
		fmt.Println(err)
	}

	s.RefreshTTL()
}

func (s *Slot) decr() {
	_, err := s.client.Decr(s.ctx, s.key).Result()
	if err != nil {
		// TODO debug logging
		fmt.Println(err)
	}

	s.RefreshTTL()
}
