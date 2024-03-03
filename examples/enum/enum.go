package main

import (
	"github.com/mjc-gh/kredis-go"
)

func main() {
	kredis.SetConfiguration("shared", "ns", "redis://localhost:6379/2")
	kredis.EnableDebugLogging()

	enum, _ := kredis.NewEnum("enum", "go", []string{"ready", "set", "go"})
	enum.Is("go")
	enum.Value()

	err := enum.SetValue("set")
	err = enum.SetValue("error")
}
