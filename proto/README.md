# 프로토콜 버퍼 정의

이 디렉토리에는 서비스 간 통신에 사용되는 Protocol Buffers 정의 파일이 포함되어 있습니다. 여기서 정의된 proto 파일은 gRPC 서비스 및 메시지 구조를 정의합니다.

## 디렉토리 구조

- **auth/**: 인증 관련 프로토콜 정의
- **notification/**: 알림 관련 프로토콜 정의
- **api/**: API 서비스 관련 프로토콜 정의

## 사용 방법

### 프로토콜 파일 작성

새로운 프로토콜 파일을 작성할 때는 다음 형식을 따라주세요:

```protobuf
syntax = "proto3";

package semo.service_name.v1;

option go_package = "github.com/wekeepgrowing/semo-backend-monorepo/proto/service_name/v1;service_namev1";

service ServiceName {
  rpc MethodName(RequestMessage) returns (ResponseMessage) {}
}

message RequestMessage {
  string field1 = 1;
  int32 field2 = 2;
}

message ResponseMessage {
  string result = 1;
}
```

### 코드 생성

프로토콜 파일에서 Go 코드를 생성하려면 다음 명령어를 사용하세요:

```bash
# 루트 디렉토리에서
make proto-gen

# 또는 직접 실행
protoc -I=. \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  ./proto/service_name/v1/*.proto
```

## 버전 관리

- 프로토콜 정의는 디렉토리 구조에 버전을 포함합니다 (예: `v1`, `v2`).
- 기존 서비스의 호환성을 깨뜨리는 변경은 새로운 버전 디렉토리에 정의해야 합니다.
- 이전 버전의 지원 계획을 명확히 문서화하세요.

## 가이드라인

1. 필드 이름과 메시지 이름은 일관된 명명 규칙을 따라야 합니다 (camelCase 또는 snake_case).
2. 각 서비스와 메시지에 주석을 통한 문서화를 제공하세요.
3. 열거형(enum)은 적절한 접두사와 함께 명확하게 정의하세요.
4. 불필요한 중복 메시지를 피하고 공통 메시지를 재사용하세요. 