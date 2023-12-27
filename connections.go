package kredis

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/redis/go-redis/v9"
)

type config struct {
	options   *redis.Options
	namespace string
}

var configs map[string]*config = map[string]*config{}
var connections map[string]*redis.Client = map[string]*redis.Client{}
var connectionLogging bool

func SetConnectionLogging(v bool) {
	connectionLogging = v
}

func SetConfiguration(name, namespace, url string) error {
	opt, err := redis.ParseURL(url)

	if err != nil {
		return err
	}

	// TODO handle redis settings
	//opt.ReadTimeout
	//opt.WriteTimeout
	//opt.PoolSize

	configs[name] = &config{options: opt, namespace: namespace}

	return nil
}

//func LoadConfigurations(filepath string) {}

func getConnection(name string) (*redis.Client, *string, error) {
	config, configured := configs[name]

	if !configured {
		return nil, nil, fmt.Errorf("%s is not a configured configuration", name)
	}

	conn, ok := connections[name]

	if ok {
		return conn, &config.namespace, nil
	}

	conn = redis.NewClient(config.options)

	if connectionLogging {
		conn.AddHook(newCmdLoggingHook())
	}

	connections[name] = conn

	//return conn
	return conn, &config.namespace, nil
}

type cmdLoggingHook struct{}

func newCmdLoggingHook() *cmdLoggingHook {
	return &cmdLoggingHook{}
}

func (c *cmdLoggingHook) DialHook(hook redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return hook(ctx, network, addr)
	}
}

var cmdColor = color.New(color.FgYellow).SprintFunc()
var keyColor = color.New(color.FgCyan).SprintFunc()
var argsColor = color.New(color.FgGreen).SprintFunc()

func (c *cmdLoggingHook) ProcessHook(hook redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		var key string
		var args []string

		cmdArgs := cmd.Args()
		name := cmdColor(strings.ToUpper(cmd.Name()))

		if len(cmdArgs) > 1 {
			key = keyColor(cmdArgs[1].(string))
		}

		if len(cmdArgs) > 2 {
			for _, cmdArg := range cmdArgs[2:] {
				args = append(args, argsColor(fmt.Sprintf("%v", cmdArg)))
			}
		}

		start := time.Now()
		err := hook(ctx, cmd)
		dur := float64(time.Since(start).Microseconds()) / 1000.0

		fmt.Printf("Kredis (%.1fms) %s %s %s\n", dur, name, key, strings.Join(args, " "))

		return err
	}
}

func (c *cmdLoggingHook) ProcessPipelineHook(hook redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return hook(ctx, cmds)
	}
}
