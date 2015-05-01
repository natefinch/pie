package main

import (
	"log"
	"net/rpc"
	"runtime"

	"github.com/natefinch/plugin"
)

func main() {
	log.SetPrefix("[master log] ")

	path := "example_plugin"
	if runtime.GOOS == "windows" {
		path = path + ".exe"
	}
	client, err := plugin.Start(path)
	if err != nil {
		log.Fatalf("Error running plugin: %s", err)
	}
	defer client.Close()
	p := plug{client}
	res, err := p.SayHi("master")
	if err != nil {
		log.Fatalf("error calling SayHi: %s", err)
	}
	log.Printf("Response from plugin: %q", res)

	res, err = p.SayHi("someone else")
	if err != nil {
		log.Fatalf("error calling SayHi: %s", err)
	}
	log.Printf("Response from plugin: %q", res)

}

type plug struct {
	client *rpc.Client
}

func (p plug) SayHi(name string) (result string, err error) {
	err = p.client.Call("Plugin.SayHi", name, &result)
	return result, err
}
