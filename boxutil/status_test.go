package boxutil

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/matryer/is"
)

func init() {
	// quicker for testing
	readyCheckInterval = 100 * time.Millisecond
}

func TestStatusChan(t *testing.T) {
	is := is.New(t)

	i := &testBox{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	status := StatusChan(ctx, i)
	is.Equal(<-status, "starting...")
	i.setReady()
	is.Equal(<-status, "ready")
	i.setError()
	is.Equal(<-status, "unavailable")
	i.clearError()
	is.Equal(<-status, "ready")

}

func TestWaitForReady(t *testing.T) {
	is := is.New(t)
	i := &testBox{}
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(300*time.Millisecond, cancel)
	go func() {
		i.setReady()
	}()
	err := WaitForReady(ctx, i)
	is.NoErr(err)
}

func TestWaitForReadyTimeout(t *testing.T) {
	is := is.New(t)
	i := &testBox{}
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(100*time.Millisecond, cancel)
	go func() {
		time.Sleep(200 * time.Millisecond)
		i.setReady()
	}()
	err := WaitForReady(ctx, i)
	is.Equal(err, context.Canceled)
}

type testBox struct {
	lock  sync.Mutex
	ready bool
	err   error
}

func (i *testBox) setReady() {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.ready = true
}

func (i *testBox) setError() {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.err = errors.New("cannot reach server")
}

func (i *testBox) clearError() {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.err = nil
}

func (i *testBox) Info() (*Info, error) {
	i.lock.Lock()
	defer i.lock.Unlock()
	if i.err != nil {
		return nil, i.err
	}
	if i.ready {
		return &Info{Status: "ready"}, nil
	}
	return &Info{Status: "starting..."}, nil
}
