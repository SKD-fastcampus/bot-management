package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/SKD-fastcampus/bot-management/pkg/logger"
	"go.uber.org/zap"
)

// Server HTTP 서버 구조체입니다.
type Server struct {
	echo   *echo.Echo
	logger *zap.Logger
	port   int
}

// ServerOption Server 생성을 위한 옵션 함수 타입입니다.
type ServerOption func(*Server)

// WithPort 서버 포트를 설정하는 옵션입니다.
func WithPort(port int) ServerOption {
	return func(s *Server) {
		s.port = port
	}
}

// WithLogger 로거를 설정하는 옵션입니다.
func WithLogger(logger *zap.Logger) ServerOption {
	return func(s *Server) {
		s.logger = logger
	}
}

// NewServer HTTP 서버를 생성합니다.
func NewServer(opts ...ServerOption) *Server {
	// 기본 서버 설정
	s := &Server{
		echo:   echo.New(),
		logger: zap.NewNop(), // 기본은 로깅 없음
		port:   8080,         // 기본 포트
	}

	// 옵션 적용
	for _, opt := range opts {
		opt(s)
	}

	// Echo 인스턴스 설정
	e := s.echo

	// 로거 설정
	logger.WithEchoLogger(e, s.logger)

	// 미들웨어 설정
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(logger.NewEchoRequestLogger(s.logger))

	// 기본 라우트 설정
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})

	// 메트릭 엔드포인트
	e.GET("/metrics", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	return s
}

// RegisterRoutes 라우트를 등록하는 메서드입니다.
// 이 메서드는 핸들러를 등록하는 함수를 받아 실행합니다.
func (s *Server) RegisterRoutes(registerFunc func(e *echo.Echo)) {
	registerFunc(s.echo)
}

// Start 서버를 시작합니다.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	s.logger.Info("HTTP 서버 시작", zap.String("addr", addr))

	return s.echo.Start(addr)
}

// Shutdown 서버를 안전하게 종료합니다.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("HTTP 서버 종료 중...")
	return s.echo.Shutdown(ctx)
}

// GetEcho 내부 Echo 인스턴스를 반환합니다.
func (s *Server) GetEcho() *echo.Echo {
	return s.echo
}

// 사용 예시:
//
// func main() {
//     // zap 로거 생성
//     zapLogger := logger.DefaultZapLogger()
//
//     // HTTP 서버 생성
//     httpServer := http.NewServer(
//         http.WithPort(8080),
//         http.WithLogger(zapLogger),
//     )
//
//     // 라우트 등록
//     httpServer.RegisterRoutes(func(e *echo.Echo) {
//         e.GET("/api/v1/users", userHandler.GetUsers)
//         e.POST("/api/v1/users", userHandler.CreateUser)
//     })
//
//     // 서버 시작
//     go func() {
//         if err := httpServer.Start(); err != nil && err != http.ErrServerClosed {
//             zapLogger.Fatal("HTTP 서버 시작 실패", zap.Error(err))
//         }
//     }()
//
//     // 종료 시그널 대기
//     quit := make(chan os.Signal, 1)
//     signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
//     <-quit
//
//     // 서버 종료
//     ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//     defer cancel()
//     if err := httpServer.Shutdown(ctx); err != nil {
//         zapLogger.Fatal("서버 강제 종료", zap.Error(err))
//     }
// }
