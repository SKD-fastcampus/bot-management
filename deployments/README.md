# 배포

이 디렉토리에는 서비스 배포에 필요한 파일과 설정이 포함되어 있습니다.

## 디렉토리 구조

- **docker/**: Docker 관련 파일 (Dockerfile 등)
- **k8s/**: Kubernetes 매니페스트 파일

## Docker

`docker` 디렉토리에는 각 서비스에 대한 Dockerfile이 포함되어 있습니다:

- **notification.Dockerfile**: 알림 서비스 Dockerfile
- **api.Dockerfile**: API 서비스 Dockerfile
- **auth.Dockerfile**: 인증 서비스 Dockerfile

Docker 이미지를 빌드하려면:

```bash
# 알림 서비스 이미지 빌드
docker build -t notification-service -f deployments/docker/notification.Dockerfile .

# API 서비스 이미지 빌드
docker build -t api-service -f deployments/docker/api.Dockerfile .

# 인증 서비스 이미지 빌드
docker build -t auth-service -f deployments/docker/auth.Dockerfile .

# 또는 Makefile 사용
make docker-all
```

## Kubernetes

`k8s` 디렉토리에는 Kubernetes 배포를 위한 매니페스트 파일이 포함되어 있습니다:

- **base/**: 기본 매니페스트 (모든 환경에 공통)
- **overlays/**: 환경별 오버레이 (kustomize 사용)
  - **dev/**: 개발 환경
  - **staging/**: 스테이징 환경
  - **prod/**: 프로덕션 환경

Kubernetes에 배포하려면:

```bash
# 개발 환경에 배포
kubectl apply -k deployments/k8s/overlays/dev

# 스테이징 환경에 배포
kubectl apply -k deployments/k8s/overlays/staging

# 프로덕션 환경에 배포
kubectl apply -k deployments/k8s/overlays/prod
```

## 환경별 배포 구성

각 환경에는 다음과 같은 구성이 포함됩니다:

- **dev**: 단일 레플리카, 디버그 모드, 메모리/CPU 제한 낮음
- **staging**: 2-3 레플리카, 프로덕션과 유사한 설정, 분리된 네임스페이스
- **prod**: 고가용성 구성, 오토스케일링, 리소스 제한, 인그레스 설정

## CI/CD 파이프라인

CI/CD 파이프라인은 GitHub Actions를 사용하여 구현되어 있으며, `.github/workflows` 디렉토리에 정의되어 있습니다. 