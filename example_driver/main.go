package main

import (
	"log"
	"net/rpc"

	"github.com/natefinch/plugin"
)

func main() {
	log.SetPrefix("[plugin log] ")

	p := plug{plugin.Drive()}
	s, err := p.SayHi("plugin")
	if err != nil {
		log.Fatalf("failed saying hi: %s", err)
	}
	log.Println("Got response from host: ", s)
	s, err = p.SayBye("plugin")
	if err != nil {
		log.Fatalf("failed saying bye: %s", err)
	}
	log.Println("Got response from host: ", s)
}

type plug struct {
	client *rpc.Client
}

func (p plug) SayHi(name string) (result string, err error) {
	err = p.client.Call("Host.SayHi", name, &result)
	return result, err
}

func (p plug) SayBye(name string) (result string, err error) {
	err = p.client.Call("Host2.SayBye", name, &result)
	return result, err
}
