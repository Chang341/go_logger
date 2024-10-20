package go_logger

import (
	"fmt"
	"os"
)

type ConsoleLogger struct {
	level int
}

func (c *ConsoleLogger) Init() {
}

func NewConsoleLoggger(config map[string]string) (log LogInterface, err error) {
	logLevel, ok := config["log_level"]
	if !ok {
		err = fmt.Errorf("not found log_level")
		return
	}
	level := getLoglevel(logLevel)
	log = &ConsoleLogger{
		level: level,
	}
	return
}

func (c *ConsoleLogger) SetLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFatal {
		level = LogLevelDebug
	}
	c.level = level
}

func (c *ConsoleLogger) Debug(format string, args ...interface{}) {
	if c.level > LogLevelDebug {
		return
	}
	logData := writeFile(LogLevelDebug, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s [%s:%s:%d] %s\n",
		logData.TimeStr, logData.LevelStr,
		logData.FileName, logData.FuncName, logData.LineNo, logData.Message)
}

func (c *ConsoleLogger) Trace(format string, args ...interface{}) {
	if c.level > LogLevelTrace {
		return
	}
	logData := writeFile(LogLevelTrace, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s [%s:%s:%d] %s\n",
		logData.TimeStr, logData.LevelStr,
		logData.FileName, logData.FuncName, logData.LineNo, logData.Message)
}

func (c *ConsoleLogger) Info(format string, args ...interface{}) {
	if c.level > LogLevelInfo {
		return
	}
	//writeFile(os.Stdout, LogLevelInfo, format, args...)
	logData := writeFile(LogLevelInfo, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s [%s:%s:%d] %s\n",
		logData.TimeStr, logData.LevelStr,
		logData.FileName, logData.FuncName, logData.LineNo, logData.Message)
}

func (c *ConsoleLogger) Warn(format string, args ...interface{}) {
	if c.level > LogLevelWarn {
		return
	}
	//writeFile(os.Stdout, LogLevelWarn, format, args...)
	logData := writeFile(LogLevelWarn, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s [%s:%s:%d] %s\n",
		logData.TimeStr, logData.LevelStr,
		logData.FileName, logData.FuncName, logData.LineNo, logData.Message)
}

func (c *ConsoleLogger) Error(format string, args ...interface{}) {
	if c.level > LogLevelError {
		return
	}
	//writeFile(os.Stdout, LogLevelError, format, args...)
	logData := writeFile(LogLevelError, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s [%s:%s:%d] %s\n",
		logData.TimeStr, logData.LevelStr,
		logData.FileName, logData.FuncName, logData.LineNo, logData.Message)
}

func (c *ConsoleLogger) Fatal(format string, args ...interface{}) {
	if c.level > LogLevelFatal {
		return
	}
	//writeFile(os.Stdout, LogLevelFatal, format, args...)
	logData := writeFile(LogLevelFatal, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s [%s:%s:%d] %s\n",
		logData.TimeStr, logData.LevelStr,
		logData.FileName, logData.FuncName, logData.LineNo, logData.Message)
}

func (c *ConsoleLogger) Close() {

}
