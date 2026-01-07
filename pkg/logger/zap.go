// File: pkg/logger/zap.go
package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config 로거 설정
type Config struct {
	// Level 로그 레벨 (debug, info, warn, error, dpanic, panic, fatal)
	Level string
	// Format 로그 포맷 (json, console)
	Format string
	// Output 로그 출력 대상 (stdout, stderr, file 등)
	Output string
	// FilePath 파일로 출력할 경우 파일 경로
	FilePath string
	// Development 개발 모드 여부
	Development bool
}

// NewZapLogger 새로운 zap 로거를 생성합니다.
func NewZapLogger(config Config) (*zap.Logger, error) {
	// 로그 레벨 설정
	level := zap.NewAtomicLevel()
	switch config.Level {
	case "debug":
		level.SetLevel(zapcore.DebugLevel)
	case "info":
		level.SetLevel(zapcore.InfoLevel)
	case "warn":
		level.SetLevel(zapcore.WarnLevel)
	case "error":
		level.SetLevel(zapcore.ErrorLevel)
	case "dpanic":
		level.SetLevel(zapcore.DPanicLevel)
	case "panic":
		level.SetLevel(zapcore.PanicLevel)
	case "fatal":
		level.SetLevel(zapcore.FatalLevel)
	default:
		level.SetLevel(zapcore.InfoLevel)
	}

	// 로그 인코더 설정
	var encoder zapcore.Encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "@timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.LevelKey = "log.level"
	encoderConfig.MessageKey = "message"
	encoderConfig.CallerKey = "caller"

	if config.Development {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	if config.Format == "console" || config.Development {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 로그 출력 설정
	var writeSyncer zapcore.WriteSyncer
	switch config.Output {
	case "stderr":
		writeSyncer = zapcore.AddSync(os.Stderr)
	case "file":
		if config.FilePath == "" {
			writeSyncer = zapcore.AddSync(os.Stdout)
		} else {
			file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return nil, err
			}
			writeSyncer = zapcore.AddSync(file)
		}
	default:
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// 로거 생성
	core := zapcore.NewCore(encoder, writeSyncer, level)
	logger := zap.New(core)

	// 개발 모드이면 호출자 정보 추가
	if config.Development {
		logger = logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
	}

	// 스택 트레이스 옵션 추가
	logger = logger.WithOptions(zap.AddStacktrace(zapcore.ErrorLevel))

	return logger, nil
}

// DefaultZapLogger 기본 설정으로 zap 로거를 생성합니다.
func DefaultZapLogger() *zap.Logger {
	config := Config{
		Level:       "info",
		Format:      "json",
		Output:      "stdout",
		Development: false,
	}

	logger, err := NewZapLogger(config)
	if err != nil {
		// 로거 생성 실패 시 기본 로거 반환
		return zap.NewExample()
	}
	return logger
}
