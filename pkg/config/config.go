// Package config는 애플리케이션 설정을 관리하는 패키지입니다.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config 인터페이스는 설정 값에 액세스하기 위한 메서드를 정의합니다.
type Config interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetFloat64(key string) float64
	GetStringSlice(key string) []string
	GetStringMap(key string) map[string]interface{}
	GetAll() map[string]interface{}
}

// viperConfig는 viper를 사용하여 Config 인터페이스를 구현합니다.
type viperConfig struct {
	v *viper.Viper
}

// GetString은 문자열 설정 값을 반환합니다.
func (c *viperConfig) GetString(key string) string {
	return c.v.GetString(key)
}

// GetInt는 정수 설정 값을 반환합니다.
func (c *viperConfig) GetInt(key string) int {
	return c.v.GetInt(key)
}

// GetBool은 불리언 설정 값을 반환합니다.
func (c *viperConfig) GetBool(key string) bool {
	return c.v.GetBool(key)
}

// GetFloat64는 부동 소수점 설정 값을 반환합니다.
func (c *viperConfig) GetFloat64(key string) float64 {
	return c.v.GetFloat64(key)
}

// GetStringSlice는 문자열 슬라이스 설정 값을 반환합니다.
func (c *viperConfig) GetStringSlice(key string) []string {
	return c.v.GetStringSlice(key)
}

// GetStringMap은 맵 설정 값을 반환합니다.
func (c *viperConfig) GetStringMap(key string) map[string]interface{} {
	return c.v.GetStringMap(key)
}

// GetAll은 전체 설정을 맵으로 반환합니다.
func (c *viperConfig) GetAll() map[string]interface{} {
	return c.v.AllSettings()
}

// 설정 디렉토리 경로
const configDir = "configs"

// Load는 지정된 서비스 이름에 해당하는 설정 파일을 로드합니다.
func Load(serviceName string) (Config, error) {
	v := viper.New()

	// 환경 변수 설정
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev" // 기본 환경은 dev
	}

	// 설정 파일 확장자 및 유형 설정
	v.SetConfigType("yaml")

	// 환경 변수 바인딩 설정
	v.SetEnvPrefix(strings.ToUpper(serviceName))
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 설정 파일 경로 설정
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		// 기본 경로는 현재 디렉토리의 configs/{env}/{service}.yaml
		configPath = filepath.Join(configDir, env)
	}

	// 설정 파일 이름 설정
	configName := serviceName
	v.SetConfigName(configName)
	v.AddConfigPath(configPath)

	// 설정 파일 로드
	if err := v.ReadInConfig(); err != nil {
		// configs/example 디렉토리에서 예제 설정 파일 시도
		v.SetConfigName(configName)
		v.AddConfigPath(filepath.Join(configDir, "example"))
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("설정 파일 로드 실패: %w", err)
		}
	}

	return &viperConfig{v: v}, nil
}
