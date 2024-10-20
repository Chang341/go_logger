### go语言日志库 go_logger
#### 相关概念说明
- 日志按照严重程度分为：debug级别、trace级别、info级别、warn级别、error级别、fatal级别
- 支持文件存储以及控制台打印两种方式，文件存储按照应用日志(debug|trace|info)和错误日志(warn|error|fatal)进行划分
  - 文件存储通过异步方式写入，不会阻塞主进程
  - 文件存储支持小时级时间切片、文件大小切片
#### 使用说明
``` go
import github.com/Chang341/go_logger // 导入日志库
go_logger.Init(name, config)         // 使用-初始化「name为日志类型file/console，config传入配置信息」
go_logger.Debug("日志内容")            // 使用-Debug()/Trace()/Info()/Warn()/Error()/Fatal()
```
```go
// 使用示例
package main

import (
  "fmt"
  "github.com/Chang341/go_logger"
)

func initLogger(name, logPath, logName, logLevel string) (err error) {
	config := make(map[string]string)
	config["log_path"] = logPath
	config["log_name"] = logName
	config["log_level"] = logLevel
	config["log_split_type"] = "size"

	err = go_logger.Init(name, config)
	if err != nil {
		return
	}
	go_logger.Debug("init logger success")
	return
}

func main(){
    err := initLogger("file", #{logPath}, #{logName}, #{logLevel})
    if err != nil {
        fmt.Println(err)
        return
    }
}

```