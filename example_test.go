package plugin_test

import (
	"log"
	"os"
	"strings"

	"net/rpc/jsonrpc"

	"github.com/natefinch/plugin"
)

// This function should be called from the master program that wants to run
// plugins to extend its functionality.
//
// StartProviderCodec starts a plugin at path "/var/lib/foo", using the JSON-RPC
// codec, and writing its output to this application's Stderr.  The application
// can then call methods on the rpc client returned using the standard rpc
// pattern.
func ExampleStartProviderCodec() {
	client, err := plugin.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, "/var/lib/foo")
	if err != nil {
		log.Fatalf("failed to load foo plugin: %s", err)
	}
	var reply string
	client.Call("Foo.ToUpper", "something", &reply)
}

// This function should be called from the plugin program that wants to provide
// functionality for the master program.
//
// NewProviderCodec starts an RPC server that reads from stdin and writes to
// stdout. It provides functions attached to the API value passed in.
// Server.Serve() will block forever, so it is common to simply put this at the
// end of the plugin's main function.
func ExampleProvider_ServeCodec() {
	p := plugin.NewProvider()
	if err := p.RegisterName("Foo", API{}); err != nil {
		log.Fatalf("can't register api: %s", err)
	}
	p.ServeCodec(jsonrpc.NewServerCodec)
}

// API is an example type to show how to serve methods over RPC.
type API struct{}

// ToUpper is an example function that gets served over RPC.  See net/rpc for
// details on how to server functionality over RPC.
func (API) ToUpper(input string, output *string) error {
	*output = strings.ToUpper(input)
	return nil
}
