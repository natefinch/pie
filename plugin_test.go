package plugin

import (
	"errors"
	"io"
	"os"
	"testing"
	"time"
)

var _ io.ReadWriteCloser = rwCloser{}
var _ io.ReadWriteCloser = ioPipe{}

func TestRWCloser(t *testing.T) {
	rc := &closeRW{}
	wc := &closeRW{}
	rwc := rwCloser{rc, wc}
	if err := rwc.Close(); err != nil {
		t.Errorf("unexpected error from rwCloser.Close: %#v", err)
	}
	if !rc.closed {
		t.Error("Close not called on ReadCloser.")
	}
	if !wc.closed {
		t.Error("Close not called on WriteCloser.")
	}
}

func TestRWCloserReadCloserError(t *testing.T) {
	readCloserErr := errors.New("read")
	rc := &closeRW{err: readCloserErr}
	wc := &closeRW{}
	rwc := rwCloser{rc, wc}
	err := rwc.Close()
	if !rc.closed {
		t.Error("Close not called on ReadCloser.")
	}
	if !wc.closed {
		t.Error("Close not called on WriteCloser.")
	}
	if err == nil {
		t.Error("ReadCloser error not passed through from rwCloser.Close")
	}
	if err != readCloserErr {
		t.Errorf("Different error returned from rwCloser than expected: %#v", err)
	}
}

func TestRWCloserWriteCloserError(t *testing.T) {
	writeCloserErr := errors.New("write")
	rc := &closeRW{}
	wc := &closeRW{err: writeCloserErr}
	rwc := rwCloser{rc, wc}
	err := rwc.Close()
	if !rc.closed {
		t.Error("Close not called on ReadCloser.")
	}
	if !wc.closed {
		t.Error("Close not called on WriteCloser.")
	}
	if err == nil {
		t.Error("ReadCloser error not passed through from rwCloser.Close")
	}
	if err != writeCloserErr {
		t.Errorf("Different error returned from rwCloser than expected: %#v", err)
	}
}

func TestRWCloserBothCloserError(t *testing.T) {
	writeCloserErr := errors.New("write")
	readCloserErr := errors.New("read")
	rc := &closeRW{err: readCloserErr}
	wc := &closeRW{err: writeCloserErr}
	rwc := rwCloser{rc, wc}
	err := rwc.Close()
	if !rc.closed {
		t.Error("Close not called on ReadCloser.")
	}
	if !wc.closed {
		t.Error("Close not called on WriteCloser.")
	}
	if err == nil {
		t.Error("Error not passed through from rwCloser.Close")
	}

	// I don't think we actually care which of these errors gets returned, as
	// long as one of them does.
	if err != writeCloserErr && err != readCloserErr {
		t.Errorf("Different error returned from rwCloser than expected: %#v", err)
	}
}

func TestIOPipeClose(t *testing.T) {
	rc := &closeRW{}
	wc := &closeRW{}
	p := &proc{}
	iop := ioPipe{rc, wc, p}
	if err := iop.Close(); err != nil {
		t.Errorf("Unexpected error from ioPipe.Close: %#v", err)
	}
	if !rc.closed {
		t.Error("Close not called on ReadCloser.")
	}
	if !wc.closed {
		t.Error("Close not called on WriteCloser.")
	}
	if p.sig == nil {
		t.Errorf("No signal sent to process")
	}
	if p.sig != os.Interrupt {
		t.Errorf("Unexpected signal sent to process, expected os.Interrupt, got %#v", p.sig)
	}
	if p.killed {
		t.Errorf("Kill() called unexpectedly on process.")
	}
}

func TestIOPipeSlowProc(t *testing.T) {
	defer func(d time.Duration) {
		procTimeout = d
	}(procTimeout)
	procTimeout = 5 * time.Millisecond
	rc := &closeRW{}
	wc := &closeRW{}
	p := &proc{delay: procTimeout * 2}
	iop := ioPipe{rc, wc, p}
	if err := iop.Close(); err != procStopTimeoutErr {
		t.Errorf("Unexpected error from ioPipe.Close, expected %#v, got: %#v", procStopTimeoutErr, err)
	}
	if !rc.closed {
		t.Error("Close not called on ReadCloser.")
	}
	if !wc.closed {
		t.Error("Close not called on WriteCloser.")
	}
	if p.sig == nil {
		t.Errorf("no signal sent to process")
	}
	if p.sig != os.Interrupt {
		t.Errorf("Unexpected signal sent to process, expected os.Interrupt, got %#v", p.sig)
	}
	if !p.killed {
		t.Errorf("Kill() unexpectedly not called on process.")
	}
}

// closeRW is a helper that fulfills io.Reader, io.Writer, and io.Closer for
// testing purposes.s
type closeRW struct {
	closed bool
	err    error
}

// Close fulfills io.Closer and will record that it was called, and return this
// value's error, if any.
func (c *closeRW) Close() error {
	c.closed = true
	return c.err
}

// Read fulfills io.Reader and does nothing.
func (*closeRW) Read(_ []byte) (int, error) {
	return 0, nil
}

// Write fulfills io.Writer and does nothing.
func (*closeRW) Write(_ []byte) (int, error) {
	return 0, nil
}

// proc is a helper that fullfills the osProcess interface for testing purposes.
type proc struct {
	delay     time.Duration
	waitErr   error
	killErr   error
	signalErr error
	sig       os.Signal
	killed    bool
}

// Wait will wait for delay time and then return waitErr.
func (p *proc) Wait() (*os.ProcessState, error) {
	<-time.After(p.delay)
	return nil, p.waitErr
}

// Kill returns killErr.
func (p *proc) Kill() error {
	p.killed = true
	return p.killErr
}

// Signal ignores the signal and returns signalErr.
func (p *proc) Signal(sig os.Signal) error {
	p.sig = sig
	return p.signalErr
}
