package logging

import (
	"errors"
	"fmt"
	"strings"
)

// LogLevel describes the severity of log messages.
type LogLevel int32

// Defined log levels.
const (
	OffLevel LogLevel = iota
	ErrorLevel
	WarningLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

var (
	levelStr = []string{
		OffLevel:     "OFF",
		ErrorLevel:   "ERROR",
		WarningLevel: "WARNING",
		InfoLevel:    "INFO",
		DebugLevel:   "DEBUG",
		TraceLevel:   "TRACE",
	}

	errInvalidLevelIdent = errors.New("invalid log level identifier (expected: off, error, warning, info, debug or trace)")
)

// String returns a string representation of the stored level. This function is
// part of the flag.Value interface.
func (l *LogLevel) String() string {
	return levelStr[*l]
}

// Set updates the log level. Valid identifiers are: off, error, warning, info,
// debug, trace. Upper/lower case is ignored. Identifiers can be shortened (e.g.
// err, i). This function is part of the flag.Value interface.
func (l *LogLevel) Set(value string) error {
	if len(value) == 0 {
		return errInvalidLevelIdent
	}
	value = strings.ToUpper(value)
	for idx, lvlstr := range levelStr {
		if strings.HasPrefix(lvlstr, value) {
			*l = LogLevel(idx)
			return nil
		}
	}
	return errInvalidLevelIdent
}

// LogLevelFlag can be used with flag.Var and sets/gets the active level.
type LogLevelFlag struct{}

// String returns a string representation of the active level. This function is
// part of the flag.Value interface.
func (*LogLevelFlag) String() string {
	// read global log level
	l := Level()
	return l.String()
}

// Set updates the active log level. Valid identifiers are: off, error, warning,
// info, debug, trace. Upper/lower case is ignored. Identifiers can be shortened
// (e.g. err, i). This function is part of the flag.Value interface.
func (*LogLevelFlag) Set(value string) error {
	var l LogLevel
	if err := l.Set(value); err != nil {
		return err
	}
	SetLevel(l)
	return nil
}

func test() {
	var v *LogLevelFlag
	fmt.Println(v.String())
}
