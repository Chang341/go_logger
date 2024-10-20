package logger

import "testing"

func TestFileLogger(t *testing.T) {
	logger := NewFileLoggger(LogLevelDebug, "D:/go_workspace/utils/logger/logs", "test")
	logger.Debug("user id[%d] is come from china", 321313)
	logger.Error("test error log")
	logger.Fatal("test fatal log")
	logger.Close()
}

func TestConsoleLogger(t *testing.T) {
	logger := NewConsoleLoggger(LogLevelDebug)
	logger.Debug("user id[%d] is come from china", 321313)
	logger.Error("test error log")
	logger.Fatal("test fatal log")
	logger.Close()

}
