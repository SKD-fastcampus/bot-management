package repository

import (
	"context"
	"time"
)

// CacheRepository 캐시 관련 저장소 인터페이스
type CacheRepository interface {
	// Get 키로 값 조회
	Get(ctx context.Context, key string) ([]byte, error)

	// Set 키-값 저장
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error

	// Delete 키 삭제
	Delete(ctx context.Context, key string) error

	// Exists 키 존재 여부 확인
	Exists(ctx context.Context, key string) (bool, error)

	// Expire 키 만료 시간 설정
	Expire(ctx context.Context, key string, expiration time.Duration) error
}
