# plugin [![GoDoc](https://godoc.org/github.com/natefinch/plugin?status.png)](https://godoc.org/github.com/natefinch/plugin) [![Build Status](https://drone.io/github.com/natefinch/plugin/status.png)](https://drone.io/github.com/natefinch/plugin/latest)

    import "github.com/natefinch/plugin"

package plugin provides a toolkit for creating plugins for Go applications.

This is a proof of concept, and still needs a lot of TLC.

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


## func Start
``` go
func Start(path string, w io.Writer) (*rpc.Client, error)
```
Start starts a plugin application at the given path and returns an RPC client
that communicates using gob encoding.  It writes to the plugin's Stdin and
reads from the plugin's Stdout.  The writer passed to w will receive
stderr output from the plugin.  Closing the RPC client returned from this
function will shut down the plugin's process.


## func StartWithCodec
``` go
func StartWithCodec(newClientCodec func(io.ReadWriteCloser) rpc.ClientCodec, path string, w io.Writer) (*rpc.Client, error)
```
StartWithCodec starts a plugin application at the given path and returns an
RPC client that communicates using the ClientCodec returned by
newClientCodec.  It writes to the plugin's Stdin and reads from the
plugin's Stdout.  The writer passed to w will receive stderr output from the
plugin.  Closing the RPC client returned from this function will shut down
the plugin's process.


## type Server
``` go
type Server struct {
    // contains filtered or unexported fields
}
```
Server is a value that will allow you to register types for the API of a
plugin and then serve those types over RPC using Stdin and Stdout.


### func NewServer
``` go
func NewServer() Server
```
NewServer returns an RPC plugin server that will serve RPC over Stdin and Stdout
using gob encoding.


### func NewServerWithCodec
``` go
func NewServerWithCodec(newServerCodec func(io.ReadWriteCloser) rpc.ServerCodec) Server
```
NewServerWithCodec returns an RPC plugin server that will serve RPC over Stdin and
Stdout using the codec returned from newServerCodec


### func (Server) Register
``` go
func (s Server) Register(rcvr interface{}) error
```
Register functions hust like net/rpc.Server's Register.


### func (Server) RegisterName
``` go
func (s Server) RegisterName(name string, rcvr interface{}) error
```
RegisterName functions just like net/rpc.Server's RegisterName.


### func (Server) Serve
``` go
func (s Server) Serve()
```
Serve starts the RPC server, listening on Stdin and writing to Stdout.  This
call will block until the client hangs up.



