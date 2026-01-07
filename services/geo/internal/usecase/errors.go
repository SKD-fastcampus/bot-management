package usecase

import "errors"

// 에러 타입 정의
var (
	ErrInvalidIPAddress    = errors.New("유효하지 않은 IP 주소입니다")
	ErrFeatureNotSupported = errors.New("지원하지 않는 기능입니다")
	ErrGeoLookupFailed     = errors.New("지리 정보 조회에 실패했습니다")
)
