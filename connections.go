package kredis

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
)

type config struct {
	options   *redis.Options
	namespace string
}

var configs map[string]*config = map[string]*config{}
var connections map[string]*redis.Client = map[string]*redis.Client{}

type RedisOption func(*redis.Options)

func SetConfiguration(name, namespace, url string, opts ...RedisOption) error {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return err
	}

	for _, optFn := range opts {
		optFn(opt)
	}

	configs[name] = &config{options: opt, namespace: namespace}
	return nil
}

func getConnection(name string) (*redis.Client, string, error) {
	config, configured := configs[name]
	if !configured {
		return nil, "", fmt.Errorf("%s is not a configured configuration", name)
	}

	conn, ok := connections[name]
	if ok {
		return conn, config.namespace, nil
	}

	conn = redis.NewClient(config.options)
	connections[name] = conn

	if debugLogger != nil {
		conn.AddHook(newCmdLoggingHook(debugLogger))
	}

	return conn, config.namespace, nil
}

type cmdLoggingHook struct {
	cmdLogger logging
}

func newCmdLoggingHook(clog logging) *cmdLoggingHook {
	return &cmdLoggingHook{clog}
}

func (c *cmdLoggingHook) DialHook(hook redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return hook(ctx, network, addr)
	}
}

func (c *cmdLoggingHook) ProcessHook(hook redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()
		err := hook(ctx, cmd)

		c.cmdLogger.Info(cmd, time.Since(start))

		return err
	}
}

func (c *cmdLoggingHook) ProcessPipelineHook(hook redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		start := time.Now()
		err := hook(ctx, cmds)

		for idx, cmd := range cmds {
			if idx == len(cmds)-1 {
				c.cmdLogger.Info(cmd, time.Since(start))
			} else {
				c.cmdLogger.Info(cmd, time.Duration(0))
			}
		}

		return err
	}
}
