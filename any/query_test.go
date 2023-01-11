package any

import (
	"encoding/json"
	"testing"
)

func TestQuery(t *testing.T) {
	var v interface{} = 123.456
	q := Q(v)
	if q.Err() != nil {
		t.Error(q.Err())
	}
	if q.Float64() != 123.456 {
		t.Error(q.Float64())
	}
	if q.Err() != nil {
		t.Error(q.Err())
	}
	v = "abc"
	q = Q(v)
	if q.String() != "abc" {
		t.Error(q.String())
	}
	if q.Err() != nil {
		t.Error(q.Err())
	}
	if q.Float64() != 0 {
		t.Error(q.Float64())
	}
	if q.Err() == nil {
		t.Error("expected error")
	}
	_ = q.String()
	if q.Err() == nil {
		t.Error("expected error")
	}
}

func TestMapQuery(t *testing.T) {
	var v interface{} = map[string]interface{}{"a": 123.456, "b": "abc"}
	q := Q(v)
	m := q.Map()
	if m.Err() != nil {
		t.Error(m.Err())
	}
	str := m.TryKey("c").String()
	if q.Err() != nil {
		t.Error(q.Err())
	}
	if str != "" {
		t.Error(str)
	}
	f := m.Key("a").Float64()
	if q.Err() != nil {
		t.Error(q.Err())
	}
	if f != 123.456 {
		t.Error(f)
	}
	b := m.Key("c").Bool()
	if q.Err() == nil {
		t.Error("expected error")
	}
	if b == true {
		t.Error(b)
	}
}

func TestMapWrap(t *testing.T) {
	var v interface{} = map[string]interface{}{"c": 42, "d": true}
	q := Q(v)
	m := q.Map().Wrap()
	e1, ok := m["c"]
	if !ok || e1.Int() != 42 {
		t.Error(e1)
	}
	e2, ok := m["d"]
	if !ok {
		t.Error("missing entry")
	}
	if q.Err() != nil {
		t.Error(q.Err())
	}
	if !e2.Bool() {
		t.Error("expected true")
	}
	if e2.Int() != 0 || q.Err() == nil {
		t.Error("expected error")
	}
}

func TestSliceQuery(t *testing.T) {
	var v interface{} = []interface{}{"a", 123.456, "b", true}
	q := Q(v)
	sl := q.Slice()
	if q.Err() != nil {
		t.Error(q.Err())
	}
	if len(sl) != 4 {
		t.Error(len(sl))
	}
	str := sl[0].String()
	if q.Err() != nil {
		t.Error(q.Err())
	}
	if str != "a" {
		t.Error(str)
	}
	b := sl[3].Bool()
	if q.Err() != nil {
		t.Error(q.Err())
	}
	if b != true {
		t.Error(b)
	}
	str = sl[1].String()
	if q.Err() == nil {
		t.Error("expected error")
	}
	if str != "" {
		t.Error(str)
	}
}

func TestQueryToFloat64(t *testing.T) {
	cases := []struct {
		val interface{}
		exp float64
		err string
	}{
		{float64(123.5), 123.5, ""},
		{int(-42), -42, ""},
		{"123.5", 123.5, ""},
		{"123.5 ", 0, `unable to cast "123.5 " of type string to float64`},
		{json.Number("33.5"), 33.5, ""},
		{json.Number(" 33.5"), 33.5, `unable to cast " 33.5" of type json.Number to float64`},
		{true, 1, ""},
		{nil, 0, ""},
	}
	for _, c := range cases {
		q := Q(c.val)
		act := q.ToFloat64()
		if q.Err() != nil {
			if c.err == "" {
				t.Errorf("case %#v: unexpected error: %v", c.val, q.Err())
			} else if q.Err().Error() != c.err {
				t.Errorf("case %#v: wrong error: %v", c.val, q.Err())
			}
		} else {
			if c.err != "" {
				t.Errorf("case %#v: expected error: %s", c.val, c.err)
			} else if act != c.exp {
				t.Errorf("case %#v: wrong value: %f", c.val, act)
			}
		}
	}
}
