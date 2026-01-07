# ────────────────────────────── 1️⃣ builder ──────────────────────────────
FROM --platform=$BUILDPLATFORM golang:1.22 AS builder
WORKDIR /workspace

# 1) go.work + 각 모듈(go.mod/go.sum)만 우선 복사 → 캐시용 레이어
COPY go.work ./

# 서비스 모듈
COPY services/*/go.mod   services/*/go.sum   ./services/

# 공통 패키지 모듈(필요한 만큼 추가)
COPY pkg/go.mod     pkg/go.sum    ./pkg/

# 2) 의존성 다운로드 (BuildKit 캐시 활용)
RUN --mount=type=cache,target=/go/pkg/mod \
    go work sync && go mod download

# ── 🆕 빌드 타임 ARG 선언 ──
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown

# ARG 값을 레이어 아래에서도 쓸 수 있게 ENV로 승격 (선택)
ENV VERSION=${VERSION} \
    COMMIT=${COMMIT} \
    BUILD_DATE=${BUILD_DATE}

# 소스 복사 후 컴파일
COPY pkg        ./pkg
COPY services   ./services

RUN cd services/geo && \
    CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH \
    go build -trimpath \
      -ldflags "-s -w \
        -X 'main.version=${VERSION}' \
        -X 'main.commit=${COMMIT}' \
        -X 'main.date=${BUILD_DATE}'" \
      -o /out/app
      
# ────────────────────────────── 2️⃣ runtime ──────────────────────────────
FROM gcr.io/distroless/static-debian11
COPY --from=builder /out/app /app
ENTRYPOINT ["/app"]
