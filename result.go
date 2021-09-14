package wait

import "sync"

// Result key must sync.Map support
type Result interface {
	Set(key, value interface{})
	Get(key interface{}) (interface{}, bool)

	GetInt(key interface{}) (int, bool)
	GetInt8(key interface{}) (int8, bool)
	GetInt16(key interface{}) (int16, bool)
	GetInt32(key interface{}) (int32, bool)
	GetInt64(key interface{}) (int64, bool)

	GetUint(key interface{}) (uint, bool)
	GetUint8(key interface{}) (uint8, bool)
	GetUint16(key interface{}) (uint16, bool)
	GetUint32(key interface{}) (uint32, bool)
	GetUint64(key interface{}) (uint64, bool)

	GetFloat32(key interface{}) (float32, bool)
	GetFloat64(key interface{}) (float64, bool)

	GetString(key interface{}) (string, bool)
	GetBool(key interface{}) (bool, bool)

	Range(f func(key, value interface{}) bool)
}

func NewResult() Result {
	r := &result{
		m: sync.Map{},
	}
	return r
}

type result struct {
	m sync.Map
}

func (r *result) GetInt(key interface{}) (int, bool) {
	obj, ok := r.Get(key)
	if !ok {
		return 0, false
	}

	return obj.(int), true
}

func (r *result) GetInt8(key interface{}) (int8, bool) {
	obj, ok := r.Get(key)
	if !ok {
		return 0, false
	}

	return obj.(int8), true
}

func (r *result) GetInt16(key interface{}) (int16, bool) {
	obj, ok := r.Get(key)
	if !ok {
		return 0, false
	}

	return obj.(int16), true
}

func (r *result) GetInt32(key interface{}) (int32, bool) {
	obj, ok := r.Get(key)
	if !ok {
		return 0, false
	}

	return obj.(int32), true
}

func (r *result) GetInt64(key interface{}) (int64, bool) {
	obj, ok := r.Get(key)
	if !ok {
		return 0, false
	}

	return obj.(int64), true
}

func (r *result) GetUint(key interface{}) (uint, bool) {
	obj, ok := r.Get(key)
	if !ok {
		return 0, false
	}

	return obj.(uint), true
}

func (r *result) GetUint8(key interface{}) (uint8, bool) {
	obj, ok := r.Get(key)
	if !ok {
		return 0, false
	}

	return obj.(uint8), true
}

func (r *result) GetUint16(key interface{}) (uint16, bool) {
	obj, ok := r.Get(key)
	if !ok {
		return 0, false
	}

	return obj.(uint16), true
}

func (r *result) GetUint32(key interface{}) (uint32, bool) {
	obj, ok := r.Get(key)
	if !ok {
		return 0, false
	}

	return obj.(uint32), true
}

func (r *result) GetUint64(key interface{}) (uint64, bool) {
	obj, ok := r.Get(key)
	if !ok {
		return 0, false
	}

	return obj.(uint64), true
}

func (r *result) GetFloat32(key interface{}) (float32, bool) {
	obj, ok := r.Get(key)
	if !ok {
		return 0, false
	}

	return obj.(float32), true
}

func (r *result) GetFloat64(key interface{}) (float64, bool) {
	obj, ok := r.Get(key)
	if !ok {
		return 0, false
	}

	return obj.(float64), true
}

func (r *result) GetString(key interface{}) (string, bool) {
	obj, ok := r.Get(key)
	if !ok {
		return "", false
	}

	return obj.(string), true
}

func (r *result) GetBool(key interface{}) (bool, bool) {
	obj, ok := r.Get(key)
	if !ok {
		return false, false
	}

	return obj.(bool), true
}

func (r *result) Set(key, value interface{}) {
	r.m.Store(key, value)
}

func (r *result) Get(key interface{}) (interface{}, bool) {
	return r.m.Load(key)
}

func (r *result) Range(f func(key, value interface{}) bool) {
	r.m.Range(f)
}
