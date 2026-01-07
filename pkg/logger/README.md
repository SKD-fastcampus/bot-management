# 로깅 패키지

이 패키지는 애플리케이션 전체에서 일관된 로깅을 제공하기 위한 래퍼 및 유틸리티 함수를 제공합니다.

## 주요 기능

- **구조화된 로깅**: JSON 형식의 로그 지원
- **zap 로깅**: 고성능 zap 로거 구현
- **GORM 로깅**: GORM ORM을 위한 로깅 어댑터
- **Echo 로깅**: Echo 웹 프레임워크를 위한 로깅 어댑터
- **로깅 레벨**: debug, info, warn, error 등 다양한 로깅 레벨 지원

## 사용 예시

### 기본 로거 생성

```go
import (
    "github.com/SKD-fastcampus/bot-management/pkg/logger"
)

func main() {
    // 기본 설정으로 로거 생성
    log := logger.DefaultZapLogger()
    
    // 로그 남기기
    log.Info("서비스 시작")
    log.Debug("디버그 정보", zap.String("key", "value"))
    log.Error("오류 발생", zap.Error(err))
}
```

### 사용자 정의 로거 설정

```go
import (
    "github.com/SKD-fastcampus/bot-management/pkg/logger"
)

func main() {
    // 로거 설정
    config := logger.Config{
        Level:       "debug",         // 로그 레벨
        Format:      "json",          // 로그 포맷 (json 또는 console)
        Output:      "stdout",        // 출력 대상 (stdout, stderr, file)
        FilePath:    "app.log",       // 파일 출력 시 경로
        Development: true,            // 개발 모드 여부
    }
    
    // 설정으로 로거 생성
    log, err := logger.NewZapLogger(config)
    if err != nil {
        panic(err)
    }
    
    // 로거 사용
    log.Info("사용자 정의 로거 시작")
}
```

### GORM 로거 설정

```go
import (
    "time"
    
    "github.com/SKD-fastcampus/bot-management/pkg/logger"
    "gorm.io/gorm"
    gormlogger "gorm.io/gorm/logger"
)

func main() {
    // zap 로거 생성
    zapLogger := logger.DefaultZapLogger()
    
    // GORM 로거 생성
    gormLogger := logger.NewGormLogger(
        zapLogger,                       // zap 로거 인스턴스
        gormlogger.Info,                 // 로그 레벨
        time.Second,                     // Slow Query 임계값
        true,                            // RecordNotFound 에러 무시 여부
    )
    
    // GORM DB 인스턴스 생성 시 로거 설정
    db, err := gorm.Open(dialect, &gorm.Config{
        Logger: gormLogger,
    })
}
```

### Echo 로거 설정

```go
import (
    "github.com/labstack/echo/v4"
    "github.com/SKD-fastcampus/bot-management/pkg/logger"
)

func main() {
    // zap 로거 생성
    zapLogger := logger.DefaultZapLogger()
    
    // Echo 인스턴스 생성
    e := echo.New()
    
    // Echo에 zap 로거 설정
    logger.WithEchoLogger(e, zapLogger)
    
    // HTTP 요청 로깅 미들웨어 설정
    e.Use(logger.NewEchoRequestLogger(zapLogger))
    
    // 라우트 등록
    e.GET("/", func(c echo.Context) error {
        return c.String(200, "Hello, World!")
    })
    
    // 서버 시작
    e.Start(":8080")
} 