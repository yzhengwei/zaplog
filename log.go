package zaplog

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger         *zap.SugaredLogger
	defaultApp     = "app"
	defaultGroup   = "group"
	defaultLogPath = "/data/logs/"
)

func InitLogger(cfg *Config) {
	logLevel := zap.DebugLevel
	switch cfg.Level {
	case "debug":
		logLevel = zap.DebugLevel
	case "info":
		logLevel = zap.InfoLevel
	case "warn":
		logLevel = zap.WarnLevel
	case "error":
		logLevel = zap.ErrorLevel
	default:
		logLevel = zap.InfoLevel
	}

	var (
		coreArr []zapcore.Core
		core    zapcore.Core
	)

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	errorWriter := getLogWriter("error", cfg)
	coreArr = append(coreArr, zapcore.NewCore(getEncoder(cfg.LogFormat), zapcore.AddSync(errorWriter), errorLevel))

	if logLevel <= zapcore.WarnLevel {
		warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl == zapcore.WarnLevel
		})
		warnWriter := getLogWriter("warn", cfg)
		coreArr = append(coreArr, zapcore.NewCore(getEncoder(cfg.LogFormat), zapcore.AddSync(warnWriter), warnLevel))
	}

	if logLevel <= zapcore.InfoLevel {
		infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl == zapcore.InfoLevel
		})
		infoWriter := getLogWriter("info", cfg)
		coreArr = append(coreArr, zapcore.NewCore(getEncoder(cfg.LogFormat), zapcore.AddSync(infoWriter), infoLevel))
	}

	if logLevel <= zapcore.DebugLevel {
		debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl == zapcore.DebugLevel
		})
		debugWriter := getLogWriter("debug", cfg)
		coreArr = append(coreArr, zapcore.NewCore(getEncoder(cfg.LogFormat), zapcore.AddSync(debugWriter), debugLevel))
	}

	core = zapcore.NewTee(coreArr...)

	sugarLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	if !cfg.Stacktrace && cfg.CallerSkip {
		sugarLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
	} else if cfg.Stacktrace && !cfg.CallerSkip {
		sugarLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel))
	} else if cfg.Stacktrace && cfg.CallerSkip {
		sugarLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2), zap.AddStacktrace(zap.ErrorLevel))
	}

	Logger = sugarLogger.Sugar()
}

func getEncoder(format string) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	//encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	if format == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(level string, cfg *Config) zapcore.WriteSyncer {
	if !cfg.Stdout {
		lumberJackLogger := &lumberjack.Logger{
			MaxSize:    1024,         // 每个日志文件保存的最大尺寸 单位：M
			MaxBackups: 10,           // 日志文件最多保存多少个备份
			MaxAge:     3,            // 文件最多保存多少天
			Compress:   cfg.Compress, // 是否压缩
			LocalTime:  true,
		}

		fileName := defaultLogPath
		if cfg.LogPath != "" {
			fileName = cfg.LogPath
		}
		if cfg.Group != "" {
			fileName = fileName + cfg.Group + "_"
		} else {
			fileName = fileName + defaultGroup + "_"
		}
		if cfg.App != "" {
			fileName = fileName + cfg.App + "_"
		} else {
			fileName = fileName + defaultApp + "_"
		}
		fileName = fileName + level + ".log"

		lumberJackLogger.Filename = fileName

		if cfg.MaxSize > 0 {
			lumberJackLogger.MaxSize = cfg.MaxSize
		}
		if cfg.MaxBackups > 0 {
			lumberJackLogger.MaxBackups = cfg.MaxBackups
		}
		if cfg.MaxAge > 0 {
			lumberJackLogger.MaxAge = cfg.MaxAge
		}

		return zapcore.AddSync(lumberJackLogger)
	}
	return zapcore.AddSync(os.Stdout)
}

func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	Logger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	Logger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	Logger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	Logger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	Logger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	Logger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	Logger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	Logger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	Logger.Fatalf(template, args...)
}

func Sync() {
	Logger.Sync()
}
