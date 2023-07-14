package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type MultiLogger struct {
	loggers loggerBundle
	level   LogLevel
}

func (logger *MultiLogger) Trace(v ...any) {
	if logger.level > Trace {
		return
	}

	_ = logger.loggers.traceLogger.Output(2, fmt.Sprintln(v...))
}

func (logger *MultiLogger) Debug(v ...any) {
	if logger.level > Debug {
		return
	}

	_ = logger.loggers.debugLogger.Output(2, fmt.Sprintln(v...))
}

func (logger *MultiLogger) Info(v ...any) {
	if logger.level > Info {
		return
	}

	_ = logger.loggers.infoLogger.Output(2, fmt.Sprintln(v...))
}

func (logger *MultiLogger) Warn(v ...any) {
	if logger.level > Warn {
		return
	}

	_ = logger.loggers.warnLogger.Output(2, fmt.Sprintln(v...))
}

func (logger *MultiLogger) Error(v ...any) {
	if logger.level > Error {
		return
	}

	_ = logger.loggers.errorLogger.Output(2, fmt.Sprintln(v...))
}

func (logger *MultiLogger) Fatal(v ...any) {
	_ = logger.loggers.fatalLogger.Output(2, fmt.Sprintln(v...))
}

var MultiLoggerInstance *MultiLogger

func GetMultiLoggerInstance() (*MultiLogger, error) {
	var once sync.Once
	var err error

	if MultiLoggerInstance == nil {
		once.Do(
			func() {
				var file *os.File
				file, err = os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					return
				}

				multiWriter := io.MultiWriter(file, os.Stdout)

				levelString := k.String("logging.log_level")

				logFlag := log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix
				loggers := &loggerBundle{
					traceLogger: log.New(multiWriter, fmt.Sprint("[  ", Trace, "  ] --  "), logFlag),
					debugLogger: log.New(multiWriter, fmt.Sprint("[  ", Debug, "  ] --  "), logFlag),
					infoLogger:  log.New(multiWriter, fmt.Sprint("[  ", Info, "   ] --  "), logFlag),
					warnLogger:  log.New(multiWriter, fmt.Sprint("[  ", Warn, "   ] --  "), logFlag),
					errorLogger: log.New(multiWriter, fmt.Sprint("[  ", Error, "  ] --  "), logFlag),
					fatalLogger: log.New(multiWriter, fmt.Sprint("[  ", Fatal, "  ] --  "), logFlag),
				}

				MultiLoggerInstance = &MultiLogger{loggers: *loggers, level: GetLogLevelByString(levelString)}
			})
	}

	return MultiLoggerInstance, err
}
