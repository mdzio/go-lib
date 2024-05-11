package testutil

import (
	"os"
	"testing"

	"github.com/mdzio/go-logging"
)

// Test configuration (environment variables)
const (
	// log level, e.g. OFF, ERROR, WARNING, INFO, DEBUG, TRACE
	logLevel = "LOG_LEVEL"
)

func init() {
	var l logging.LogLevel
	err := l.Set(os.Getenv(logLevel))
	if err == nil {
		logging.SetLevel(l)
	}
}

// Config reads a test configuration. Test is skipped, if the environment
// variable with the specified name is not found.
func Config(t *testing.T, name string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		t.Skip("Environment variable " + name + " not set")
	}
	return v
}
