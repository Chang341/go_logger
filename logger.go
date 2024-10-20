package go_logger

import "fmt"

var log LogInterface

// name参数传入 file/console 来确定初始化哪个对象
func Init(name string, config map[string]string) (err error) {
	switch name {
	case "file":
		log, err = NewFileLoggger(config)
	case "console":
		log, err = NewConsoleLoggger(config)
	default:
		err = fmt.Errorf("unsupport log name %s", name)
	}
	return
}

func Debug(format string, args ...interface{}) {
	log.Debug(format, args...)
}

func Trace(format string, args ...interface{}) {
	log.Trace(format, args...)
}

func Info(format string, args ...interface{}) {
	log.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	log.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	log.Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	log.Fatal(format, args...)
}
