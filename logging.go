package kredis

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/redis/go-redis/v9"
)

type logging interface {
	Info(redis.Cmder, time.Duration)
	Warn(string, error)
}

type stdLogger struct{}

var debugLogger logging

// Enable logging of Redis commands and errors that are ignored by functions
// with failsafes. This is useful for development and tests. Logging is
// disabled by default.
func EnableDebugLogging() {
	debugLogger = stdLogger{}
}

// Turn off debug logging.
func DisableDebugLogging() {
	debugLogger = nil
}

// TODO add a way to configure a user provided value that implements the logging interface
// func SetCommandLogger(userLogger cmdLogging)

func (log stdLogger) Info(cmd redis.Cmder, dur time.Duration) {
	name, key, args := cmdLogParts(cmd)

	if dur == 0 {
		fmt.Printf("Kredis (tx)    %s %s %s\n", name, key, args)
	} else {
		msec := float64(dur.Microseconds()) / 1000.0

		fmt.Printf("Kredis (%.1fms) %s %s %s\n", msec, name, key, args)
	}
}

func (log stdLogger) Warn(msg string, err error) {

}

var cmdColor = color.New(color.FgYellow).SprintFunc()
var keyColor = color.New(color.FgCyan).SprintFunc()
var argsColor = color.New(color.FgGreen).SprintFunc()

func cmdLogParts(cmd redis.Cmder) (string, string, string) {
	var key string
	var args []string

	cmdArgs := cmd.Args()
	name := cmdColor(strings.ToUpper(cmd.Name()))

	if len(cmdArgs) > 1 {
		switch cmdArgs[1].(type) {
		case int64:
			key = keyColor(cmdArgs[1].(int64))
		default:
			key = keyColor(cmdArgs[1].(string))
		}
	}

	if len(cmdArgs) > 2 {
		for _, cmdArg := range cmdArgs[2:] {
			args = append(args, argsColor(fmt.Sprintf("%v", cmdArg)))
		}
	}

	return name, key, strings.Join(args, " ")
}
