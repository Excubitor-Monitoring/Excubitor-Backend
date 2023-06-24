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
	w.logger.Trace(w.formatMessage(msg, args))
}

func (w *LogWrapper) Debug(msg string, args ...interface{}) {
	w.logger.Debug(w.formatMessage(msg, args))
}

func (w *LogWrapper) Info(msg string, args ...interface{}) {
	w.logger.Info(w.formatMessage(msg, args))
}

func (w *LogWrapper) Warn(msg string, args ...interface{}) {
	w.logger.Warn(w.formatMessage(msg, args))
}

func (w *LogWrapper) Error(msg string, args ...interface{}) {
	w.logger.Error(w.formatMessage(msg, args))
}

func (w *LogWrapper) IsTrace() bool {
	return true
}

func (w *LogWrapper) IsDebug() bool {
	return true
}

func (w *LogWrapper) IsInfo() bool {
	return true
}

func (w *LogWrapper) IsWarn() bool {
	return true
}

func (w *LogWrapper) IsError() bool {
	return true
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
	return
}

func (w *LogWrapper) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	return log.New(os.Stdout, "StandardLogger", 0)
}

func (w *LogWrapper) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
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
