# GEMINI.md

이 문서는 **semo-backend-monorepo** 프로젝트를 이해하고, 개발하고, 유지보수하기 위한 핵심 가이드입니다. **GEMINI** 및 이 프로젝트에 참여하는 모든 AI 에이전트와 개발자는 이 문서를 기준으로 작업을 수행해야 합니다.

---

## 1. 프로젝트 개요 (Project Overview)

이 레포지토리는 **Go Workspaces**를 기반으로 한 **마이크로서비스 아키텍처(Microservices Architecture)** 모노레포입니다.
각 서비스는 독립적으로 배포 가능하며, **Clean Architecture** 원칙을 엄격하게 따릅니다.

### 핵심 철학
*   **독립성**: 각 서비스는 서로의 내부 구현을 알지 못하며, gRPC/HTTP로만 통신합니다.
*   **일관성**: 모든 서비스는 동일한 프로젝트 구조와 코딩 컨벤션을 따릅니다.
*   **명시성**: 의존성은 암시적이지 않고 명시적으로 주입됩니다(Manual Dependency Injection).

### 현재 서비스 목록 (Applications)
| 서비스명 | 설명 | 상태 | 위치 |
| :--- | :--- | :--- | :--- |
| **geo** | 위치 기반 서비스 (GeoLite2 활용) | 운영중 | `services/geo` |

---

## 2. 아키텍처 및 디자인 패턴 (Architecture & Design Patterns)

이 프로젝트는 **Clean Architecture (Hexagonal Architecture)**를 따릅니다. 코드를 작성할 때 다음 계층 구조와 의존성 방향을 반드시 준수해야 합니다.

### 2.1 계층 구조 (Layered Structure)

의존성 방향: **Infrastructure → Adapter → Usecase → Domain**
(내부 계층은 외부 계층을 알지 못합니다.)

1.  **Domain (`internal/domain`)**
    *   **역할**: 핵심 비즈니스 로직과 엔티티 정의.
    *   **구성**: Entity(Struct), Repository Interface.
    *   **특징**: 외부 패키지(DB 드라이버, 프레임워크 등)에 의존하지 않는 순수 Go 코드.

2.  **Usecase (`internal/usecase`)**
    *   **역할**: 비즈니스 흐름 제어 (Application Business Rules).
    *   **구성**: Usecase Interface 및 구현체.
    *   **특징**: Domain Repository Interface를 사용하여 데이터를 처리.

3.  **Adapter (`internal/adapter`)**
    *   **역할**: 외부 세계와 유스케이스 간의 변환기.
    *   **구성**:
        *   **Handler**: HTTP(Echo), gRPC 요청을 받아 Usecase 호출.
        *   **Repository**: Domain Repository Interface의 구현체 (DB 접근).

4.  **Infrastructure (`internal/infrastructure`)**
    *   **역할**: 기술적인 세부 사항 구현.
    *   **구성**: HTTP 서버(Echo), gRPC 서버, DB 연결 설정 등.
    *   **특징**: 프레임워크나 드라이버에 직접적으로 의존하는 코드가 위치.

### 2.2 의존성 주입 (Dependency Injection)

우리는 **수동 의존성 주입(Manual Dependency Injection)**을 사용합니다. `uber-go/fx` 같은 DI 프레임워크를 사용하지 **않습니다**.

*   **진입점**: 모든 의존성 조립은 `cmd/server/main.go`에서 이루어집니다.
*   **생성자 주입**: 모든 struct는 생성자(`New...`)를 통해 의존성을 주입받아야 합니다.

**`main.go` Wiring 예시:**
```go
// 1. Repository 생성 (Infrastructure -> Adapter)
repo := repository.NewPostgresRepository(dbConn)

// 2. Usecase 생성 (Adapter -> Usecase)
usecase := usecase.NewUserUsecase(repo)

// 3. Handler 생성 (Usecase -> Adapter)
handler := httpHandler.NewUserHandler(usecase)

// 4. Server에 Handler 등록
server.Register(handler)
```

---

## 3. 프로젝트 구조 (Project Structure)

### `/pkg` (Shared Packages)
모든 서비스에서 공통으로 사용되는 유틸리티 라이브러리입니다.
*   `logger`: Zap 기반 구조화된 로깅. `logger.DefaultZapLogger()` 사용.
*   `config`: Viper 기반 설정 관리.
*   `errors`, `middleware` 등.

### `/configs` (Configurations)
환경별 설정 파일(`config.yaml`)을 관리합니다.
*   구조: `/configs/{service_name}/{environment}/config.yaml` (예: `configs/geo/dev/config.yaml`) (실제 구조는 프로젝트별로 상이할 수 있으나 기본적으로 `dev`, `prod` 등으로 나뉨)
*   **주의**: 비밀번호 등 민감 정보는 환경변수(`${ENV_VAR}`)로 대체해야 합니다.

### `/deployments`
*   `docker/`: 서비스별 Dockerfile. `distroless` 이미지를 사용하여 경량화.
*   `k8s/helm/`: Helm 차트.

### `/proto`
Protocol Buffers 정의 파일(`.proto`)이 위치합니다.
*   변경 시 `make proto-gen`으로 Go 코드를 재생성해야 합니다.

---

## 4. 개발 워크플로우 (Development Workflow)

### 4.1 새 기능 추가 시 (Feature Development)
1.  **Domain 정의**: `internal/domain`에 엔티티와 인터페이스 작성.
2.  **Usecase 구현**: `internal/usecase`에 비즈니스 로직 구현 및 테스트 작성.
3.  **Infrastructure/Adapter 구현**: `internal/adapter`에 Repository 구현 및 Handler 작성.
4.  **Wiring**: `cmd/server/main.go`에 의존성 연결.

### 4.2 새 서비스 추가 시 (New Service)
1.  `services/{service_name}` 디렉토리 생성.
2.  기존 `services/geo` 구조를 복사하여 템플릿으로 활용.
3.  `go.work` 파일에 새 모듈 경로 추가 (`go work use ./services/{service_name}`).
4.  `Makefile`에 빌드 및 실행 명령어 추가.

### 4.3 코딩 컨벤션 (Coding Standards)
*   **로깅**: `log.Print` 대신 `pkg/logger` 사용. 항상 `ctx`를 포함하거나 구조화된 필드(`zap.String` 등) 사용.
*   **에러 핸들링**: 에러를 무시하지 말고 `pkg/errors`나 `fmt.Errorf`로 래핑하여 문맥 정보 추가.
*   **테스트**: 비즈니스 로직(Usecase)은 100% 단위 테스트 커버리지를 목표로 함. `mockery`로 모킹.

### 4.4 데이터베이스 관리 (Database & Migration)
*   **마이그레이션 전략**: **GORM AutoMigrate**를 사용합니다.
    *   별도의 마이그레이션 도구(flyway 등)를 사용하지 않고, 애플리케이션 시작 시 `gorm.AutoMigrate(&Entity{})`를 호출하여 스키마를 동기화합니다.
    *   **주의**: 컬럼 삭제나 데이터 타입 변경 같은 파괴적인 변경은 주의가 필요하며, 필요 시 수동으로 SQL을 실행해야 할 수 있습니다.
*   **초기화**: `scripts/create_multiple_dbs.sh` 스크립트를 통해 로컬 개발용 DB를 생성할 수 있습니다.

---

## 5. 명령어 및 운영 가이드 (Operational Guide)

루트 디렉토리의 `Makefile`을 통해 주요 작업을 수행합니다.

### 개발 (Development)
```bash
make setup          # 개발 환경 설정
make run            # Docker Compose로 모든 서비스 실행
make air-geo        # Geo 서비스 Hot Reload 실행
```

### 빌드 및 배포 (Build & Deploy)
```bash
make build          # 모든 서비스 빌드
make docker-geo     # Geo 서비스 Docker 이미지 빌드
```

### 코드 품질 및 생성 (Quality & Gen)
```bash
make lint           # golangci-lint 실행
make test           # 모든 테스트 실행
make proto-gen      # gRPC 코드 생성
make mock           # Mockery로 Mock 생성
make tidy           # Go 모듈 정리
```

### 5.1 필수 도구 (Prerequisites)
프로젝트 개발 및 빌드를 위해 다음 도구들이 반드시 설치되어 있어야 합니다.

1.  **Go** (1.23+)
2.  **Docker & Docker Compose**
3.  **Protocol Buffers Implements**
    *   `.proto` 파일 수정 및 코드 생성을 위해 필요합니다.
    *   **설치 (macOS)**:
        ```bash
        brew install protobuf
        go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
        go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
        ```
    *   **확인**: `protoc --version` 실행 시 버전 정보가 출력되어야 합니다.


---

## 6. 참고 리소스 (References)
*   **Reference Implementation**: `services/geo` 디렉토리가 이 프로젝트의 표준 구현체입니다. 구조가 헷갈릴 때는 이 서비스를 참고하세요.
*   **Dependency Injection Rules**: `.cursor/rules/2-dependency-injection-rule.mdc`
*   **Documentation**: `CLAUDE.md`, `configs/README.md`
