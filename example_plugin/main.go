package main

import (
	"github.com/natefinch/plugin"
)

func main() {
	stdplug.Provide("Plugin", api{})
}

type api struct{}

func (api) SayHi(name string, response *string) error {
	*response = "Hi " + name
	return nil
}
