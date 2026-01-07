# 공유 패키지

이 디렉토리에는 여러 서비스에서 공유하는 패키지가 포함되어 있습니다. 이 패키지들은 기본적인 기능을 제공하며 모든 서비스에서 일관되게 사용할 수 있도록 설계되었습니다.

## 디렉토리 구조

- **logger**: 구조화된 로깅을 위한 패키지
- **errors**: 표준화된 에러 처리를 위한 패키지
- **middleware**: HTTP 및 gRPC 미들웨어
- **database**: 데이터베이스 연결 및 공통 함수
- **validator**: 입력 유효성 검사를 위한 패키지
- **crypto**: 암호화 및 보안 관련 기능
- **tracing**: 분산 추적 기능

## 사용 방법

각 서비스에서 필요한 공유 패키지를 가져와서 사용할 수 있습니다:

```go
import (
    "github.com/SKD-fastcampus/bot-management/pkg/logger"
    "github.com/SKD-fastcampus/bot-management/pkg/errors"
)

func main() {
    log := logger.New()
    log.Info("서비스 시작")
    
    if err := doSomething(); err != nil {
        log.WithError(errors.Wrap(err, "작업 실패")).Error("서비스 실패")
    }
}
```

## 패키지 추가하기

새로운 공유 패키지를 추가할 때는 다음 가이드라인을 따라주세요:

1. 패키지는 단일 책임 원칙을 따라야 합니다.
2. 패키지는 명확한 API와 충분한 문서를 제공해야 합니다.
3. 특정 서비스에 종속되는 코드는 포함하지 마세요.
4. 테스트 코드를 반드시 작성하세요.

## 주의사항

- 공유 패키지에서 다른 공유 패키지를 사용할 때 순환 의존성을 주의하세요.
- 새로운 버전이 기존 서비스의 동작을 깨뜨리지 않도록 호환성을 유지하세요. 