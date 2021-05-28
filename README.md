# zaplog
基于 zap 和 lumberjack 配置的日志输出工具

#### Example

```go
	logConfig := &zaplog.Config{
		Level:      "info",		
		LogFormat:  "json",			// 日志格式：json 或 文本格式，默认为文本
		LogPath:    "/data/logs/",
		Stacktrace: false,			// Error 及以上是否输出堆栈信息
		CallerSkip: false,			// 是否输出Caller上级调用位置
		Stdout:     false,			// 是否到标准输出，若到标准输出则不生成日志文件
		App:        "app",			// app  
		Group:      "",
		MaxSize:    0,
		MaxBackups: 0,
		MaxAge:     0,
		Compress:   false,
	}

	zaplog.InitLogger(logConfig)
	zaplog.Sync()

	zaplog.Infof("This is info ....")
	zaplog.Warnf("This is warn ....")
```

