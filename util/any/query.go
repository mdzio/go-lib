package any

import (
	"errors"
	"fmt"
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
