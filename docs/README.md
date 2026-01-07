# 프로젝트 문서

이 디렉토리에는 프로젝트와 관련된 문서가 포함되어 있습니다. 개발자와 사용자를 위한 가이드, 아키텍처 설계 문서, API 명세 등을 제공합니다.

## 디렉토리 구조

- **architecture/**: 시스템 아키텍처 문서
- **api/**: API 명세 및 설명서
- **guidelines/**: 개발 가이드라인 및 모범 사례

## 아키텍처 문서

`architecture/` 디렉토리에는 시스템 설계에 관한 문서가 포함되어 있습니다:

- **overview.md**: 시스템 전체 구조에 대한 개요
- **services.md**: 각 서비스에 대한 설명 및 책임
- **data-flow.md**: 데이터 흐름 및 서비스 간 통신
- **deployment.md**: 배포 구조 및 환경 설정

## API 문서

`api/` 디렉토리에는 외부 API에 대한 명세와 사용 예제가 포함되어 있습니다:

- **endpoints.md**: API 엔드포인트 목록 및 매개변수
- **authentication.md**: 인증 및 권한 부여 방법
- **examples/**: API 호출 예제 (curl, httpie 등)
- **swagger/**: Swagger 문서 (OpenAPI 명세)

## 개발 가이드라인

`guidelines/` 디렉토리에는 개발 팀을 위한 가이드라인이 포함되어 있습니다:

- **CONTRIBUTING.md**: 기여 가이드
- **code-style.md**: 코드 스타일 가이드
- **testing.md**: 테스트 작성 가이드
- **pr-process.md**: PR 프로세스 및 리뷰 기준

## 문서 갱신하기

- 코드를 변경할 때 관련 문서도 함께 갱신해주세요.
- 마크다운(.md) 형식을 사용하여 문서를 작성해주세요.
- 복잡한 다이어그램은 [Mermaid](https://mermaid-js.github.io/mermaid/#/) 또는 [PlantUML](https://plantuml.com/)을 사용하여 작성하세요.
- 코드 예제는 실제 작동하는 코드를 포함하고 주석을 통해 설명하세요. 