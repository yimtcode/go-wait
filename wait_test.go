package wait

import (
	"testing"
	"time"
)

func TestWait_Wait(t *testing.T) {
	key, value := 0x1, 999
	w := NewWait()
	w.SetTimeout(10 * time.Second)
	w.InitKey(key)
	r, err := w.Wait(
		func() {
			time.AfterFunc(2*time.Second, func() {
				w.TriggerValue(key, value)
			})
		},
		key)

	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(r)
}

func TestWait_WaitAll(t *testing.T) {
	keys := []interface{}{1, 2, 3}
	values := []interface{}{11, 22, 33}

	w := NewWait()
	w.SetTimeout(10 * time.Second)
	w.InitKey(keys...)
	r, err := w.WaitAll(func() {
		time.AfterFunc(1*time.Second, func() {
			w.TriggerValue(keys[0], values[0])
		})
		time.AfterFunc(2*time.Second, func() {
			w.TriggerValue(keys[1], values[1])
		})
		time.AfterFunc(3*time.Second, func() {
			w.TriggerValue(keys[2], values[2])
		})
	}, keys...)
	if err != nil {
		t.Fatal(err.Error())
	}

	v1, _ := r.GetInt(keys[0])
	v2, _ := r.GetInt(keys[1])
	v3, _ := r.GetInt(keys[2])
	t.Logf("%d %d %d\n", v1, v2, v3)
}

func TestWait_WaitAny(t *testing.T) {
	keys := []interface{}{1, 2, 3}
	values := []interface{}{11, 22, 33}

	w := NewWait()
	w.SetTimeout(10 * time.Second)
	w.InitKey(keys...)
	r, err := w.WaitAny(func() {
		time.AfterFunc(2*time.Second, func() {
			w.TriggerValue(keys[0], values[0])
		})
		time.AfterFunc(1*time.Second, func() {
			w.TriggerValue(keys[1], values[1])
		})
		time.AfterFunc(3*time.Second, func() {
			w.TriggerValue(keys[2], values[2])
		})
	}, keys...)
	if err != nil {
		t.Fatal(err.Error())
	}

	v1, b1 := r.GetInt(keys[0])
	v2, b2 := r.GetInt(keys[1])
	v3, b3 := r.GetInt(keys[2])
	t.Logf("key1: %d,%t key2: %d,%t key3: %d,%t\n", v1, b1, v2, b2, v3, b3)
}
