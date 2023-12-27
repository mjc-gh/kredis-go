package kredis

import "time"

// variadic configuration functions and structs

// general proxy options (all factories)

type ProxyOptions struct {
	config    *string
	expiresIn *time.Duration
}

func (o *ProxyOptions) Config() string {
	if o.config != nil {
		return *o.config
	}

	return "shared"
}

func (o *ProxyOptions) ExpiresIn() time.Duration {
	if o.expiresIn != nil {
		return *o.expiresIn
	}

	return time.Duration(0)
}

type ProxyOption func(*ProxyOptions)

func WithConfigName(name string) ProxyOption {
	return func(o *ProxyOptions) {
		o.config = &name
	}
}
func WithExpiry(expires string) ProxyOption {
	return func(o *ProxyOptions) {
		duration, err := time.ParseDuration(expires)
		if err != nil {
			return
		}

		o.expiresIn = &duration
	}
}

// TODO WithContext ProxyOption for providing a user-defined context.Context for the Proxy
// struct and thus the redis client commands

// func WithContext(ctx context.Context) ProxyOption {
//   return func(o *ProxyOptions) {
//   }
// }

// For range options (list, unique lists)

type RangeOptions struct {
	start int64
}

type RangeOption func(*RangeOptions)

func WithRangeStart(s int64) RangeOption {
	return func(o *RangeOptions) {
		o.start = s
	}
}
