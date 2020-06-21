package jsonutil

import "testing"

func TestEqual(t *testing.T) {
	a := `{"a":1,"b":false,"c":"abc"}`
	b := `{"c":"abc","a":1,"b":false}`
	c := `[1,2,3]`
	d := `[3,2,1]`
	e := `[1.0,2.0,3.0]`
	cases := []struct {
		c1, c2 string
		er     bool
	}{
		{a, b, true},
		{a, c, false},
		{c, d, false},
		{c, e, true},
	}
	for _, c := range cases {
		if Equal([]byte(c.c1), []byte(c.c2)) != c.er {
			t.Error(c.c1, c.c2, c.er)
		}
	}
}
