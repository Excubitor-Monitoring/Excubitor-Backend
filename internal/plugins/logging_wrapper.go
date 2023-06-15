package plugins

import (
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/hashicorp/go-hclog"
	"io"
	"log"
	"os"
)

type LogWrapper struct {
	logger logging.Logger
}

func (w *LogWrapper) Log(level hclog.Level, msg string, args ...interface{}) {
	switch level {
	case hclog.Error:
		w.Error(msg, args)
	case hclog.Warn:
		w.Warn(msg, args)
	case hclog.Info:
		w.Info(msg, args)
	case hclog.Debug:
		w.Debug(msg, args)
	case hclog.Trace:
		w.Trace(msg, args)
	default:
		w.Info("[UNKNOWN LOG LEVEL] "+msg, args)
	}
}

func (w *LogWrapper) Trace(msg string, args ...interface{}) {
	w.logger.Trace(msg)
}

func (w *LogWrapper) Debug(msg string, args ...interface{}) {
	w.logger.Debug(msg)
}

func (w *LogWrapper) Info(msg string, args ...interface{}) {
	w.logger.Info(msg)
}

func (w *LogWrapper) Warn(msg string, args ...interface{}) {
	w.logger.Warn(msg)
}

func (w *LogWrapper) Error(msg string, args ...interface{}) {
	w.logger.Debug(msg)
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
	return nil
}

func (w *LogWrapper) With(args ...interface{}) hclog.Logger {
	return w
}

func (w *LogWrapper) Name() string {
	return "WrapperLogger"
}

func (w *LogWrapper) Named(name string) hclog.Logger {
	return w
}

func (w *LogWrapper) ResetNamed(name string) hclog.Logger {
	return w
}

func (w *LogWrapper) SetLevel(level hclog.Level) {
	return
}

func (w *LogWrapper) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	return log.New(os.Stdout, "StandardLogger", log.LstdFlags)
}

func (w *LogWrapper) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return os.Stdout
}
