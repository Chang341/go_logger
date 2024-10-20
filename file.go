package go_logger

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type FileLogger struct {
	level         int
	logPath       string
	logName       string
	appFile       *os.File
	errFile       *os.File
	logDataChal   chan *LogData // 声明管道/队列，存放LogData日志数据的指针
	IsErrLog      bool          // 标识写到哪个log文件中
	logSplitType  int           // 日志拆分类型：小时/大小
	logSplitSize  int64         // 按大小拆分，多大
	lastSplitHour int           // 上一次拆分的时间
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

	logSplitTypeStr, ok := config["log_split_type"]
	var logSplitSize int64
	if !ok {
		logSplitTypeStr = "hour" // 默认按照小时拆分
	} else { // 按大小拆分
		logSplitSizeStr, ok := config["log_split_size"]
		if !ok {
			logSplitSizeStr = "104857600" // 100M
		}
		logSplitSize, err = strconv.ParseInt(logSplitSizeStr, 10, 64)
		if err != nil {
			logSplitSize = 104857600
		}
	}
	var logSplitType int
	if logSplitTypeStr == "size" {
		logSplitType = SplitLogBySize
	} else {
		logSplitType = SplitLogByHour
	}

	log = &FileLogger{
		level:        level,
		logPath:      logPath,
		logName:      logName,
		logDataChal:  make(chan *LogData, chanSize),
		IsErrLog:     isErrLog,
		logSplitType: logSplitType,
		logSplitSize: logSplitSize,
	}
	log.Init()
	return
}

func (f *FileLogger) Init() {
	// 创建appFile日志文件：存放debug、trace、info级别的日志
	fileName := fmt.Sprintf(AppLogFormat, f.logPath, f.logName)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Printf("open file %s failed, err:%v", fileName, err)
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
		f.checkSplitFile(f.IsErrLog)
		fmt.Fprintf(file, LogDataFormat,
			logData.TimeStr, logData.LevelStr,
			logData.FileName, logData.FuncName, logData.LineNo, logData.Message)
	}
}

func (f *FileLogger) checkSplitFile(isErrLog bool) {
	hour := time.Now().Hour()
	if f.logSplitType == SplitLogByHour {
		// 判断是否需要拆分
		if f.lastSplitHour == hour {
			return // 不需要拆分
		} else {
			// 按照小时拆分
			f.splitFile(isErrLog)
			f.lastSplitHour = hour
		}
	} else {
		// 判断是否需要拆分
		file := f.appFile
		if isErrLog {
			file = f.errFile
		}
		stat, err := file.Stat()
		if err != nil {
			return
		}
		fileSize := stat.Size()
		if fileSize <= f.logSplitSize {
			return
		} else {
			// 按传入大小拆分
			f.splitFile(isErrLog)
		}
	}
}

func (f *FileLogger) splitFile(isErrLog bool) {
	now := time.Now()
	fileName := fmt.Sprintf(AppLogFormat, f.logPath, f.logName)
	backupFileName := fmt.Sprintf(AppLogFormat+"%04d%02d%02d%02d%02d%02d", f.logPath, f.logName,
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	file := f.appFile
	if isErrLog {
		file = f.errFile
		fileName = fmt.Sprintf(ErrLogFormat, f.logPath, f.logName)
		backupFileName = fmt.Sprintf(ErrLogFormat+"%04d%02d%02d%02d%02d%02d", f.logPath, f.logName,
			now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	}
	file.Close()
	os.Rename(fileName, backupFileName) // 备份

	// 重新开一个文件，并将句柄赋给相应的参数
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return
	}
	f.appFile = file
	if isErrLog {
		f.errFile = file
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
