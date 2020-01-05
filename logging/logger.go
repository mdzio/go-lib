package logging

// A Logger is used to generate log messages.
type Logger interface {

	// Will a specific log level output a log message to the Writer? These
	// functions can be used to avoid costly builds of log messages.
	ErrorEnabled() bool
	WarningEnabled() bool
	InfoEnabled() bool
	DebugEnabled() bool
	TraceEnabled() bool

	// Create a log message using fmt.Sprint for formatting.
	Error(values ...interface{})
	Warning(values ...interface{})
	Info(values ...interface{})
	Debug(values ...interface{})
	Trace(values ...interface{})

	// Create a log message using fmt.Sprintf for formatting.
	Errorf(format string, values ...interface{})
	Warningf(format string, values ...interface{})
	Infof(format string, values ...interface{})
	Debugf(format string, values ...interface{})
	Tracef(format string, values ...interface{})
}

type logger string

// Get creates an Logger with the specified identifier.
func Get(id string) Logger {
	l := logger(id)
	return &l
}

func (l *logger) ErrorEnabled() bool {
	return Level() >= ErrorLevel
}

func (l *logger) WarningEnabled() bool {
	return Level() >= WarningLevel
}

func (l *logger) InfoEnabled() bool {
	return Level() >= InfoLevel
}

func (l *logger) DebugEnabled() bool {
	return Level() >= DebugLevel
}

func (l *logger) TraceEnabled() bool {
	return Level() >= TraceLevel
}

func (l *logger) Error(values ...interface{}) {
	if !l.ErrorEnabled() {
		return
	}
	log(ErrorLevel, string(*l), values...)
}

func (l *logger) Warning(values ...interface{}) {
	if !l.WarningEnabled() {
		return
	}
	log(WarningLevel, string(*l), values...)
}

func (l *logger) Info(values ...interface{}) {
	if !l.InfoEnabled() {
		return
	}
	log(InfoLevel, string(*l), values...)
}

func (l *logger) Debug(values ...interface{}) {
	if !l.DebugEnabled() {
		return
	}
	log(DebugLevel, string(*l), values...)
}

func (l *logger) Trace(values ...interface{}) {
	if !l.TraceEnabled() {
		return
	}
	log(TraceLevel, string(*l), values...)
}

func (l *logger) Errorf(format string, values ...interface{}) {
	if !l.ErrorEnabled() {
		return
	}
	logf(ErrorLevel, string(*l), format, values...)
}

func (l *logger) Warningf(format string, values ...interface{}) {
	if !l.WarningEnabled() {
		return
	}
	logf(WarningLevel, string(*l), format, values...)
}

func (l *logger) Infof(format string, values ...interface{}) {
	if !l.InfoEnabled() {
		return
	}
	logf(InfoLevel, string(*l), format, values...)
}

func (l *logger) Debugf(format string, values ...interface{}) {
	if !l.DebugEnabled() {
		return
	}
	logf(DebugLevel, string(*l), format, values...)
}

func (l *logger) Tracef(format string, values ...interface{}) {
	if !l.TraceEnabled() {
		return
	}
	logf(TraceLevel, string(*l), format, values...)
}
