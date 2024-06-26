package any

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// Query helps extracting values from interface{} based models.
type Query struct {
	value interface{}
	err   *error
}

// Q creates a new Query for the specified Value.
func Q(v interface{}) *Query {
	var err error
	return &Query{value: v, err: &err}
}

// Err returns the first encountered error.
func (q *Query) Err() error {
	return *q.err
}

// SetErr sets an error from external.
func (q *Query) SetErr(err error) {
	*q.err = err
}

// Bool gets a bool.
func (q *Query) Bool() (b bool) {
	// previous error or empty?
	if q.Err() != nil || q.value == nil {
		return
	}
	// extract
	b, ok := q.value.(bool)
	if !ok {
		*q.err = errors.New("not a bool")
		return
	}
	return
}

// Float64 gets a float64.
func (q *Query) Float64() (f float64) {
	// previous error or empty?
	if q.Err() != nil || q.value == nil {
		return
	}
	// extract
	f, ok := q.value.(float64)
	if !ok {
		*q.err = errors.New("not a float64")
		return
	}
	return
}

// ToFloat64 converts the value to a float64.
func (q *Query) ToFloat64() (f float64) {
	// previous error or empty?
	if q.Err() != nil || q.value == nil {
		return
	}
	// convert value
	switch v := q.value.(type) {
	case float64:
		f = float64(v)
	case float32:
		f = float64(v)
	case int:
		f = float64(v)
	case int64:
		f = float64(v)
	case int32:
		f = float64(v)
	case int16:
		f = float64(v)
	case int8:
		f = float64(v)
	case uint:
		f = float64(v)
	case uint64:
		f = float64(v)
	case uint32:
		f = float64(v)
	case uint16:
		f = float64(v)
	case uint8:
		f = float64(v)
	case string:
		var err error
		f, err = strconv.ParseFloat(v, 64)
		if err != nil {
			*q.err = fmt.Errorf("unable to cast %#v of type %T to float64", q.value, q.value)
		}
	case json.Number:
		var err error
		f, err = v.Float64()
		if err != nil {
			*q.err = fmt.Errorf("unable to cast %#v of type %T to float64", q.value, q.value)
		}
	case bool:
		if v {
			f = 1
		} else {
			f = 0
		}
	default:
		*q.err = fmt.Errorf("unable to cast %#v of type %T to float64", q.value, q.value)
	}
	return
}

// Int gets an int.
func (q *Query) Int() (i int) {
	// previous error or empty?
	if q.Err() != nil || q.value == nil {
		return
	}
	// extract
	i, ok := q.value.(int)
	if !ok {
		*q.err = errors.New("not an int")
		return
	}
	return
}

// String gets a string.
func (q *Query) String() (s string) {
	// previous error or empty?
	if q.Err() != nil || q.value == nil {
		return
	}
	// extract
	s, ok := q.value.(string)
	if !ok {
		*q.err = errors.New("not a string")
		return
	}
	return
}

// Slice gets a []interface{} as []*Query.
func (q *Query) Slice() []*Query {
	// previous error or empty?
	if q.Err() != nil || q.value == nil {
		return nil
	}
	// extract
	s, ok := q.value.([]interface{})
	if !ok {
		*q.err = errors.New("not a slice")
		return nil
	}
	ns := make([]*Query, len(s))
	for i, v := range s {
		ns[i] = &Query{value: v, err: q.err}
	}
	return ns
}

// Map gets a map[string]interface{} as MapQuery.
func (q *Query) Map() *MapQuery {
	// previous error or empty?
	if q.Err() != nil || q.value == nil {
		return &MapQuery{err: q.err}
	}
	// extract
	m, ok := q.value.(map[string]interface{})
	if !ok {
		*q.err = errors.New("not a map")
		return &MapQuery{err: q.err}
	}
	return &MapQuery{value: m, err: q.err}
}

// Unwrap returns the wrapped value.
func (q *Query) Unwrap() interface{} {
	return q.value
}

// MapQuery helps extracting values from a map[string]interface{}.
type MapQuery struct {
	value map[string]interface{}
	err   *error
}

// Err returns the first encountered error.
func (q *MapQuery) Err() error {
	return *q.err
}

func (q *MapQuery) key(name string, must bool) *Query {
	// previous error?
	if q.Err() != nil {
		return &Query{err: q.err}
	}
	// lookup
	v, ok := q.value[name]
	if !ok {
		if must {
			*q.err = fmt.Errorf("field not found: %s", name)
		}
		return &Query{err: q.err}
	}
	return &Query{value: v, err: q.err}
}

// Key sets an error, if the specified member is missing.
func (q *MapQuery) Key(name string) *Query {
	return q.key(name, true)
}

// TryKey does not set an error, if the specified member is missing.
func (q *MapQuery) TryKey(name string) *Query {
	return q.key(name, false)
}

// Has returns true, if the the specified key exists.
func (q *MapQuery) Has(name string) bool {
	// previous error?
	if q.Err() != nil {
		return false
	}
	// lookup
	_, ok := q.value[name]
	return ok
}

// Wrap returns a new map with all values wrapped as Query.
func (q *MapQuery) Wrap() map[string]*Query {
	// previous error?
	if q.Err() != nil {
		// return empty map
		return nil
	}
	// wrap map values
	r := make(map[string]*Query)
	for k, v := range q.value {
		r[k] = &Query{value: v, err: q.err}
	}
	return r
}

// Wrap returns a new map with all values wrapped as Query.
func (q *MapQuery) Unwrap() map[string]interface{} {
	// nil on previous error
	return q.value
}
