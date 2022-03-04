package wait

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

const defaultTimeout = time.Second * 30

// Wait
// key must sync.Map support
type Wait interface {
	SetTimeout(timeout time.Duration) Wait
	SetContext(ctx context.Context) Wait

	InitKey(keys ...interface{})
	Wait(key interface{}) (interface{}, error)
	WaitContext(ctx context.Context, key interface{}) (interface{}, error)
	WaitTimeout(timeout time.Duration, key interface{}) (interface{}, error)
	WaitContextTimeout(ctx context.Context, timeout time.Duration, key interface{}) (interface{}, error)

	WaitAny(keys ...interface{}) (Result, error)
	WaitAnyContext(ctx context.Context, keys ...interface{}) (Result, error)
	WaitAnyTimeout(timeout time.Duration, keys ...interface{}) (Result, error)
	WaitAnyContextTimeout(ctx context.Context, timeout time.Duration, keys ...interface{}) (Result, error)

	WaitAll(keys ...interface{}) (Result, error)
	WaitAllContext(ctx context.Context, keys ...interface{}) (Result, error)
	WaitAllTimeout(timeout time.Duration, keys ...interface{}) (Result, error)
	WaitAllContextTimeout(ctx context.Context, timeout time.Duration, keys ...interface{}) (Result, error)

	Trigger(key interface{})
	TriggerValue(key, value interface{})
}

func NewWait() Wait {
	w := &wait{
		timeout: defaultTimeout,
		ctx:     context.Background(),
		m:       sync.Map{},
	}

	return w
}

type wait struct {
	timeout time.Duration
	ctx     context.Context
	logger  log.Logger
	m       sync.Map
}

type waitResult struct {
	Result
	error
}

func (w *wait) SetTimeout(timeout time.Duration) Wait {
	w.timeout = timeout
	return w
}

func (w *wait) SetContext(ctx context.Context) Wait {
	w.ctx = ctx
	return w
}

func (w *wait) Wait(key interface{}) (interface{}, error) {
	return w.WaitContextTimeout(w.ctx, w.timeout, key)
}

func (w *wait) WaitContext(ctx context.Context, key interface{}) (interface{}, error) {
	return w.WaitContextTimeout(ctx, w.timeout, key)
}

func (w *wait) WaitTimeout(timeout time.Duration, key interface{}) (interface{}, error) {
	return w.WaitContextTimeout(w.ctx, timeout, key)
}

func (w *wait) WaitContextTimeout(ctx context.Context, timeout time.Duration, key interface{}) (interface{}, error) {
	e, ok := w.getEvent(key)
	if !ok {
		return nil, fmt.Errorf("Not init key %v ", key)
	}

	return e.WaitContextTimeout(ctx, timeout)
}

func (w *wait) WaitAny(keys ...interface{}) (Result, error) {
	return w.WaitAnyContextTimeout(w.ctx, w.timeout, keys...)
}

func (w *wait) WaitAnyContext(ctx context.Context, keys ...interface{}) (Result, error) {
	return w.WaitAnyContextTimeout(ctx, w.timeout, keys...)
}

func (w *wait) WaitAnyTimeout(timeout time.Duration, keys ...interface{}) (Result, error) {
	return w.WaitAnyContextTimeout(w.ctx, timeout, keys...)
}

func (w *wait) WaitAnyContextTimeout(ctx context.Context, timeout time.Duration, keys ...interface{}) (Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	ch := make(chan waitResult)
	commander, cancel := context.WithCancel(ctx)
	for _, key := range keys {
		go func(key interface{}) {
			e, ok := w.getEvent(key)
			if !ok {
				return
			}
			obj, err := e.WaitContextTimeout(commander, timeout)
			var r Result = nil
			if err == nil {
				r = NewResult()
				r.Set(key, obj)

				select {
				case ch <- waitResult{r, err}:
				case <-commander.Done():
				}
			}
		}(key)
	}

	defer cancel()

	select {
	case wr := <-ch:
		return wr, nil
	case <-time.After(timeout):
		return nil, ErrTimeout
	}
}

func (w *wait) WaitAll(keys ...interface{}) (Result, error) {
	return w.WaitAllContextTimeout(w.ctx, w.timeout, keys...)
}

func (w *wait) WaitAllContext(ctx context.Context, keys ...interface{}) (Result, error) {
	return w.WaitAllContextTimeout(ctx, w.timeout, keys...)
}

func (w *wait) WaitAllTimeout(timeout time.Duration, keys ...interface{}) (Result, error) {
	return w.WaitAllContextTimeout(w.ctx, timeout, keys...)
}

func (w *wait) WaitAllContextTimeout(ctx context.Context, timeout time.Duration, keys ...interface{}) (Result, error) {

	var parentContext context.Context = nil
	if ctx == nil {
		parentContext = context.Background()
	} else {
		parentContext = ctx
	}
	commander, cancel := context.WithCancel(parentContext)

	ch := make(chan waitResult)
	timer := time.NewTimer(w.timeout)
	for _, key := range keys {
		go func(key interface{}) {
			e, ok := w.getEvent(key)
			if !ok {
				return
			}

			obj, err := e.WaitContextTimeout(commander, w.timeout)
			var r Result = nil
			if err == nil {
				r = NewResult()
				r.Set(key, obj)

				select {
				case ch <- waitResult{r, err}:
				case <-commander.Done():
				}
			}
		}(key)
	}

	// 停止
	defer cancel()

	result := NewResult()
	for _, _ = range keys {
		select {
		case item := <-ch:
			item.Result.Range(func(key, value interface{}) bool {
				result.Set(key, value)
				return true
			})
		case <-timer.C:
			return nil, ErrTimeout
		}
	}

	return result, nil
}

func (w *wait) Trigger(key interface{}) {
	w.TriggerValue(key, nil)
}

func (w *wait) TriggerValue(key, value interface{}) {
	obj, ok := w.m.Load(key)
	if !ok {
		return
	}

	e := obj.(Event)
	e.Trigger(value)
}

func (w *wait) InitKey(keys ...interface{}) {
	for _, key := range keys {
		w.m.Store(key, NewEvent())
	}
}

func (w *wait) getEvent(key interface{}) (Event, bool) {
	obj, ok := w.m.Load(key)
	if !ok {
		return nil, false
	}

	return obj.(Event), true
}
