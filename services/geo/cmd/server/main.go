package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	pb "github.com/SKD-fastcampus/bot-management/proto/geo/v1"
	grpcHandler "github.com/SKD-fastcampus/bot-management/services/geo/internal/adapter/handler/grpc"
	httpHandler "github.com/SKD-fastcampus/bot-management/services/geo/internal/adapter/handler/http"
	"github.com/SKD-fastcampus/bot-management/services/geo/internal/adapter/repository"
	"github.com/SKD-fastcampus/bot-management/services/geo/internal/config"
	grpcServer "github.com/SKD-fastcampus/bot-management/services/geo/internal/infrastructure/grpc"
	httpServer "github.com/SKD-fastcampus/bot-management/services/geo/internal/infrastructure/http"
	"github.com/SKD-fastcampus/bot-management/services/geo/internal/usecase"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// 1. 설정 로드
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("설정 로드 실패: %v", err))
	}

	// 2. 로거 가져오기
	log := cfg.Logger
	log.Info("GEO 서비스 시작")

	// 3. GeoLite2 데이터베이스 경로 설정
	dataDir := filepath.Join("services", "geo", "data")
	if cfg.GeoLite.DbPath != "" {
		dataDir = cfg.GeoLite.DbPath
	}

	cityDbPath := filepath.Join(dataDir, "GeoLite2-City.mmdb")
	countryDbPath := filepath.Join(dataDir, "GeoLite2-Country.mmdb")
	asnDbPath := filepath.Join(dataDir, "GeoLite2-ASN.mmdb")

	// 4. GeoLite2 리포지토리 초기화
	log.Info("GeoLite2 데이터베이스 초기화 중...")
	geoRepo, err := repository.NewGeoLite2Repository(cityDbPath, countryDbPath, asnDbPath)
	if err != nil {
		log.Fatal("GeoLite2 리포지토리 초기화 실패", zap.Error(err))
	}
	defer geoRepo.Close()
	log.Info("GeoLite2 데이터베이스 초기화 완료")

	// 5. 유스케이스 초기화
	geoUseCase := usecase.NewGeoUseCaseWithGeoLite2(geoRepo)
	defer geoUseCase.Close()

	// 6. HTTP 핸들러 초기화
	geoHttpHandler := httpHandler.NewGeoHandler(geoUseCase)

	// 7. gRPC 핸들러 초기화
	geoGrpcHandler := grpcHandler.NewGeoHandler(geoUseCase)

	// 8. HTTP 서버 포트 설정
	httpPort := 8080
	if cfg.Server.HTTP.Port != "" {
		httpPort = parseInt(cfg.Server.HTTP.Port, 8080)
	}

	// 9. gRPC 서버 포트 설정
	grpcPort := 9090
	if cfg.Server.GRPC.Port != "" {
		grpcPort = parseInt(cfg.Server.GRPC.Port, 9090)
	}

	// 10. HTTP 서버 초기화 및 시작
	httpSrv := httpServer.NewServer(
		httpServer.WithPort(httpPort),
		httpServer.WithLogger(log),
	)

	// 라우트 등록
	httpSrv.RegisterRoutes(geoHttpHandler.RegisterRoutes)

	// HTTP 서버 시작
	go func() {
		if err := httpSrv.Start(); err != nil {
			log.Error("HTTP 서버 에러", zap.Error(err))
		}
	}()

	// 11. gRPC 서버 초기화 및 시작
	grpcSrv := grpcServer.NewServer(
		grpcServer.WithPort(grpcPort),
		grpcServer.WithLogger(log),
	)

	// gRPC 서비스 등록
	grpcSrv.RegisterService(func(server *grpc.Server) {
		pb.RegisterGeoServiceServer(server, geoGrpcHandler)
	})

	// gRPC 서버 시작
	go func() {
		if err := grpcSrv.Start(); err != nil {
			log.Error("gRPC 서버 에러", zap.Error(err))
		}
	}()

	log.Info("서버 실행 중...",
		zap.Int("http_port", httpPort),
		zap.Int("grpc_port", grpcPort),
	)

	// 12. 종료 시그널 처리
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("서버 종료 중...")

	// 13. 종료 타임아웃 설정
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 14. HTTP 서버 종료
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Error("HTTP 서버 종료 실패", zap.Error(err))
	}

	// 15. gRPC 서버 종료
	if err := grpcSrv.Shutdown(ctx); err != nil {
		log.Error("gRPC 서버 종료 실패", zap.Error(err))
	}

	log.Info("서버 정상 종료")
}

// parseInt는 문자열을 정수로 변환하고, 변환 실패 시 기본값을 반환합니다.
func parseInt(s string, defaultVal int) int {
	var val int
	if _, err := fmt.Sscanf(s, "%d", &val); err != nil {
		return defaultVal
	}
	return val
}
