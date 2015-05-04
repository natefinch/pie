// package plugin provides a toolkit for creating plugins for Go applications.

// This is a proof of concept, and still needs a lot of TLC.

// plugin uses go's standard net/rpc to wrap os.Stdin and os.Stdout from a
// subprocess, to enable easy RPC with zero configuration.

// The plugin is implemented as a regular Go application, which provides a Go type
// to be used over RPC.  The master process (i.e. your main application, into which
// the plugin is plugging) then runs the plugin as a subprocess, and uses RPC to
// communicate with the subprocess.  From a developers' point of view, this makes
// using functionality from plugins very simple.

// Included in this repo is a very simple example of a master process and a plugin
// process, to see how the library can be used.  example_master expects
// example_plugin to be in the same directory.  You can just go install both of
// them, and it'll work correctly.

// The really nice thing about this library is how simple the code is for the
// plugins. It's just a few very simple lines of boilerplate.

package plugin
