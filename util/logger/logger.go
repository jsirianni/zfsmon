package logger

import (
    "os"
    "io"
    "io/ioutil"
    "log"

    "github.com/pkg/errors"
)

// log levels
const (
    TraceLVL   = "trace"
    InfoLVL    = "info"
    WarningLVL = "warning"
    ErrorLVL   = "error"
)

// Logger type contains methods for logging
type Logger struct {
    trace   *log.Logger
    info    *log.Logger
    warning *log.Logger
    error   *log.Logger

    // set to true when init() has
    // been ran
    configured bool

    logLevel    string
}

// Trace prints trace logs if they are enabled
func (l Logger) Trace(message interface{}) {
    l.trace.Println(message)
}

// Info prints info logs if they are enabled
func (l Logger) Info(message interface{}) {
    l.info.Println(message)
}

// Warning prints warning logs if they are enabled
func (l Logger) Warning(message interface{}) {
    l.warning.Println(message)
}

// Error prints error logs if they are enabled
func (l Logger) Error(message interface{}) {
    l.error.Println(message)
}

// Configure configures the io.Writer values for each log level
func (l *Logger) Configure(logLevel string) error {
    switch logLevel {
	case TraceLVL:
		l.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	case InfoLVL:
		l.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	case WarningLVL:
		l.Init(ioutil.Discard, ioutil.Discard, os.Stdout, os.Stderr)
	case ErrorLVL:
		l.Init(ioutil.Discard, ioutil.Discard, ioutil.Discard, os.Stderr)
    default:
        supported := TraceLVL + ", " + InfoLVL + ", " + WarningLVL + ", " + ErrorLVL
        return errors.New("log level is not valid, supported levels are " + supported)
	}

    l.logLevel = logLevel

	return nil
}

// Init takes io.Writers for each type of log level. logger.Configure() should
// be used unless you have very specific io.Writer requirements
func (l *Logger) Init(
    traceHandle io.Writer,
    infoHandle io.Writer,
    warningHandle io.Writer,
    errorHandle io.Writer) {

    l.trace = log.New(traceHandle,
        "TRACE: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    l.info = log.New(infoHandle,
        "INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    l.warning = log.New(warningHandle,
        "WARNING: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    l.error = log.New(errorHandle,
        "ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    l.configured = true
}

// Configured returns the status of the logger
func (l Logger) Configured() bool {
    return l.configured
}

// Level returns the configured log level
func (l Logger) Level() string {
    if l.Configured() == false {
        return "not configured"
    }
    return l.logLevel
}
