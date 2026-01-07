//go:build tools

package tools

import (
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
)

// 이 파일은 도구 버전 관리를 위해 사용됩니다.
// go.mod에 개발에 필요한 도구 의존성을 포함하도록 합니다.
// 이 파일은 실제로 컴파일되지 않습니다.
