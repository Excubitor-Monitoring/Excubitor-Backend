package logging

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/configuration"
	"log"
	"os"
	"sync"
)

type FileLogger struct {
	loggers loggerBundle
	level   LogLevel
}

func (logger *FileLogger) Trace(v ...any) {
	if logger.level > Trace {
		return
	}

	_ = logger.loggers.traceLogger.Output(2, fmt.Sprintln(v...))
}

func (logger *FileLogger) Debug(v ...any) {
	if logger.level > Debug {
		return
	}

	_ = logger.loggers.debugLogger.Output(2, fmt.Sprintln(v...))
}

func (logger *FileLogger) Info(v ...any) {
	if logger.level > Info {
		return
	}

	_ = logger.loggers.infoLogger.Output(2, fmt.Sprintln(v...))
}

func (logger *FileLogger) Warn(v ...any) {
	if logger.level > Warn {
		return
	}

	_ = logger.loggers.warnLogger.Output(2, fmt.Sprintln(v...))
}

func (logger *FileLogger) Error(v ...any) {
	if logger.level > Error {
		return
	}

	_ = logger.loggers.errorLogger.Output(2, fmt.Sprintln(v...))
}

func (logger *FileLogger) Fatal(v ...any) {
	_ = logger.loggers.fatalLogger.Output(2, fmt.Sprintln(v...))
}

var fileLoggerInstance *FileLogger

func GetFileLoggerInstance() (*FileLogger, error) {
	var once sync.Once
	var err error

	if fileLoggerInstance == nil {
		once.Do(
			func() {
				var configurator *configuration.ConcreteConfigurator
				configurator, err = configuration.GetInstance()
				if err != nil {
					return
				}

				var file *os.File
				file, err = os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					return
				}

				levelString := configurator.GetConfig().Logging.LogLevel

				logFlag := log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix
				loggers := &loggerBundle{
					traceLogger: log.New(file, fmt.Sprint("[  ", Trace, "  ] --  "), logFlag),
					debugLogger: log.New(file, fmt.Sprint("[  ", Debug, "  ] --  "), logFlag),
					infoLogger:  log.New(file, fmt.Sprint("[  ", Info, "   ] --  "), logFlag),
					warnLogger:  log.New(file, fmt.Sprint("[  ", Warn, "   ] --  "), logFlag),
					errorLogger: log.New(file, fmt.Sprint("[  ", Error, "  ] --  "), logFlag),
					fatalLogger: log.New(file, fmt.Sprint("[  ", Fatal, "  ] --  "), logFlag),
				}

				fileLoggerInstance = &FileLogger{loggers: *loggers, level: GetLogLevelByString(levelString)}
			})
	}

	return fileLoggerInstance, err
}

func compressFile(path string) error {
	input, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	defer func(writer *gzip.Writer) {
		e := writer.Close()
		err = e
		if err != nil {
			return
		}
	}(writer)

	_, err = writer.Write(input)
	if err != nil {
		return err
	}

	err = os.WriteFile("log.txt.gz", buffer.Bytes(), 0666)
	if err != nil {
		return err
	}

	return nil
}
