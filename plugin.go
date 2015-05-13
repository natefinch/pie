package plugin

import (
	"errors"
	"fmt"
	"io"
	"net/rpc"
	"os"
	"os/exec"
	"time"
)

var procStopTimeoutErr = errors.New("process killed after timeout waiting for process to stop")

// NewServer returns an RPC server that will serve RPC over this application's
// Stdin and Stdout using gob encoding.
func NewServer() Server {
	return Server{
		server: rpc.NewServer(),
		rwc:    rwCloser{os.Stdin, os.Stdout},
	}
}

// NewServerWithCodec returns an RPC server that will serve RPC over this
// application's Stdin and Stdout using the ServerCodec returned from codec.
func NewServerWithCodec(codec func(io.ReadWriteCloser) rpc.ServerCodec) Server {
	return Server{
		server: rpc.NewServer(),
		codec:  codec,
		rwc:    rwCloser{os.Stdin, os.Stdout},
	}
}

// Server is a type that will allow you to register types for the API of a
// plugin and then serve those types over RPC.  It encompasses the functionality
// to talk to a plugin/master process.
type Server struct {
	server *rpc.Server
	codec  func(io.ReadWriteCloser) rpc.ServerCodec
	rwc    io.ReadWriteCloser
}

// Serve starts the RPC server.  This call will block until the client hangs up.
func (s Server) Serve() {
	if s.codec != nil {
		s.server.ServeCodec(s.codec(s.rwc))
	}
	s.server.ServeConn(s.rwc)
}

// Register publishes in the server the set of methods of the receiver value
// that satisfy the following conditions:
//
//	- exported method
//	- two arguments, both of exported type
//	- the second argument is a pointer
//	- one return value, of type error
//
// It returns an error if the receiver is not an exported type or has no
// suitable methods. It also logs the error using package log. The client
// accesses each method using a string of the form "Type.Method", where Type is
// the receiver's concrete type.
func (s Server) Register(rcvr interface{}) error {
	return s.server.Register(rcvr)
}

// RegisterName is like Register but uses the provided name for the type
// instead of the receiver's concrete type.
func (s Server) RegisterName(name string, rcvr interface{}) error {
	return s.server.RegisterName(name, rcvr)
}

// Start starts an application (plugin) at the given path and returns an RPC
// client that communicates using gob encoding.  It writes to the plugin's Stdin
// and reads from the plugin's Stdout.  The writer passed to w will receive
// stderr output from the plugin.  Closing the RPC client returned from this
// function will shut down the plugin's process.
func Start(path string, w io.Writer) (*rpc.Client, error) {
	rwc, err := start(path, w)
	if err != nil {
		return nil, err
	}
	return rpc.NewClient(rwc), nil
}

// StartWithCodec starts an application (plugin) at the given path and returns
// an RPC client that communicates using the ClientCodec returned by codec.  It
// writes to the plugin's Stdin and reads from the plugin's Stdout.  The writer
// passed to w will receive stderr output from the plugin.  Closing the RPC
// client returned from this function will shut down the plugin's process.
func StartWithCodec(codec func(io.ReadWriteCloser) rpc.ClientCodec, path string, w io.Writer) (*rpc.Client, error) {
	rwc, err := start(path, w)
	if err != nil {
		return nil, err
	}
	return rpc.NewClientWithCodec(codec(rwc)), nil
}

// StartDriver starts a plugin application that consumes an API this application
// provides.  In effect, the plugin is "driving" this application.
func StartDriver(path string, w io.Writer) (Server, error) {
	rwc, err := start(path, w)
	if err != nil {
		return Server{}, err
	}
	return Server{
		server: rpc.NewServer(),
		rwc:    rwc,
	}, nil
}

// StartDriverWithCodec starts a plugin application that consumes an API this
// application provides using RPC with the ServerCodec returned by codec.  In
// effect, the plugin is "driving" this application.
func StartDriverWithCodec(codec func(io.ReadWriteCloser) rpc.ServerCodec, path string, w io.Writer) (Server, error) {
	rwc, err := start(path, w)
	if err != nil {
		return Server{}, err
	}
	return Server{
		server: rpc.NewServer(),
		codec:  codec,
		rwc:    rwc,
	}, nil
}

// Drive returns an rpc.Client that will drive the host process over this
// application's Stdin and Stdout using gob encoding.
func Drive() *rpc.Client {
	return rpc.NewClient(rwCloser{os.Stdin, os.Stdout})
}

// DriveWithCodec returs an rpc.Client that will drive the host process over
// this application's Stdin and Stdout using the ClientCodec returned by codec.
func DriveWithCodec(codec func(io.ReadWriteCloser) rpc.ClientCodec) *rpc.Client {
	return rpc.NewClientWithCodec(codec(rwCloser{os.Stdin, os.Stdout}))
}

// start runs the plugin and returns a ReadWriteCloser that can be used to
// control the plugin.
func start(path string, w io.Writer) (io.ReadWriteCloser, error) {
	cmd := exec.Command(path)
	in, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			in.Close()
		}
	}()
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			out.Close()
		}
	}()

	cmd.Stderr = w
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return ioPipe{out, in, cmd.Process}, nil
}

// osProcess is an interface that is fullfilled by *os.Process and makes our
// testing a little easier.
type osProcess interface {
	Wait() (*os.ProcessState, error)
	Kill() error
	Signal(os.Signal) error
}

// ioPipe simply wraps a ReadCloser, WriteCloser, and a Process, and coordinates
// them so they all close together.
type ioPipe struct {
	io.ReadCloser
	io.WriteCloser
	proc osProcess
}

// Close closes the pipe's WriteCloser, ReadClosers, and process.
func (iop ioPipe) Close() error {
	err := iop.ReadCloser.Close()
	if writeErr := iop.WriteCloser.Close(); writeErr != nil {
		err = writeErr
	}
	if procErr := iop.closeProc(); procErr != nil {
		err = procErr
	}
	return err
}

// procTimeout is the timeout to wait for a process to stop after being
// signalled.  It is adjustable to keep tests fast.
var procTimeout = time.Second

// closeProc sends an interrupt signal to the pipe's process, and if it doesn't
// respond in one second, kills the process.
func (iop ioPipe) closeProc() error {
	result := make(chan error, 1)
	go func() { _, err := iop.proc.Wait(); result <- err }()
	if err := iop.proc.Signal(os.Interrupt); err != nil {
		return err
	}
	select {
	case err := <-result:
		return err
	case <-time.After(procTimeout):
		if err := iop.proc.Kill(); err != nil {
			return fmt.Errorf("error killing process after timeout: %s", err)
		}
		return procStopTimeoutErr
	}
}

// rwCloser just merges a ReadCloser and a WriteCloser into a ReadWriteCloser.
type rwCloser struct {
	io.ReadCloser
	io.WriteCloser
}

// Close closes both the ReadCloser and the WriteCloser, returning the last
// error from either.
func (rw rwCloser) Close() error {
	err := rw.ReadCloser.Close()
	if err := rw.WriteCloser.Close(); err != nil {
		return err
	}
	return err
}
