package go_logger

import "testing"

func TestFileLogger(t *testing.T) {
	config := make(map[string]string)
	config["log_path"] = "D:/go_workspace/utils/logger/logs"
	config["log_name"] = "test"
	config["log_level"] = "debug"
	logger, _ := NewFileLoggger(config)
	logger.Debug("user id[%d] is come from china", 321313)
	logger.Error("test error log")
	logger.Fatal("test fatal log")
	logger.Close()
}

func TestConsoleLogger(t *testing.T) {
	config := make(map[string]string)
	config["log_level"] = "debug"
	logger, _ := NewConsoleLoggger(config)
	logger.Debug("user id[%d] is come from china", 321313)
	logger.Error("test error log")
	logger.Fatal("test fatal log")
	logger.Close()

}
