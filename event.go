package wait

import (
	"context"
	"errors"
	"time"
)

var (
	ErrTimeout = errors.New("Time out ")
)

type Event interface {
	WaitContext(ctx context.Context) (interface{}, error)
	WaitTimeout(timeout time.Duration) (interface{}, error)
	WaitContextTimeout(ctx context.Context, timeout time.Duration) (interface{}, error)
	Trigger(value interface{})
}

func NewEvent() Event {
	w := &event{}
	w.c = make(chan interface{})
	return w
}

type event struct {
	c chan interface{}
}

func (e *event) WaitContext(ctx context.Context) (interface{}, error) {
	select {
	case obj := <-e.c:
		return obj, nil
	case <-ctx.Done():
		return nil, context.Canceled
	}
}

func (e *event) WaitTimeout(timeout time.Duration) (interface{}, error) {
	select {
	case obj := <-e.c:
		return obj, nil
	case <-time.After(timeout):
		return nil, ErrTimeout
	}
}

func (e *event) WaitContextTimeout(ctx context.Context, timeout time.Duration) (interface{}, error) {
	select {
	case obj := <-e.c:
		return obj, nil
	case <-ctx.Done():
		return nil, context.Canceled
	case <-time.After(timeout):
		return nil, ErrTimeout
	}
}

func (e *event) Trigger(value interface{}) {
	select {
	case e.c <- value:
	default:
	}
}
