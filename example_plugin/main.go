package main

import (
	"github.com/natefinch/stdplug"
)

func main() {
	stdplug.Provide("Plugin", plugin{})
}

type plugin struct{}

func (plugin) SayHi(name string, response *string) error {
	*response = "Hi " + name
	return nil
}
