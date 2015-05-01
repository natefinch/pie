# plugin [![GoDoc](https://godoc.org/github.com/natefinch/plugin?status.png)](https://godoc.org/github.com/natefinch/plugin) [![Build Status](https://drone.io/github.com/natefinch/plugin/status.png)](https://drone.io/github.com/natefinch/plugin/latest)

### A framework for go plugins that use RPC over stdin/stdout

This is a proof of concept, and still needs a lot of TlC.

plugin uses go's standard net/rpc to wrap os.Stdin and os.Stdout from a
subprocess, to enable easy RPC with zero configuration.

The plugin is implemented as a regular Go application, which provides a Go type
to be used over RPC.  The master process (i.e. your main application, into which
the plugin is plugging) then runs the plugin as a subprocess, and uses RPC to
communicate with the subprocess.  From a developers' point of view, this makes
using functionality from plugins very simple.

Included in this repo is a very simple example of a master process and a plugin
process, to see how the library can be used.  example_master expects
example_plugin to be in the same directory.  You can just go install both of
them, and it'll work correctly.

The really nice thing about this library is how simple the code is for the
plugins. It's just a few very simple lines of boilerplate.

# Godoc

## func Provide
``` go
func Provide(name string, rcvr interface{})
```

Provide starts an rpc server providing the given interface over Stdin and
Stdout.  This call will block forever.


## func Start
``` go
func Start(path string) (client *rpc.Client, err error)
```

Start starts a plugin application at the given path and returns an RPC client
that talks to it over Stdin and Stdout.



