package kredis

// From the Kredis code docs:
//
// A limiter is a specialized form of a counter that can be checked whether it
// has been exceeded and is provided fail safe. This means it can be used to
// guard login screens from brute force attacks without denying access in case
// Redis is offline.
type Limiter struct {
	Proxy
	limit uint64
}

func NewLimiter(key string, limit uint64, opts ...ProxyOption) (*Limiter, error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return nil, err
	}

	return &Limiter{Proxy: *proxy, limit: limit}, nil
}

// Similarly to the Poke func, the error be ignored, depending on the use case
func NewLimiterWithDefault(key string, limit uint64, defaultValue int64, opts ...ProxyOption) (l *Limiter, err error) {
	proxy, err := NewProxy(key, opts...)
	if err != nil {
		return
	}

	l = &Limiter{Proxy: *proxy, limit: limit}
	err = proxy.watch(func() error {
		return l.increment(defaultValue)
	})

	return
}

// Depending on the use case, the error can and should be ignored.
//
// You can "poke" yourself above the limit and should use the func IsExceeded()
// to know if the limit is reached.
func (l *Limiter) Poke() error {
	return l.increment(1)
}

// It's possible for this func to return false even if we're above the limit
// due to a potential race condition.
func (l *Limiter) IsExceeded() bool {
	v, err := l.client.Get(l.ctx, l.key).Uint64()
	if err != nil {
		// key either does not exist or Redis is down; return false as a failsafe
		return false
	}

	return v >= l.limit
}

func (l *Limiter) Reset() (err error) {
	_, err = l.client.Del(l.ctx, l.key).Result()
	return
}

func (l *Limiter) increment(by int64) (err error) {
	pipe := l.client.TxPipeline()
	if l.expiresIn > 0 {
		pipe.SetNX(l.ctx, l.key, 0, l.expiresIn)
	}

	pipe.IncrBy(l.ctx, l.key, 1)
	_, err = pipe.Exec(l.ctx)
	return
}
