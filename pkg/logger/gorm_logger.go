// File: pkg/logger/gorm_logger.go
package logger

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// ZapGormLogger는 gorm의 logger.Interface를 구현하며, 모든 GORM 로그를 zap 로 기록합니다.
// 추가로 slow query 임계시간, RecordNotFound 에러 무시 옵션 등을 제공합니다.
type ZapGormLogger struct {
	// LogLevel은 기록할 로그의 최소 레벨을 지정합니다.
	// (gormlogger.Silent, Error, Warn, Info 중 하나)
	LogLevel gormlogger.LogLevel
	// SlowThreshold는 쿼리 실행 시간이 이 시간보다 길면 slow query로 판단하여 Warn 레벨 로그를 남깁니다.
	// 0이면 사용하지 않습니다.
	SlowThreshold time.Duration
	// IgnoreRecordNotFoundError가 true이면 gorm.ErrRecordNotFound 에러는 로그에 남기지 않습니다.
	IgnoreRecordNotFoundError bool
	// zap 로거 인스턴스
	Logger *zap.Logger
}

// NewGormLogger는 지정한 옵션을 가진 ZapGormLogger 인스턴스를 생성합니다.
func NewGormLogger(logger *zap.Logger, level gormlogger.LogLevel, slowThreshold time.Duration, ignoreRecordNotFoundError bool) *ZapGormLogger {
	return &ZapGormLogger{
		LogLevel:                  level,
		SlowThreshold:             slowThreshold,
		IgnoreRecordNotFoundError: ignoreRecordNotFoundError,
		Logger:                    logger,
	}
}

// LogMode는 로그 레벨을 변경한 새로운 로거 인스턴스를 반환합니다.
func (z *ZapGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *z
	newLogger.LogLevel = level
	return &newLogger
}

// Info는 일반 정보를 zap을 통해 기록합니다.
func (z *ZapGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if z.LogLevel < gormlogger.Info {
		return
	}
	// context에서 추가 정보를 추출할 수 있다면 여기에 삽입 (예: request id 등)
	z.Logger.Sugar().Infof(msg, data...)
}

// Warn는 경고 로그를 zap을 통해 기록합니다.
func (z *ZapGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if z.LogLevel < gormlogger.Warn {
		return
	}
	z.Logger.Sugar().Warnf(msg, data...)
}

// Error는 에러 로그를 zap을 통해 기록합니다.
func (z *ZapGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if z.LogLevel < gormlogger.Error {
		return
	}
	// 에러 발생시 스택 트레이스 정보도 함께 남깁니다.
	z.Logger.Sugar().Errorf(msg, data...)
}

// Trace는 쿼리 실행 시간, SQL, 영향을 받은 행 수, 에러 등 상세 정보를 zap을 통해 기록합니다.
func (z *ZapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// 쿼리 실행 시간 측정
	elapsed := time.Since(begin)
	// 실행된 SQL과 영향을 받은 행 수를 얻음
	sql, rows := fc()

	// context의 deadline 정보가 있다면 로그에 포함 (디버그 목적으로)
	if deadline, ok := ctx.Deadline(); ok {
		z.Logger.Debug("Context deadline", zap.Time("deadline", deadline))
	}

	// 에러가 발생한 경우
	if err != nil && (!z.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)) {
		// RecordNotFound 에러가 무시 대상이 아니거나, 그 외의 에러인 경우 에러 로그 기록
		z.Logger.Error("GORM Trace Error",
			zap.Error(err),
			zap.Duration("elapsed", elapsed),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Stack("stack"),
		)
		return
	}

	// 쿼리 실행 시간이 설정된 임계시간보다 길 경우 (slow query)
	if z.SlowThreshold != 0 && elapsed > z.SlowThreshold {
		z.Logger.Warn("GORM Slow Query",
			zap.Duration("elapsed", elapsed),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
		)
		return
	}

	// 일반 쿼리 로그 기록 (LogLevel이 Info 이상인 경우)
	if z.LogLevel >= gormlogger.Info {
		z.Logger.Info("GORM Query",
			zap.Duration("elapsed", elapsed),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
		)
	}
}
