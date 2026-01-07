# Semo 백엔드 Helm 차트

Semo 백엔드 모노레포 애플리케이션을 Kubernetes에 배포하기 위한 Helm 차트입니다.

## 개요

이 Helm 차트는 Semo 백엔드 애플리케이션 구성 요소를 Kubernetes에 배포하기 위한 템플릿을 제공합니다. 모노레포 구조를 반영하여 각 서비스는 독립적인 서브차트로 구성되어 있으며, 필요에 따라 개별적으로 배포하거나 전체를 함께 배포할 수 있습니다.

### 서브차트 구성

- **api**: Semo API 서비스
- **auth**: 인증 서비스
- **admin**: 관리자 대시보드 서비스
- **worker**: 백그라운드 작업 처리 서비스

## 사용 방법

### 전체 서비스 배포

전체 백엔드 스택을 한 번에 배포하려면:

```bash
# 기본 설정으로 배포
helm install semo ./deployments/k8s/helm

# 사용자 정의 values 파일 사용
helm install semo ./deployments/k8s/helm -f my-values.yaml

# 특정 네임스페이스에 배포
helm install semo ./deployments/k8s/helm --namespace semo --create-namespace
```

### 특정 서비스만 배포

특정 서비스만 배포하려면 values 파일에서 필요한 서비스만 활성화하면 됩니다:

```bash
# API와 인증 서비스만 배포
helm install semo ./deployments/k8s/helm --set admin.enabled=false,worker.enabled=false
```

또는 서브차트만 직접 배포할 수도 있습니다:

```bash
# 인증 서비스만 배포
helm install auth-service ./deployments/k8s/helm/charts/auth
```

### 배포 업그레이드

기존 배포를 업그레이드하려면:

```bash
helm upgrade semo ./deployments/k8s/helm
```

### 배포 삭제

배포를 삭제하려면:

```bash
helm uninstall semo
```

## 설정

### 글로벌 설정

`values.yaml` 파일에서 모든 서브차트에 적용되는 글로벌 설정을 정의할 수 있습니다:

```yaml
global:
  imageRegistry: "my-registry.com"
  imagePullSecrets:
    - regcred
  storageClass: "standard"
  environment: dev
  labels:
    org: semo
    project: backend
```

### 서비스 활성화/비활성화

각 서비스는 개별적으로 활성화하거나 비활성화할 수 있습니다:

```yaml
api:
  enabled: true
auth:
  enabled: true
admin:
  enabled: false
worker:
  enabled: false
```

### 서비스별 설정

각 서브차트에는 자체 values.yaml 파일이 있어 서비스별 설정을 정의할 수 있습니다.

#### auth 서비스 설정 예시

```yaml
replicaCount: 2
image:
  tag: v1.0.1
service:
  type: ClusterIP
ingress:
  enabled: true
  hosts:
    - host: auth.example.com
      paths:
        - path: /
          pathType: Prefix
```

## 환경별 설정

각 환경(개발, 스테이징, 프로덕션)에 맞게 values 파일을 별도로 관리하는 것을 권장합니다:

- `values-dev.yaml`
- `values-staging.yaml`
- `values-prod.yaml`

예를 들어, 프로덕션 환경에 배포할 때:

```bash
helm install semo ./deployments/k8s/helm -f values-prod.yaml
```

## FAQ

### Q: 서브차트 간의 통신은 어떻게 설정하나요?

A: 서브차트 간 통신은 Kubernetes 서비스를 통해 이루어집니다. 예를 들어, API 서비스가 인증 서비스의 gRPC를 호출해야 한다면:

```yaml
# api 서비스의 values.yaml
config:
  services:
    auth:
      grpcUrl: semo-auth:8082
```

### Q: 시크릿과 민감한 정보는 어떻게 관리하나요?

A: 개발 환경에서는 values 파일에 직접 시크릿 값을 포함할 수 있지만, 프로덕션 환경에서는 외부 시크릿 관리자(예: HashiCorp Vault, AWS Secrets Manager)를 사용하거나 Kubernetes Secrets를 별도로 관리하는 것이 좋습니다.

Helm 차트에서는 다음과 같이 외부 시크릿을 참조할 수 있습니다:

```yaml
auth:
  secrets:
    existingSecret: "auth-credentials"
```

### Q: 데이터베이스와 같은 외부 서비스는 어떻게 설정하나요?

A: 외부 서비스는 일반적으로 별도의 Helm 차트나 관리형 서비스로 제공됩니다. 예를 들어 PostgreSQL을 사용한다면:

```yaml
# PostgreSQL을 의존성으로 추가
dependencies:
  - name: postgresql
    version: 10.16.2
    repository: https://charts.bitnami.com/bitnami
    condition: postgresql.enabled

# PostgreSQL 비활성화하고 외부 인스턴스 사용
postgresql:
  enabled: false

# 외부 데이터베이스 설정
externalDatabase:
  host: my-postgres.example.com
  port: 5432
  user: postgres
  password: postgres
  database: semo
```

### Q: Ingress 설정은 어떻게 하나요?

A: 각 서브차트에는 Ingress 설정을 위한 템플릿이 포함되어 있습니다. 활성화하려면:

```yaml
ingress:
  enabled: true
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: api.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: api-tls
      hosts:
        - api.example.com
```

### Q: 로그와 모니터링은 어떻게 설정하나요?

A: 서비스 로그는 영구 볼륨에 저장하거나 EFK(Elasticsearch, Fluentd, Kibana) 또는 ELK(Elasticsearch, Logstash, Kibana) 스택을 사용하여 중앙 집중식으로 관리할 수 있습니다. 모니터링은 Prometheus와 Grafana를 사용하는 것이 일반적입니다.

기본 로그 설정:

```yaml
volumes:
  logs:
    enabled: true
    size: 10Gi
``` 