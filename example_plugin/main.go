package main

import (
	"log"

	"github.com/natefinch/plugin"
)

func main() {
	log.SetPrefix("[plugin log] ")

	s := plugin.NewServer()
	if err := s.RegisterName("Plugin", api{}); err != nil {
		log.Fatalf("failed to register Plugin: %s", err)
	}
	if err := s.RegisterName("Plugin2", api2{}); err != nil {
		log.Fatalf("failed to register Plugin2: %s", err)
	}
	s.Serve()
}

type api struct{}

func (api) SayHi(name string, response *string) error {
	log.Printf("got call for SayHi with name %q", name)

	*response = "Hi " + name
	return nil
}

type api2 struct{}

func (api2) SayBye(name string, response *string) error {
	log.Printf("got call for SayBye with name %q", name)

	*response = "Bye " + name
	return nil
}
