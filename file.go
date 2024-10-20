package go_logger

import (
	"fmt"
	"os"
	"strconv"
)

type FileLogger struct {
	level       int
	logPath     string
	logName     string
	appFile     *os.File
	errFile     *os.File
	logDataChal chan *LogData // 声明管道/队列，存放LogData日志数据的指针
	IsErrLog    bool          // 标识写到哪个log文件中
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
	chanSizeStr, ok := config["log_chan_size"]
	if !ok {
		chanSizeStr = "50000"
	}
	chanSize, err := strconv.Atoi(chanSizeStr)
	if err != nil {
		chanSize = 50000
	}
	level := getLoglevel(logLevel)
	isErrLog := false
	if level == LogLevelWarn || level == LogLevelFatal || level == LogLevelError {
		isErrLog = true
	}
	log = &FileLogger{
		level:       level,
		logPath:     logPath,
		logName:     logName,
		logDataChal: make(chan *LogData, chanSize),
		IsErrLog:    isErrLog,
	}
	log.Init()
	return
}

func (f *FileLogger) Init() {
	// 创建appFile日志文件：存放debug、trace、info级别的日志
	fileName := fmt.Sprintf(AppLogFormat, f.logPath, f.logName)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed, err:%v", fileName, err))
	}
	f.appFile = file

	// 创建errFile日志文件：存放warn、error、fatal级别的日志
	fileName = fmt.Sprintf(ErrLogFormat, f.logPath, f.logName)
	file, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed, err:%v", fileName, err))
	}
	f.errFile = file
	go f.writeLogAsync()
}

func (f *FileLogger) writeLogAsync() {
	for logData := range f.logDataChal { // 遍历队列，取出数据完成写入
		file := f.appFile
		if f.IsErrLog {
			file = f.errFile
		}
		fmt.Fprintf(file, LogDataFormat,
			logData.TimeStr, logData.LevelStr,
			logData.FileName, logData.FuncName, logData.LineNo, logData.Message)
	}
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
	//writeFile(f.appFile, LogLevelDebug, format, args...)
	logData := writeFile(LogLevelDebug, format, args...)
	select {
	case f.logDataChal <- logData:
	default: // 如果logData超出缓冲区，丢弃该条日志信息
	}
}

func (f *FileLogger) Trace(format string, args ...interface{}) {
	if f.level > LogLevelTrace {
		return
	}
	//writeFile(f.appFile, LogLevelTrace, format, args...)
	logData := writeFile(LogLevelTrace, format, args...)
	select {
	case f.logDataChal <- logData:
	default: // 如果logData超出缓冲区，丢弃该条日志信息
	}
}

func (f *FileLogger) Info(format string, args ...interface{}) {
	if f.level > LogLevelInfo {
		return
	}
	//writeFile(f.appFile, LogLevelInfo, format, args...)
	logData := writeFile(LogLevelInfo, format, args...)
	select {
	case f.logDataChal <- logData:
	default: // 如果logData超出缓冲区，丢弃该条日志信息
	}
}

func (f *FileLogger) Warn(format string, args ...interface{}) {
	if f.level > LogLevelWarn {
		return
	}
	//writeFile(f.errFile, LogLevelWarn, format, args...)
	logData := writeFile(LogLevelWarn, format, args...)
	select {
	case f.logDataChal <- logData:
	default: // 如果logData超出缓冲区，丢弃该条日志信息
	}
}

func (f *FileLogger) Error(format string, args ...interface{}) {
	if f.level > LogLevelError {
		return
	}
	//writeFile(f.errFile, LogLevelError, format, args...)
	logData := writeFile(LogLevelError, format, args...)
	select {
	case f.logDataChal <- logData:
	default: // 如果logData超出缓冲区，丢弃该条日志信息
	}
}

func (f *FileLogger) Fatal(format string, args ...interface{}) {
	if f.level > LogLevelFatal {
		return
	}
	//writeFile(f.errFile, LogLevelFatal, format, args...)
	logData := writeFile(LogLevelFatal, format, args...)
	select {
	case f.logDataChal <- logData:
	default: // 如果logData超出缓冲区，丢弃该条日志信息
	}
}

func (f *FileLogger) Close() {
	f.appFile.Close()
	f.errFile.Close()
}
