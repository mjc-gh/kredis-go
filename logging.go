package kredis

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/redis/go-redis/extra/rediscmd/v9"
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

func SetDebugLogger(userLogger logging) {
	debugLogger = userLogger
}

func (log stdLogger) Info(cmd redis.Cmder, dur time.Duration) {
	name, key, args := cmdLogParts(cmd)

	if dur == 0 {
		fmt.Printf("Kredis (tx)    %s %s %s\n", name, key, args)
	} else {
		msec := float64(dur.Microseconds()) / 1000.0

		fmt.Printf("Kredis (%.1fms) %s %s %s\n", msec, name, key, args)
	}
}

func (log stdLogger) Warn(fnName string, err error) {
	fmt.Printf("Kredis error in %s: %s", fnName, err.Error())
}

var cmdColor = color.New(color.FgYellow).SprintFunc()
var keyColor = color.New(color.FgCyan).SprintFunc()
var argsColor = color.New(color.FgGreen).SprintFunc()

func cmdLogParts(cmd redis.Cmder) (string, string, string) {
	var name string
	var key string
	var args []string

	cmdStr := rediscmd.CmdString(cmd)
	for idx, str := range strings.Split(cmdStr, " ") {
		if idx == 0 {
			name = cmdColor(strings.ToUpper(str))
		} else if idx == 1 {
			key = keyColor(str)
		} else {
			args = append(args, argsColor(str))
		}
	}

	return name, key, strings.Join(args, " ")
}
