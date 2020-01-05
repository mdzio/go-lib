package any

import "testing"

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
