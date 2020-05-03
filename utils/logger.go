package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

type LogLevel string

// Keys for log level
const (
	logLevelInfo  LogLevel = "INFO"
	logLevelWarn  LogLevel = "WARN"
	logLevelError LogLevel = "ERROR"
	logLevelDebug LogLevel = "DEBUG"
)

// Log instance
var Log = Logger()

// Logger creates a new logger instance
func Logger() *logInstance {
	return &logInstance{}
}

// LogInstance model of logging instance
type logInstance struct{}

// Warn prints a warning log message to stdout.
func (log *logInstance) Warn(object interface{}) {
	log.writeToLog(logLevelWarn, object)
}

// Warnf prints a warning log message to stdout using any further arguments to format the message.
// See fmt.Sprintf.
func (log *logInstance) Warnf(message string, args ...interface{}) {
	log.Warn(fmt.Sprintf(message, args...))
}

// Info prints a info log message to stdout.
func (log *logInstance) Info(object interface{}) {
	log.writeToLog(logLevelInfo, object)
}

// Infof prints a info log message to stdout using any further arguments to format the message.
// See fmt.Sprintf.
func (log *logInstance) Infof(message string, args ...interface{}) {
	log.Info(fmt.Sprintf(message, args...))
}

// Debug prints a debug log message to stdout.
func (log *logInstance) Debug(object interface{}) {
	log.writeToLog(logLevelDebug, object)
}

// Debugf prints a debug log message to stdout using any further arguments to format the message.
// See fmt.Sprintf.
func (log *logInstance) Debugf(message string, args ...interface{}) {
	log.Debug(fmt.Sprintf(message, args...))
}

// Error prints a error log message to stdout.
func (log *logInstance) Error(object interface{}) {
	log.writeToLog(logLevelError, object)
}

// Errorf prints a error log message to stdout using any further arguments to format the message.
// See fmt.Sprintf.
func (log *logInstance) Errorf(message string, args ...interface{}) {
	log.Error(fmt.Sprintf(message, args...))
}

// Fatal prints a error message to stdout and exits the main process.
func (log *logInstance) Fatal(object interface{}) {
	log.Error(object)
	os.Exit(1)
}

// Fatalf prints a error message to stdout using any further arguments to format the message and exits the main process.
// See fmt.Sprintf.
func (log *logInstance) Fatalf(message string, args ...interface{}) {
	log.Errorf(message, args...)
	os.Exit(1)
}

func (log *logInstance) writeToLog(loglevel LogLevel, object interface{}) {
	var logmap map[string]interface{}

	switch object.(type) {
	case map[string]interface{}:
		logmap = object.(map[string]interface{})
		break
	case string:
		logmap = map[string]interface{}{}
		logmap["text"] = object.(string)
		break
	default:
		logmap = map[string]interface{}{}
		logmap["reference"] = object
		break
	}
	logmap["level"] = loglevel

	logJSON, err := json.Marshal(logmap)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(logJSON))
	}
}
