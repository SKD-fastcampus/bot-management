// File: pkg/logger/echo_logger.go
package logger

import (
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

// NewEchoRequestLogger는 Echo 서버를 위한 Request Logger를 생성합니다.
// zap을 사용하여 HTTP 요청과 응답을 로깅합니다.
func NewEchoRequestLogger(logger *zap.Logger) echo.MiddlewareFunc {
	config := middleware.RequestLoggerConfig{
		// 특정 경로(예: /)를 로그에서 제외하고 싶다면 Skipper 사용
		Skipper: func(c echo.Context) bool {
			return c.Request().URL.Path == "/health" || c.Request().URL.Path == "/metrics"
		},
		// 다음 미들웨어나 핸들러가 실행되기 전 실행되는 함수
		BeforeNextFunc: func(c echo.Context) {
			// Request 시작 시간을 context에 저장
			c.Set("request-start-time", time.Now())
		},
		// 에러도 글로벌 핸들러에게 넘길지 여부 (원하는 동작에 맞추어 선택)
		HandleError: true,

		// 로그로 남길 항목 설정 (필요한 것들을 true로)
		LogLatency:       true, // 핸들러 체인을 실행한 뒤 소요된 시간 기록
		LogProtocol:      true, // HTTP/1.1, HTTP/2 등 프로토콜 정보
		LogRemoteIP:      true, // 클라이언트 IP (echo.Context.RealIP() 기준)
		LogHost:          true, // Host 정보 (예: example.com)
		LogMethod:        true, // HTTP 메서드 (GET, POST 등)
		LogURI:           true, // 요청 URI (/users?lang=en 등)
		LogURIPath:       true, // 요청 Path 부분 (/users 등)
		LogRoutePath:     true, // echo 라우팅 경로 (/users/:id 등)
		LogRequestID:     true, // X-Request-ID 헤더 또는 자동 생성된 Request ID
		LogReferer:       true, // Referer 정보
		LogUserAgent:     true, // User-Agent 정보
		LogStatus:        true, // 응답 상태 코드
		LogError:         true, // next(...)에서 발생한 에러
		LogContentLength: true, // 요청 헤더의 Content-Length
		LogResponseSize:  true, // 실제 응답의 바이트 수

		// 특정 Header/Query Param/Form Value 등을 추가로 기록하고 싶다면
		// 필요한 이름들을 아래와 같이 배열에 추가
		LogHeaders:     []string{"Content-Type", "Accept", "Authorization"},
		LogQueryParams: []string{"lang", "page", "limit"},
		LogFormValues:  []string{},

		// 중요: 요청/응답 정보를 실제로 어떻게 로깅할지 결정하는 함수
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			// BeforeNextFunc에서 설정한 정보 가져오기
			startTime, _ := c.Get("request-start-time").(time.Time)
			elapsed := time.Since(startTime)

			// 로그 필드를 zap.Field 형태로 구성
			fields := []zap.Field{
				zap.String("request.remote_ip", v.RemoteIP),
				zap.String("request.host", v.Host),
				zap.String("request.protocol", v.Protocol),
				zap.String("request.method", v.Method),
				zap.String("request.uri", v.URI),
				zap.String("request.path", v.URIPath),
				zap.String("request.route", v.RoutePath),
				zap.String("request.user_agent", v.UserAgent),
				zap.String("request.referer", v.Referer),
				zap.Int("response.status", v.Status),
				zap.Duration("response.latency", v.Latency),
				zap.String("response.latency_human", v.Latency.String()),
				zap.Duration("response.elapsed_since_before_next", elapsed),
				zap.String("request.request_id", v.RequestID),
				zap.Int64("response.response_size", v.ResponseSize),
				zap.String("request.content_length", v.ContentLength),
			}

			// Header, QueryParam, FormValue 같은 slice 형태 데이터들을 로깅하는 예시
			if len(v.Headers) > 0 {
				// Authorization 헤더 내용은 마스킹 처리
				headers := make(map[string]string)
				for k, values := range v.Headers {
					if len(values) > 0 {
						if k == "Authorization" {
							// Bearer 토큰 일부만 표시 (예: "Bearer xxxx...xxxx")
							val := values[0]
							if len(val) > 15 {
								headers[k] = val[:10] + "..." + val[len(val)-5:]
							} else {
								headers[k] = "[MASKED]"
							}
						} else {
							headers[k] = values[0]
						}
					}
				}
				fields = append(fields, zap.Any("request.headers", headers))
			}

			if len(v.QueryParams) > 0 {
				fields = append(fields, zap.Any("request.query_params", v.QueryParams))
			}

			if len(v.FormValues) > 0 {
				fields = append(fields, zap.Any("request.form_values", v.FormValues))
			}

			// 에러가 있는 경우
			if v.Error != nil {
				fields = append(fields, zap.Error(v.Error))
				// 에러이므로 Warn/Error 레벨로 로그를 찍을 수도 있습니다.
				logger.Error("Request failed", fields...)
				return nil
			}

			// 4XX 에러는 Warn 레벨로 기록
			if v.Status >= 400 && v.Status < 500 {
				logger.Warn("Client error", fields...)
				return nil
			}

			// 5XX 에러는 Error 레벨로 기록
			if v.Status >= 500 {
				logger.Error("Server error", fields...)
				return nil
			}

			// 정상 응답의 경우 Info 레벨로 기록
			logger.Info("Request completed", fields...)
			return nil
		},
	}

	return middleware.RequestLoggerWithConfig(config)
}

// WithEchoLogger Echo에 대한 커스텀 에러 핸들러를 설정합니다.
func WithEchoLogger(e *echo.Echo, logger *zap.Logger) {
	// zap 로거를 사용하는 Echo 내장 Logger 구현체
	e.Logger = NewEchoZapLogger(logger)

	// Echo 에러 핸들러 설정
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}

		// 에러 로그 기록
		logger.Error("HTTP error",
			zap.Error(err),
			zap.Int("status", code),
			zap.String("method", c.Request().Method),
			zap.String("path", c.Request().URL.Path),
			zap.String("ip", c.RealIP()),
		)

		// 에러 응답
		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead {
				err = c.NoContent(code)
			} else {
				err = c.JSON(code, map[string]interface{}{
					"error": http.StatusText(code),
				})
			}
			if err != nil {
				logger.Error("Failed to send error response", zap.Error(err))
			}
		}
	}
}

// EchoZapLogger는 echo.Logger 인터페이스를 구현한 zap 로거 래퍼입니다.
type EchoZapLogger struct {
	Logger *zap.Logger
}

// NewEchoZapLogger는 Echo의 Logger 인터페이스를 구현한 zap 로거 래퍼를 생성합니다.
func NewEchoZapLogger(logger *zap.Logger) *EchoZapLogger {
	return &EchoZapLogger{Logger: logger}
}

// Output Echo 로깅을 위한 Writer를 반환합니다.
func (l *EchoZapLogger) Output() io.Writer {
	return &zapWriter{logger: l.Logger}
}

// SetOutput Echo 로깅을 위한 Writer를 설정합니다. (zap에서는 무시됨)
func (l *EchoZapLogger) SetOutput(w io.Writer) {
	// zap에서는 무시됨
}

// Level Echo 로깅 레벨을 반환합니다.
func (l *EchoZapLogger) Level() log.Lvl {
	return log.INFO
}

// SetLevel Echo 로깅 레벨을 설정합니다. (zap에서는 무시됨)
func (l *EchoZapLogger) SetLevel(v log.Lvl) {
	// zap에서는 무시됨
}

// SetHeader Echo 로그 헤더를 설정합니다. (zap에서는 무시됨)
func (l *EchoZapLogger) SetHeader(h string) {
	// zap에서는 무시됨
}

// Prefix 로그 프리픽스를 반환합니다. (zap에서는 사용되지 않음)
func (l *EchoZapLogger) Prefix() string {
	return ""
}

// SetPrefix 로그 프리픽스를 설정합니다. (zap에서는 무시됨)
func (l *EchoZapLogger) SetPrefix(p string) {
	// zap에서는 무시됨
}

// Print zap 로거로 INFO 레벨 로그를 기록합니다.
func (l *EchoZapLogger) Print(i ...interface{}) {
	l.Logger.Sugar().Info(i...)
}

// Printf zap 로거로 INFO 레벨 로그를 기록합니다. (포맷 지정)
func (l *EchoZapLogger) Printf(format string, i ...interface{}) {
	l.Logger.Sugar().Infof(format, i...)
}

// Printj zap 로거로 INFO 레벨 로그를 기록합니다. (JSON 형식)
func (l *EchoZapLogger) Printj(j log.JSON) {
	l.Logger.Info("json_message", zap.Any("json", j))
}

// Debug zap 로거로 DEBUG 레벨 로그를 기록합니다.
func (l *EchoZapLogger) Debug(i ...interface{}) {
	l.Logger.Sugar().Debug(i...)
}

// Debugf zap 로거로 DEBUG 레벨 로그를 기록합니다. (포맷 지정)
func (l *EchoZapLogger) Debugf(format string, i ...interface{}) {
	l.Logger.Sugar().Debugf(format, i...)
}

// Debugj zap 로거로 DEBUG 레벨 로그를 기록합니다. (JSON 형식)
func (l *EchoZapLogger) Debugj(j log.JSON) {
	l.Logger.Debug("json_message", zap.Any("json", j))
}

// Info zap 로거로 INFO 레벨 로그를 기록합니다.
func (l *EchoZapLogger) Info(i ...interface{}) {
	l.Logger.Sugar().Info(i...)
}

// Infof zap 로거로 INFO 레벨 로그를 기록합니다. (포맷 지정)
func (l *EchoZapLogger) Infof(format string, i ...interface{}) {
	l.Logger.Sugar().Infof(format, i...)
}

// Infoj zap 로거로 INFO 레벨 로그를 기록합니다. (JSON 형식)
func (l *EchoZapLogger) Infoj(j log.JSON) {
	l.Logger.Info("json_message", zap.Any("json", j))
}

// Warn zap 로거로 WARN 레벨 로그를 기록합니다.
func (l *EchoZapLogger) Warn(i ...interface{}) {
	l.Logger.Sugar().Warn(i...)
}

// Warnf zap 로거로 WARN 레벨 로그를 기록합니다. (포맷 지정)
func (l *EchoZapLogger) Warnf(format string, i ...interface{}) {
	l.Logger.Sugar().Warnf(format, i...)
}

// Warnj zap 로거로 WARN 레벨 로그를 기록합니다. (JSON 형식)
func (l *EchoZapLogger) Warnj(j log.JSON) {
	l.Logger.Warn("json_message", zap.Any("json", j))
}

// Error zap 로거로 ERROR 레벨 로그를 기록합니다.
func (l *EchoZapLogger) Error(i ...interface{}) {
	l.Logger.Sugar().Error(i...)
}

// Errorf zap 로거로 ERROR 레벨 로그를 기록합니다. (포맷 지정)
func (l *EchoZapLogger) Errorf(format string, i ...interface{}) {
	l.Logger.Sugar().Errorf(format, i...)
}

// Errorj zap 로거로 ERROR 레벨 로그를 기록합니다. (JSON 형식)
func (l *EchoZapLogger) Errorj(j log.JSON) {
	l.Logger.Error("json_message", zap.Any("json", j))
}

// Fatal zap 로거로 FATAL 레벨 로그를 기록하고 프로그램을 종료합니다.
func (l *EchoZapLogger) Fatal(i ...interface{}) {
	l.Logger.Sugar().Fatal(i...)
}

// Fatalf zap 로거로 FATAL 레벨 로그를 기록하고 프로그램을 종료합니다. (포맷 지정)
func (l *EchoZapLogger) Fatalf(format string, i ...interface{}) {
	l.Logger.Sugar().Fatalf(format, i...)
}

// Fatalj zap 로거로 FATAL 레벨 로그를 기록하고 프로그램을 종료합니다. (JSON 형식)
func (l *EchoZapLogger) Fatalj(j log.JSON) {
	l.Logger.Fatal("json_message", zap.Any("json", j))
}

// Panic zap 로거로 PANIC 레벨 로그를 기록하고 패닉을 발생시킵니다.
func (l *EchoZapLogger) Panic(i ...interface{}) {
	l.Logger.Sugar().Panic(i...)
}

// Panicf zap 로거로 PANIC 레벨 로그를 기록하고 패닉을 발생시킵니다. (포맷 지정)
func (l *EchoZapLogger) Panicf(format string, i ...interface{}) {
	l.Logger.Sugar().Panicf(format, i...)
}

// Panicj zap 로거로 PANIC 레벨 로그를 기록하고 패닉을 발생시킵니다. (JSON 형식)
func (l *EchoZapLogger) Panicj(j log.JSON) {
	l.Logger.Panic("json_message", zap.Any("json", j))
}

// zapWriter는 io.Writer 인터페이스를 구현한 zap 로거 래퍼입니다.
type zapWriter struct {
	logger *zap.Logger
}

// Write는 io.Writer 인터페이스 구현을 위한 메서드입니다.
func (w *zapWriter) Write(p []byte) (n int, err error) {
	w.logger.Info(string(p))
	return len(p), nil
}
