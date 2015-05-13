package main

import (
	"log"
	"os"
	"runtime"

	"github.com/natefinch/plugin"
)

func main() {
	log.SetPrefix("[host log] ")
	path := "example_driver"
	if runtime.GOOS == "windows" {
		path = path + ".exe"
	}

	s, err := plugin.StartDriver(os.Stderr, path)
	if err != nil {
		log.Fatalf("failed to start driver: %s", err)
	}
	if err := s.RegisterName("Host", api{}); err != nil {
		log.Fatalf("failed to register Host: %s", err)
	}
	if err := s.RegisterName("Host2", api2{}); err != nil {
		log.Fatalf("failed to register Host2: %s", err)
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
