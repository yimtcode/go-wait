package wait

import (
	"context"
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
	Wait(work func(), key interface{}) (interface{}, error)
	WaitContext(work func(), ctx context.Context, key interface{}) (interface{}, error)
	WaitTimeout(work func(), timeout time.Duration, key interface{}) (interface{}, error)
	WaitContextTimeout(work func(), ctx context.Context, timeout time.Duration, key interface{}) (interface{}, error)

	WaitAny(work func(), keys ...interface{}) (Result, error)
	WaitAnyContext(work func(), ctx context.Context, keys ...interface{}) (Result, error)
	WaitAnyTimeout(work func(), timeout time.Duration, keys ...interface{}) (Result, error)
	WaitAnyContextTimeout(work func(), ctx context.Context, timeout time.Duration, keys ...interface{}) (Result, error)

	WaitAll(work func(), keys ...interface{}) (Result, error)
	WaitAllContext(work func(), ctx context.Context, keys ...interface{}) (Result, error)
	WaitAllTimeout(work func(), timeout time.Duration, keys ...interface{}) (Result, error)
	WaitAllContextTimeout(work func(), ctx context.Context, timeout time.Duration, keys ...interface{}) (Result, error)

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

func (w *wait) Wait(work func(), key interface{}) (interface{}, error) {
	return w.WaitContextTimeout(work, w.ctx, w.timeout, key)
}

func (w *wait) WaitContext(work func(), ctx context.Context, key interface{}) (interface{}, error) {
	return w.WaitContextTimeout(work, ctx, w.timeout, key)
}

func (w *wait) WaitTimeout(work func(), timeout time.Duration, key interface{}) (interface{}, error) {
	return w.WaitContextTimeout(work, w.ctx, timeout, key)
}

func (w *wait) WaitContextTimeout(work func(), ctx context.Context, timeout time.Duration, key interface{}) (interface{}, error) {
	e, _ := w.getEvent(key)

	go work()

	return e.WaitContextTimeout(ctx, timeout)
}

func (w *wait) WaitAny(work func(), keys ...interface{}) (Result, error) {
	return w.WaitAnyContextTimeout(work, w.ctx, w.timeout, keys...)
}

func (w *wait) WaitAnyContext(work func(), ctx context.Context, keys ...interface{}) (Result, error) {
	return w.WaitAnyContextTimeout(work, ctx, w.timeout, keys...)
}

func (w *wait) WaitAnyTimeout(work func(), timeout time.Duration, keys ...interface{}) (Result, error) {
	return w.WaitAnyContextTimeout(work, w.ctx, timeout, keys...)
}

func (w *wait) WaitAnyContextTimeout(work func(), ctx context.Context, timeout time.Duration, keys ...interface{}) (Result, error) {

	parentContext := ctx
	if parentContext == nil {
		parentContext = context.Background()
	}
	ch := make(chan waitResult)
	commander, cancel := context.WithCancel(parentContext)
	for _, key := range keys {
		go func(key interface{}) {
			e, _ := w.getEvent(key)
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

	go work()

	defer cancel()

	select {
	case wr := <-ch:
		return wr, nil
	case <-time.After(timeout):
		return nil, ErrTimeout
	}
}

func (w *wait) WaitAll(work func(), keys ...interface{}) (Result, error) {
	return w.WaitAllContextTimeout(work, w.ctx, w.timeout, keys...)
}

func (w *wait) WaitAllContext(work func(), ctx context.Context, keys ...interface{}) (Result, error) {
	return w.WaitAllContextTimeout(work, ctx, w.timeout, keys...)
}

func (w *wait) WaitAllTimeout(work func(), timeout time.Duration, keys ...interface{}) (Result, error) {
	return w.WaitAllContextTimeout(work, w.ctx, timeout, keys...)
}

func (w *wait) WaitAllContextTimeout(work func(), ctx context.Context, timeout time.Duration, keys ...interface{}) (Result, error) {

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
			e, _ := w.getEvent(key)
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

	go work()

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
