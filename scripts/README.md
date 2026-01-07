# 스크립트

이 디렉토리에는 개발 및 배포에 유용한 유틸리티 스크립트가 포함되어 있습니다.

## 주요 스크립트

- **create_multiple_dbs.sh**: Docker 컨테이너에서 여러 PostgreSQL 데이터베이스를 생성하는 스크립트
- **migrations/**: 데이터베이스 마이그레이션 스크립트

## 사용 방법

### 데이터베이스 마이그레이션

마이그레이션 스크립트를 실행하려면:

```bash
# 마이그레이션 적용
./scripts/migrations/run.sh up

# 특정 서비스에 대한 마이그레이션 적용
./scripts/migrations/run.sh up notification

# 마이그레이션 롤백
./scripts/migrations/run.sh down

# 마이그레이션 상태 확인
./scripts/migrations/run.sh status
```

### 개발 환경 초기화

개발 환경을 초기화하려면:

```bash
# 초기 설정
./scripts/setup.sh

# 개발용 인증서 생성
./scripts/generate_certs.sh
```

## 스크립트 추가 가이드라인

1. 모든 스크립트에는 사용법 설명이 포함되어야 합니다 (`--help` 플래그 지원).
2. 스크립트는 오류 처리와 로깅을 포함해야 합니다.
3. 환경 변수를 통해 구성 가능하게 만드세요.
4. 스크립트가 특정 의존성을 요구하는 경우, 시작 시 이를 확인하고 명확한 오류 메시지를 표시하세요. 