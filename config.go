package zaplog

type Config struct {
	Level      string `yaml:"level"`
	LogFormat  string `yaml:"log_format"` // 日志格式：json 或 文本格式，默认为文本
	LogPath    string `yaml:"log_path"`
	Stacktrace bool   `yaml:"stacktrace"`  // Error 及以上是否输出堆栈信息
	CallerSkip bool   `yaml:"caller_skip"` // 是否输出Caller上级调用位置
	Stdout     bool   `yaml:"stdout"`      // 是否到标准输出，若到标准输出则不生成日志文件

	App        string `yaml:"app"`
	Group      string `yaml:"group"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}
