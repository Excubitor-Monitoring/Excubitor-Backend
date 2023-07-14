package plugins

import (
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/hashicorp/go-hclog"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type LogWrapper struct {
	logger         logging.Logger
	name           string
	persistentArgs []interface{}
	level          hclog.Level
}

func NewLogWrapper(logger logging.Logger, name string, level logging.LogLevel) *LogWrapper {
	var hclogLevel hclog.Level

	switch level {
	case logging.Trace:
		hclogLevel = hclog.Trace
	case logging.Debug:
		hclogLevel = hclog.Debug
	case logging.Info:
		hclogLevel = hclog.Info
	case logging.Warn:
		hclogLevel = hclog.Warn
	case logging.Error:
		hclogLevel = hclog.Error
	case logging.Fatal:
		hclogLevel = hclog.Error
	}

	return &LogWrapper{
		logger: logger,
		name:   name,
		level:  hclogLevel,
	}
}

func (w *LogWrapper) Log(level hclog.Level, msg string, args ...interface{}) {
	switch level {
	case hclog.Error:
		w.Error(msg, args...)
	case hclog.Warn:
		w.Warn(msg, args...)
	case hclog.Info:
		w.Info(msg, args...)
	case hclog.Debug:
		w.Debug(msg, args...)
	case hclog.Trace:
		w.Trace(msg, args...)
	default:
		w.Info("[UNKNOWN LOG LEVEL] "+msg, args)
	}
}

func (w *LogWrapper) Trace(msg string, args ...interface{}) {
	if w.level > hclog.Trace {
		return
	}

	w.logger.Trace(w.formatMessage(msg, args))
}

func (w *LogWrapper) Debug(msg string, args ...interface{}) {
	if w.level > hclog.Debug {
		return
	}

	w.logger.Debug(w.formatMessage(msg, args))
}

func (w *LogWrapper) Info(msg string, args ...interface{}) {
	if w.level > hclog.Info {
		return
	}

	w.logger.Info(w.formatMessage(msg, args))
}

func (w *LogWrapper) Warn(msg string, args ...interface{}) {
	if w.level > hclog.Warn {
		return
	}

	w.logger.Warn(w.formatMessage(msg, args))
}

func (w *LogWrapper) Error(msg string, args ...interface{}) {
	if w.level > hclog.Error {
		return
	}

	w.logger.Error(w.formatMessage(msg, args))
}

func (w *LogWrapper) IsTrace() bool {
	return w.level <= hclog.Trace
}

func (w *LogWrapper) IsDebug() bool {
	return w.level <= hclog.Debug
}

func (w *LogWrapper) IsInfo() bool {
	return w.level <= hclog.Info
}

func (w *LogWrapper) IsWarn() bool {
	return w.level <= hclog.Warn
}

func (w *LogWrapper) IsError() bool {
	return w.level <= hclog.Error
}

func (w *LogWrapper) ImpliedArgs() []interface{} {
	return w.persistentArgs
}

func (w *LogWrapper) With(args ...interface{}) hclog.Logger {
	return &LogWrapper{logger: w.logger, name: w.name, persistentArgs: args}
}

func (w *LogWrapper) Name() string {
	return w.name
}

func (w *LogWrapper) Named(name string) hclog.Logger {
	if w.name == "" {
		return w.ResetNamed(name)
	} else {
		return &LogWrapper{logger: w.logger, name: w.name + " > " + name, persistentArgs: w.persistentArgs}
	}
}

func (w *LogWrapper) ResetNamed(name string) hclog.Logger {
	return &LogWrapper{logger: w.logger, name: name, persistentArgs: w.persistentArgs}
}

func (w *LogWrapper) SetLevel(level hclog.Level) {
	w.level = level
}

func (w *LogWrapper) StandardLogger(_ *hclog.StandardLoggerOptions) *log.Logger {
	return log.New(os.Stdout, "StandardLogger", 0)
}

func (w *LogWrapper) StandardWriter(_ *hclog.StandardLoggerOptions) io.Writer {
	return os.Stdout
}

func (w *LogWrapper) formatMessage(msg string, args []interface{}) string {
	var output strings.Builder

	if w.name != "" {
		output.WriteString("[ " + w.name + " ] - ")
	}

	output.WriteString(msg)

	args = append(w.persistentArgs, args...)

	if len(args) == 1 && reflect.TypeOf(args[0]).String() != "string" {
		return msg
	}

	opened := false
	skip := false

	if len(args) > 0 {
		for index, arg := range args {
			if skip {
				skip = false
				continue
			}

			if index%2 == 0 {
				switch arg.(type) {
				case string:
					if !opened {
						output.WriteString(" (")
						opened = true
					}
					output.WriteString(arg.(string) + " = ")
				default:
					skip = true
					continue
				}
			} else {
				switch arg.(type) {
				case string:
					output.WriteString("\"" + arg.(string) + "\"")
					if index != len(args)-1 {
						output.WriteString(", ")
					}
				case []string:
					output.WriteString("[")

					for contentIndex, content := range arg.([]string) {
						output.WriteString("\"" + content + "\"")

						if contentIndex != len(arg.([]string))-1 {
							output.WriteString(", ")
						}
					}

					output.WriteString("]")
				case int:
					output.WriteString(strconv.Itoa(arg.(int)))
				}

				if index == len(args)-1 {
					output.WriteString(")")
				}
			}
		}
	}

	return output.String()
}
