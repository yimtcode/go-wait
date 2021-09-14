package wait

import (
	"testing"
	"time"
)

func TestEvent(t *testing.T) {
	e := NewEvent()
	time.AfterFunc(2 * time.Second, func() {
		e.Trigger(123)
	})
	obj, err := e.WaitTimeout(3 * time.Second)
	t.Log(obj, err)
}