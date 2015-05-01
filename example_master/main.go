package main

import (
	"log"
	"net/rpc"
	"os"
	"path/filepath"
	"runtime"

	"github.com/natefinch/stdplug"
)

func main() {
	path := filepath.Join(os.Getenv("GOPATH"), "bin", "example_plugin")
	if runtime.GOOS == "windows" {
		path = path + ".exe"
	}
	client, err := stdplug.Start(path)
	if err != nil {
		log.Fatalf("Error running plugin: %s", err)
	}
	defer client.Close()
	p := plugin{client}
	res, err := p.SayHi("master")
	if err != nil {
		log.Fatalf("error calling SayHi: %s", err)
	}
	log.Printf("Response from plugin: %q", res)
}

type plugin struct {
	client *rpc.Client
}

func (p plugin) SayHi(name string) (result string, err error) {
	err = p.client.Call("Plugin.SayHi", name, &result)
	return result, err
}
