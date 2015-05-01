package main

import (
	"github.com/natefinch/plugin"
)

func main() {
	plugin.Provide("Plugin", api{})
}

type api struct{}

func (api) SayHi(name string, response *string) error {
	*response = "Hi " + name
	return nil
}
