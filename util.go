package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

func GetLineInfo() (fileName string, funcName string, lineNo int) {
	pc, file, line, ok := runtime.Caller(3)
	if ok {
		fileName = file
		funcName = runtime.FuncForPC(pc).Name()
		lineNo = line
	}
	return file, funcName, lineNo
}

func writeFile(file *os.File, level int, format string, args ...interface{}) {
	now := time.Now()
	nowStr := now.Format("2006-01-02 15:04:05.999") // 必须传这个时间点，格式可变/-等
	levelStr := getLevelText(level)
	fileName, funcName, lineNo := GetLineInfo()
	fileName = path.Base(fileName)
	funcName = path.Base(funcName)
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(file, "%s %s [%s:%s:%d] %s\n", nowStr, levelStr, fileName, funcName, lineNo, msg)
}
