package main

import (
	"log"

	"github.com/natefinch/plugin"
)

func main() {
	log.SetPrefix("[plugin log] ")

	plugin.Provide("Plugin", api{})
}

type api struct{}

func (api) SayHi(name string, response *string) error {
	log.Printf("got call for SayHi with name %q", name)

	*response = "Hi " + name
	return nil
}
