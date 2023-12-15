package kredis

import (
	"context"
	"errors"
	"fmt"
	"reflect"
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
	defaultValue any // TODO deprecate this field
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

	if reflect.ValueOf(options.DefaultValue).Kind() == reflect.Ptr {
		return nil, errors.New("default value cannot be a pointer")
	}

	return &Proxy{
		// TODO figure out the best way to handle context
		ctx:          context.TODO(),
		client:       client,
		key:          key,
		expiresIn:    duration,
		defaultValue: options.DefaultValue,
	}, nil
}
