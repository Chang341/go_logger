package logger

import (
	"fmt"
	"os"
)

type FileLogger struct {
	level   int
	logPath string
	logName string
	appFile *os.File
	errFile *os.File
}

func NewFileLoggger(config map[string]string) (log LogInterface, err error) {
	logPath, ok := config["log_path"]
	if !ok {
		err = fmt.Errorf("not found log_path")
		return
	}
	logName, ok := config["log_name"]
	if !ok {
		err = fmt.Errorf("not found log_name")
		return
	}
	logLevel, ok := config["log_level"]
	if !ok {
		err = fmt.Errorf("not found log_level")
		return
	}

	level := getLoglevel(logLevel)
	log = &FileLogger{
		level:   level,
		logPath: logPath,
		logName: logName,
	}
	log.Init()
	return
}

func (f *FileLogger) Init() {
	// 创建appFile日志文件：存放debug、trace、info级别的日志
	fileName := fmt.Sprintf("%s/%s.app.log", f.logPath, f.logName)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed, err:%v", fileName, err))
	}
	f.appFile = file

	// 创建errFile日志文件：存放warn、error、fatal级别的日志
	fileName = fmt.Sprintf("%s/%s.error.log", f.logPath, f.logName)
	file, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed, err:%v", fileName, err))
	}
	f.errFile = file
}

func (f *FileLogger) SetLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFatal {
		level = LogLevelDebug
	}
	f.level = level
}

func (f *FileLogger) Debug(format string, args ...interface{}) {
	if f.level > LogLevelDebug {
		return
	}
	writeFile(f.appFile, LogLevelDebug, format, args...)
}

func (f *FileLogger) Trace(format string, args ...interface{}) {
	if f.level > LogLevelTrace {
		return
	}
	writeFile(f.appFile, LogLevelTrace, format, args...)
}

func (f *FileLogger) Info(format string, args ...interface{}) {
	if f.level > LogLevelInfo {
		return
	}
	writeFile(f.appFile, LogLevelInfo, format, args...)
}

func (f *FileLogger) Warn(format string, args ...interface{}) {
	if f.level > LogLevelWarn {
		return
	}
	writeFile(f.errFile, LogLevelWarn, format, args...)
}

func (f *FileLogger) Error(format string, args ...interface{}) {
	if f.level > LogLevelError {
		return
	}
	writeFile(f.errFile, LogLevelError, format, args...)
}

func (f *FileLogger) Fatal(format string, args ...interface{}) {
	if f.level > LogLevelFatal {
		return
	}
	writeFile(f.errFile, LogLevelFatal, format, args...)
}

func (f *FileLogger) Close() {
	f.appFile.Close()
	f.errFile.Close()
}
