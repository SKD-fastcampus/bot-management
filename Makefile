.PHONY: setup run build test lint proto-gen docker-all docker-notification docker-api docker-auth mock tidy

# 기본 명령어
all: setup build

# 개발 환경 설정
setup:
	@echo "개발 환경을 설정합니다..."
	go mod download
	cd services/geo && go mod download
	cd tools && go mod download
	go install github.com/vektra/mockery/v2@latest

# 모든 모듈 tidy
tidy:
	@echo "모든 모듈의 go.mod 파일을 정리합니다..."
	cd pkg && go mod tidy
	cd services/geo && go mod tidy
	cd tools && go mod tidy
	go work sync

# 모든 서비스 실행
run:
	@echo "모든 서비스를 실행합니다..."
	docker-compose up

# 서비스 빌드
build:
	@echo "모든 서비스를 빌드합니다..."
	go build -o bin/geo services/geo/cmd/server/main.go

# 프로토콜 버퍼 코드 생성
proto-gen:
	@echo "Protocol Buffers 코드를 생성합니다..."
	protoc -I=. \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		./proto/geo/v1/*.proto

# 테스트 실행
test:
	@echo "모든 테스트를 실행합니다..."
	go test ./pkg/...
	go test ./services/geo/...

# 린트 체크
lint:
	@echo "코드 린트 검사를 실행합니다..."
	golangci-lint run ./...

# Mock 생성
mock:
	@echo "Mock 객체를 생성합니다..."
	mockery --config=.mockery.yaml

# Docker 이미지 빌드
docker-geo:
	@echo "지리 서비스 Docker 이미지를 빌드합니다..."
	docker build -t geo-service -f deployments/docker/geo.Dockerfile .

air-geo:
	APP_SERVICE=geo air -c .air.toml -build.args_bin="--config=configs/dev/geo.yaml"