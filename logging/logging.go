/*
Package logging provides a logging API and an implementation to output log
messages to stderr or any other Writer. The complete logging package can be used
concurrently from multiple goroutines. The performance costs for handling
suppressed log messages are minimal (one non-blocking atomic integer
comparison).

Each part of the application creates its own Logger in the global var section or
within the init function:
    var (
        log := logging.Get("my_id")
    )

For each log statement a level is set by the invoked log function. The log
levels are ordered as follows: error (highest severity) < warning < info < debug
< trace (lowest severity).

The Logger provides various functions to log and format messages at different
levels. The log messages are formatted through fmt.Sprint and fmt.Sprintf:
    log.Error("my error message: ", myValue)
    log.Debugf("my debug message: %v", myValue)
*/
package logging

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// LogFlags are used to modify the output format of log messages.
type LogFlags int32

// Defined flags for the output format.
const (
	// Includes a timestamp in the log line.
	TimeFlag LogFlags = 1 << iota
	// Includes the log level.
	LevelFlag
	// Includes the identifer.
	IdentifierFlag
)

const (
	timeFmt = "2006-01-02 15:04:05"
)

var (
	// use atomic access for level and flags
	level = int32(InfoLevel)
	flags = int32(TimeFlag | LevelFlag | IdentifierFlag)

	// use the mutex for the Writer
	mutex  sync.Mutex
	writer io.Writer = os.Stderr
)

// SetLevel sets the required log level of the messages, which should be logged.
func SetLevel(lvl LogLevel) {
	atomic.StoreInt32(&level, int32(lvl))
}

// Level gets the current log level.
func Level() LogLevel {
	return LogLevel(atomic.LoadInt32(&level))
}

// SetFlags sets flags for the output format.
func SetFlags(f LogFlags) {
	atomic.StoreInt32(&flags, int32(f))
}

// Flags gets the current flags.
func Flags() LogFlags {
	return LogFlags(atomic.LoadInt32(&flags))
}

// SetWriter sets a new Writer for outputting the log messages.
func SetWriter(w io.Writer) {
	mutex.Lock()
	defer mutex.Unlock()
	writer = w
}

func logf(lvl LogLevel, ident string, format string, values ...interface{}) {
	logstr(lvl, ident, fmt.Sprintf(format, values...))
}

func log(lvl LogLevel, ident string, values ...interface{}) {
	logstr(lvl, ident, fmt.Sprint(values...))
}

func logstr(lvl LogLevel, ident string, msg string) {
	var b strings.Builder
	sep := false
	f := Flags()

	if (f & TimeFlag) != 0 {
		b.WriteString(time.Now().Format(timeFmt))
		sep = true
	}
	if (f & LevelFlag) != 0 {
		if sep {
			b.WriteRune('|')
		}
		fmt.Fprintf(&b, "%-7s", lvl.String())
		sep = true
	}
	if (f & IdentifierFlag) != 0 {
		if sep {
			b.WriteRune('|')
		}
		fmt.Fprintf(&b, "%-15s", ident)
		sep = true
	}
	if sep {
		b.WriteRune('|')
	}

	// replace special chars
	r := strings.NewReplacer("\n", `\n`, "\r", ``)
	b.WriteString(r.Replace(msg))

	b.WriteRune('\n')

	mutex.Lock()
	defer mutex.Unlock()
	writer.Write([]byte(b.String()))
}
