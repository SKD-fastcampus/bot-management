# 인프라스트럭처 레이어

이 디렉토리는 외부 시스템과의 통신 및 기술적 세부 사항을 구현하는 인프라스트럭처 레이어를 포함합니다.

## 구조

- **http/**: HTTP 서버 및 핸들러 구현
- **grpc/**: gRPC 서버 및 서비스 구현
- **geolite/**: GeoLite2 데이터베이스 관련 구현

## HTTP 서버

`http/server.go`는 Echo 프레임워크를 사용한 HTTP 서버를 제공합니다:

```go
// HTTP 서버 생성 예시
httpServer := http.NewServer(
    http.WithPort(8080),
    http.WithLogger(zapLogger),
)

// 라우트 등록
httpServer.RegisterRoutes(func(e *echo.Echo) {
    e.GET("/api/v1/geo", geoHandler.GetGeoInfo)
})

// 서버 시작
go httpServer.Start()
```

## gRPC 서버

`grpc/server.go`는 gRPC 서버를 제공합니다:

```go
// gRPC 서버 생성 예시
grpcServer := grpc.NewServer(
    grpc.WithPort(9090),
    grpc.WithLogger(zapLogger),
)

// 서비스 등록
grpcServer.RegisterService(func(server *grpc.Server) {
    pb.RegisterGeoServiceServer(server, geoServiceImpl)
})

// 서버 시작
go grpcServer.Start()
```

## 로깅

서버 구현은 `pkg/logger` 패키지의 zap 로거를 사용하여 구조화된 로깅을 제공합니다:

- HTTP 요청/응답 로깅
- gRPC 요청/응답 로깅
- 에러 로깅 