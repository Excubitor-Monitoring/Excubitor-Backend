package logging

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strings"
)

// LOG LEVELS

// LogLevel is an enum type for the available log levels.
type LogLevel int

const (
	Trace LogLevel = iota
	Debug
	Info
	Warn
	Error
	Fatal
)

func (level LogLevel) String() string {
	switch level {
	case Trace:
		return "TRACE"
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	case Fatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// GetLogLevelByString returns the LogLevel matching the given string.
func GetLogLevelByString(level string) LogLevel {
	switch strings.ToUpper(level) {
	case "TRACE":
		return Trace
	case "DEBUG":
		return Debug
	case "Info":
		return Info
	case "Warn":
		return Warn
	case "Error":
		return Error
	case "Fatal":
		return Fatal
	default:
		return Info
	}
}

// LOGGER INTERFACE

// Logger is the interface for all loggers.
// It defines the various methods to log statements.
type Logger interface {
	Trace(v ...any)
	Debug(v ...any)
	Info(v ...any)
	Warn(v ...any)
	Error(v ...any)
	Fatal(v ...any)
}

// LOGGER BUNDLE

type loggerBundle struct {
	traceLogger *log.Logger
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
}

func GetLogger() (Logger, error) {
	var logger Logger
	var err error

	loggingMethod := viper.GetString("logging.method")

	switch strings.ToUpper(loggingMethod) {
	case "CONSOLE":
		logger, err = GetConsoleLoggerInstance()
		break
	case "FILE":
		logger, err = GetFileLoggerInstance()
		break
	case "HYBRID":
		logger, err = GetMultiLoggerInstance()
		break
	default:
		fmt.Printf("Could not identify logging method %s! Falling back to console logging.\n", loggingMethod)
		logger, err = GetConsoleLoggerInstance()
		break
	}

	return logger, err
}
