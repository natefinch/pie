// package plugin provides helper functions for creating plugins using RPC over
// stdin/stdout.
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

// Provide starts an rpc server providing the given interface over Stdin and Stdout.  This call will block forever.
func Provide(name string, rcvr interface{}) {
	s := rpc.NewServer()
	s.RegisterName(name, rcvr)
	s.ServeConn(rwCloser{os.Stdin, os.Stdout})
}

// Start starts a plugin application at the given path and returns an RPC client
// that talks to it over Stdin and Stdout.
func Start(path string) (client *rpc.Client, err error) {
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

	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return rpc.NewClient(ioPipe{out, in, cmd.Process}), nil
}

type ioPipe struct {
	io.ReadCloser
	io.WriteCloser
	proc *os.Process
}

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

func (iop ioPipe) closeProc() error {
	result := make(chan error, 1)
	go func() { _, err := iop.proc.Wait(); result <- err }()
	if err := iop.proc.Signal(os.Interrupt); err != nil {
		return err
	}
	select {
	case err := <-result:
		return err
	case <-time.After(time.Second):
		if err := iop.proc.Kill(); err != nil {
			return fmt.Errorf("error killing process after timeout: %s", err)
		}
		return errors.New("timed out waiting for process to stop")
	}
}

type rwCloser struct {
	io.ReadCloser
	io.WriteCloser
}

func (rw rwCloser) Close() error {
	err := rw.ReadCloser.Close()
	if err := rw.WriteCloser.Close(); err != nil {
		return err
	}
	return err
}
