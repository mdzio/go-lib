package logging

import "testing"

func TestLogLevel(t *testing.T) {
	cases := []struct {
		ident string
		err   bool
		lvl   LogLevel
		str   string
	}{
		{"OFF", false, OffLevel, "OFF"},
		{"ERROR", false, ErrorLevel, "ERROR"},
		{"warn", false, WarningLevel, "WARNING"},
		{"inF", false, InfoLevel, "INFO"},
		{"DEB", false, DebugLevel, "DEBUG"},
		{"t", false, TraceLevel, "TRACE"},
		{"warnings", true, 0, ""},
		{"", true, 0, ""},
		{"n", true, 0, ""},
	}
	for _, c := range cases {
		var l LogLevel
		err := l.Set(c.ident)
		// test parsing of level
		if (err != nil) != c.err {
			t.Errorf("invalid error: %+v", c)
		}
		if err == nil {
			// test correct level
			if l != c.lvl {
				t.Errorf("invalid level: %+v", c)
			}
			// test string of level
			if l.String() != c.str {
				t.Errorf("invalid string: %+v", c)
			}
		}
	}
}

func TestLogLevelFlag(t *testing.T) {
	pl := Level()
	var l LogLevelFlag
	err := l.Set("tr")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if Level() != TraceLevel {
		t.Errorf("unexpected level: %v", Level())
	}
	if l.String() != "TRACE" {
		t.Errorf("unexpected level string: %s", l.String())
	}
	SetLevel(pl)
}
