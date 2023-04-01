package kredis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type config struct {
	options   *redis.Options
	namespace string
}

var configs map[string]*config = map[string]*config{}
var connections map[string]*redis.Client = map[string]*redis.Client{}

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
	connections[name] = conn

	//return conn
	return conn, &config.namespace, nil
}
