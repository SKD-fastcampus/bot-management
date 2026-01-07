package logger

import (
	"context"
	"path"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewGrpcUnaryServerInterceptor는 단일 요청/응답 gRPC 메서드에 대한 로깅 인터셉터를 생성합니다.
func NewGrpcUnaryServerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// 요청 시작 시간 기록
		startTime := time.Now()

		// 메서드 이름 추출
		fullMethod := info.FullMethod
		service := path.Dir(fullMethod)[1:]
		method := path.Base(fullMethod)

		// 핸들러 호출 및 응답/에러 캡처
		resp, err = handler(ctx, req)

		// 소요 시간 계산
		duration := time.Since(startTime)

		// 상태 코드 추출
		statusCode := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				statusCode = st.Code()
			} else {
				statusCode = codes.Unknown
			}
		}

		// 로그 레벨 결정
		if statusCode == codes.OK {
			logger.Info("gRPC 요청 완료",
				zap.String("grpc.service", service),
				zap.String("grpc.method", method),
				zap.String("grpc.code", statusCode.String()),
				zap.Duration("grpc.duration", duration),
			)
		} else if statusCode == codes.Canceled || statusCode == codes.DeadlineExceeded || statusCode == codes.ResourceExhausted ||
			statusCode == codes.Aborted || statusCode == codes.Unavailable || statusCode == codes.DataLoss {
			logger.Warn("gRPC 요청 실패",
				zap.String("grpc.service", service),
				zap.String("grpc.method", method),
				zap.String("grpc.code", statusCode.String()),
				zap.Error(err),
				zap.Duration("grpc.duration", duration),
			)
		} else {
			logger.Error("gRPC 요청 오류",
				zap.String("grpc.service", service),
				zap.String("grpc.method", method),
				zap.String("grpc.code", statusCode.String()),
				zap.Error(err),
				zap.Duration("grpc.duration", duration),
			)
		}

		return resp, err
	}
}

// NewGrpcStreamServerInterceptor는 스트리밍 gRPC 메서드에 대한 로깅 인터셉터를 생성합니다.
func NewGrpcStreamServerInterceptor(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// 요청 시작 시간 기록
		startTime := time.Now()

		// 메서드 이름 추출
		fullMethod := info.FullMethod
		service := path.Dir(fullMethod)[1:]
		method := path.Base(fullMethod)

		// 시작 로그
		logger.Info("gRPC 스트림 시작",
			zap.String("grpc.service", service),
			zap.String("grpc.method", method),
			zap.Bool("grpc.is_client_stream", info.IsClientStream),
			zap.Bool("grpc.is_server_stream", info.IsServerStream),
		)

		// ServerStream을 래핑하여 메시지 카운팅
		wrappedStream := &wrappedServerStream{
			ServerStream: ss,
			recvCount:    0,
			sendCount:    0,
		}

		// 핸들러 호출 및 에러 캡처
		err := handler(srv, wrappedStream)

		// 소요 시간 계산
		duration := time.Since(startTime)

		// 상태 코드 추출
		statusCode := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				statusCode = st.Code()
			} else {
				statusCode = codes.Unknown
			}
		}

		// 로그 레벨 결정
		if statusCode == codes.OK {
			logger.Info("gRPC 스트림 완료",
				zap.String("grpc.service", service),
				zap.String("grpc.method", method),
				zap.String("grpc.code", statusCode.String()),
				zap.Int("grpc.recv_count", wrappedStream.recvCount),
				zap.Int("grpc.send_count", wrappedStream.sendCount),
				zap.Duration("grpc.duration", duration),
			)
		} else if statusCode == codes.Canceled || statusCode == codes.DeadlineExceeded || statusCode == codes.ResourceExhausted ||
			statusCode == codes.Aborted || statusCode == codes.Unavailable || statusCode == codes.DataLoss {
			logger.Warn("gRPC 스트림 실패",
				zap.String("grpc.service", service),
				zap.String("grpc.method", method),
				zap.String("grpc.code", statusCode.String()),
				zap.Error(err),
				zap.Int("grpc.recv_count", wrappedStream.recvCount),
				zap.Int("grpc.send_count", wrappedStream.sendCount),
				zap.Duration("grpc.duration", duration),
			)
		} else {
			logger.Error("gRPC 스트림 오류",
				zap.String("grpc.service", service),
				zap.String("grpc.method", method),
				zap.String("grpc.code", statusCode.String()),
				zap.Error(err),
				zap.Int("grpc.recv_count", wrappedStream.recvCount),
				zap.Int("grpc.send_count", wrappedStream.sendCount),
				zap.Duration("grpc.duration", duration),
			)
		}

		return err
	}
}

// wrappedServerStream은 ServerStream을 래핑하여 메시지 송수신 횟수를 추적합니다.
type wrappedServerStream struct {
	grpc.ServerStream
	recvCount int
	sendCount int
}

// RecvMsg는 메시지 수신 횟수를 추적합니다.
func (w *wrappedServerStream) RecvMsg(m interface{}) error {
	err := w.ServerStream.RecvMsg(m)
	if err == nil {
		w.recvCount++
	}
	return err
}

// SendMsg는 메시지 송신 횟수를 추적합니다.
func (w *wrappedServerStream) SendMsg(m interface{}) error {
	err := w.ServerStream.SendMsg(m)
	if err == nil {
		w.sendCount++
	}
	return err
}

// WithGrpcLogger는 gRPC 서버에 로깅 인터셉터를 설정합니다.
func WithGrpcLogger(server *grpc.Server, logger *zap.Logger) *grpc.Server {
	// 이미 생성된 서버에 인터셉터 추가는 지원되지 않으므로 주의 필요
	// 대신 사용 예시를 제공합니다.
	return server
}

// 사용 예시:
//
// import (
//     "google.golang.org/grpc"
//     "github.com/wekeepgrowing/semo-backend-monorepo/pkg/logger"
// )
//
// func NewGrpcServer(logger *zap.Logger) *grpc.Server {
//     // 인터셉터 생성
//     unaryInterceptor := logger.NewGrpcUnaryServerInterceptor(logger)
//     streamInterceptor := logger.NewGrpcStreamServerInterceptor(logger)
//
//     // gRPC 서버 생성 시 인터셉터 설정
//     server := grpc.NewServer(
//         grpc.UnaryInterceptor(unaryInterceptor),
//         grpc.StreamInterceptor(streamInterceptor),
//     )
//
//     return server
// }
