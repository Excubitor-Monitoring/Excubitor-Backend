package logging

import (
	"fmt"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/configuration"
	"log"
	"os"
	"sync"
)

// CONSOLE LOGGING COLORS

type ConsoleColor string

const (
	Reset  ConsoleColor = "\033[0m"
	Red    ConsoleColor = "\033[31m"
	Green  ConsoleColor = "\033[32m"
	Yellow ConsoleColor = "\033[33m"
	Blue   ConsoleColor = "\033[34m"
	Purple ConsoleColor = "\033[35m"
	Cyan   ConsoleColor = "\033[36m"
	Gray   ConsoleColor = "\033[37m"
	White  ConsoleColor = "\033[97m"
)

type ConsoleLogger struct {
	loggers loggerBundle
	level   LogLevel
}

func (logger *ConsoleLogger) Trace(v ...any) {
	if logger.level > Trace {
		return
	}

	_ = logger.loggers.traceLogger.Output(2, fmt.Sprintln(v...))
}

func (logger *ConsoleLogger) Debug(v ...any) {
	if logger.level > Debug {
		return
	}

	_ = logger.loggers.debugLogger.Output(2, fmt.Sprintln(v...))
}

func (logger *ConsoleLogger) Info(v ...any) {
	if logger.level > Info {
		return
	}

	_ = logger.loggers.infoLogger.Output(2, fmt.Sprintln(v...))
}

func (logger *ConsoleLogger) Warn(v ...any) {
	if logger.level > Warn {
		return
	}

	_ = logger.loggers.warnLogger.Output(2, fmt.Sprintln(v...))
}

func (logger *ConsoleLogger) Error(v ...any) {
	if logger.level > Error {
		return
	}

	_ = logger.loggers.errorLogger.Output(2, fmt.Sprintln(v...))
}

func (logger *ConsoleLogger) Fatal(v ...any) {
	_ = logger.loggers.fatalLogger.Output(2, fmt.Sprintln(v...))
}

var consoleLoggerInstance *ConsoleLogger

func GetConsoleLoggerInstance() (*ConsoleLogger, error) {
	var once sync.Once
	var err error

	if consoleLoggerInstance == nil {
		once.Do(
			func() {
				var configurator *configuration.ConcreteConfigurator
				configurator, err = configuration.GetInstance()
				if err != nil {
					return
				}

				levelString := configurator.GetConfig().Logging.LogLevel

				logFlag := log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix

				loggers := &loggerBundle{
					traceLogger: log.New(os.Stdout, fmt.Sprint("[  ", Cyan, Trace, Reset, "  ] --  "), logFlag),
					debugLogger: log.New(os.Stdout, fmt.Sprint("[  ", Green, Debug, Reset, "  ] --  "), logFlag),
					infoLogger:  log.New(os.Stdout, fmt.Sprint("[  ", Reset, Info, Reset, "   ] --  "), logFlag),
					warnLogger:  log.New(os.Stdout, fmt.Sprint("[  ", Yellow, Warn, Reset, "   ] --  "), logFlag),
					errorLogger: log.New(os.Stdout, fmt.Sprint("[  ", Red, Error, Reset, "  ] --  "), logFlag),
					fatalLogger: log.New(os.Stdout, fmt.Sprint("[  ", Purple, Fatal, Reset, "  ] --  "), logFlag),
				}

				consoleLoggerInstance = &ConsoleLogger{*loggers, GetLogLevelByString(levelString)}
			})
	}

	return consoleLoggerInstance, err
}
