package kredis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *KredisTestSuite) TestGetConfigurationUsesCacheMap() {
	_, ok := connections["shared"]
	s.False(ok)

	c, _, e := getConnection("shared")
	s.NoError(e)

	c2, _, e := getConnection("shared")
	s.Same(c, c2)
}

func (s *KredisTestSuite) TestSetConfigurationWithRedisOptions() {
	SetConfiguration("test", "", "redis://localhost:6379/0", func(opts *redis.Options) {
		opts.ReadTimeout = time.Duration(1)
		opts.WriteTimeout = time.Duration(1)
		opts.PoolSize = 1
	})

	cfg := configs["test"]

	s.Equal(time.Duration(1), cfg.options.ReadTimeout)
	s.Equal(time.Duration(1), cfg.options.WriteTimeout)
	s.Equal(1, cfg.options.PoolSize)

	delete(configs, "test")
}
