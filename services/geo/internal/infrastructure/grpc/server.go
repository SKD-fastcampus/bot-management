package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/SKD-fastcampus/bot-management/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Server gRPC 서버 구조체입니다.
type Server struct {
	grpcServer *grpc.Server
	logger     *zap.Logger
	port       int
	listener   net.Listener
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

// NewServer gRPC 서버를 생성합니다.
func NewServer(opts ...ServerOption) *Server {
	// 기본 서버 설정
	s := &Server{
		logger: zap.NewNop(), // 기본은 로깅 없음
		port:   9090,         // 기본 포트
	}

	// 옵션 적용
	for _, opt := range opts {
		opt(s)
	}

	// gRPC 인터셉터 설정
	unaryInterceptor := logger.NewGrpcUnaryServerInterceptor(s.logger)
	streamInterceptor := logger.NewGrpcStreamServerInterceptor(s.logger)

	// gRPC 서버 생성
	s.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor),
		grpc.StreamInterceptor(streamInterceptor),
	)

	// 헬스 체크 서비스 등록
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(s.grpcServer, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// 리플렉션 서비스 등록 (gRPC 서버 탐색용, 개발 환경에서 유용)
	reflection.Register(s.grpcServer)

	return s
}

// RegisterService gRPC 서비스를 등록하는 메서드입니다.
// 이 메서드는 서비스 등록 함수를 받아 실행합니다.
func (s *Server) RegisterService(registerFunc func(server *grpc.Server)) {
	registerFunc(s.grpcServer)
}

// Start 서버를 시작합니다.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("gRPC 서버 리스닝 실패: %w", err)
	}
	s.listener = lis

	s.logger.Info("gRPC 서버 시작", zap.String("addr", addr))
	return s.grpcServer.Serve(lis)
}

// Shutdown 서버를 안전하게 종료합니다.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("gRPC 서버 종료 중...")
	stopped := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		// 컨텍스트 타임아웃 시 강제 종료
		s.logger.Warn("gRPC 서버 강제 종료")
		s.grpcServer.Stop()
		return ctx.Err()
	case <-stopped:
		// 정상 종료
		s.logger.Info("gRPC 서버 종료 완료")
		return nil
	}
}

// GetGrpcServer 내부 gRPC 서버 인스턴스를 반환합니다.
func (s *Server) GetGrpcServer() *grpc.Server {
	return s.grpcServer
}

// 사용 예시:
//
// func main() {
//     // zap 로거 생성
//     zapLogger := logger.DefaultZapLogger()
//
//     // gRPC 서버 생성
//     grpcServer := grpc.NewServer(
//         grpc.WithPort(9090),
//         grpc.WithLogger(zapLogger),
//     )
//
//     // gRPC 서비스 등록
//     grpcServer.RegisterService(func(server *grpc.Server) {
//         pb.RegisterGeoServiceServer(server, geoServiceImpl)
//     })
//
//     // 서버 시작
//     go func() {
//         if err := grpcServer.Start(); err != nil {
//             zapLogger.Fatal("gRPC 서버 시작 실패", zap.Error(err))
//         }
//     }()
//
//     // 종료 시그널 대기
//     quit := make(chan os.Signal, 1)
//     signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
//     <-quit
//
//     // 서버 종료
//     ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//     defer cancel()
//     if err := grpcServer.Shutdown(ctx); err != nil {
//         zapLogger.Fatal("gRPC 서버 강제 종료", zap.Error(err))
//     }
// }
