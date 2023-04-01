package kredis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// TODO Add Pipelining support
// https://redis.uptrace.dev/guide/go-redis-pipelines.html#pipelines

type Proxy struct {
	ctx          context.Context
	client       *redis.Client
	key          string
	expiresIn    time.Duration
	defaultValue any
}

func NewProxy(key string, options Options) (*Proxy, error) {
	client, namespace, err := getConnection(options.GetConfig())

	if err != nil {
		return nil, err
	}

	duration, _ := time.ParseDuration(options.ExpiresIn)

	if namespace != nil {
		key = fmt.Sprintf("%s:%s", *namespace, key)
	}

	return &Proxy{
		// TODO figure out the best way to handle context
		ctx:          context.Background(),
		client:       client,
		key:          key,
		expiresIn:    duration,
		defaultValue: options.DefaultValue,
	}, nil
}
