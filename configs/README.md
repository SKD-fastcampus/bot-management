# 설정 파일

이 디렉토리에는 서비스 설정을 위한 파일들이 포함되어 있습니다. 각 환경(개발, 스테이징, 프로덕션)에 맞는 설정을 제공합니다.

## 디렉토리 구조

- **dev/**: 개발 환경 설정
- **staging/**: 스테이징 환경 설정
- **prod/**: 프로덕션 환경 설정

## 설정 파일 형식

설정 파일은 YAML 형식으로 작성되며, 서비스별로 구분됩니다:

```yaml
# 서비스 이름 (예: auth-service.yaml)
service:
  name: auth-service
  version: 1.0.0

server:
  port: 8082
  timeout: 30s

database:
  host: localhost
  port: 5432
  name: auth_db
  user: postgres
  
log:
  level: info
  format: json
```

## 설정 로드 방법

서비스는 다음과 같은 방법으로 설정을 로드합니다:

1. 환경 변수 `ENV`에 따라 적절한 환경 설정 디렉토리 선택 (기본값: dev)
2. 서비스는 pkg/config 패키지를 사용하여 설정을 로드합니다.

```go
import "github.com/SKD-fastcampus/bot-management/pkg/config"

func main() {
    cfg, err := config.Load("auth-service")
    if err != nil {
        log.Fatalf("설정 로드 실패: %v", err)
    }
    
    // 설정 사용
    port := cfg.Get("server.port")
}
```

## 민감한 정보 처리

- 비밀번호, API 키 등 민감한 정보는 설정 파일에 직접 포함하지 마세요.
- 민감한 정보는 환경 변수 또는 보안 저장소(Vault, Kubernetes Secrets 등)를 통해 제공하세요.
- 설정에서 민감한 정보를 참조할 때 다음 형식을 사용하세요:

```yaml
database:
  password: ${DB_PASSWORD}  # 환경 변수에서 로드
```

## 가이드라인

1. 설정 파일은 각 서비스별로 분리하세요.
2. 주석을 통해 각 설정 옵션을 문서화하세요.
3. 기본값을 제공하고 명확히 문서화하세요.
4. 설정을 변경할 때는 프로덕션 환경에 영향을 주기 전에 철저히 테스트하세요. 