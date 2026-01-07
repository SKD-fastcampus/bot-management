package config

import (
	"github.com/SKD-fastcampus/bot-management/pkg/config"
	"github.com/SKD-fastcampus/bot-management/pkg/logger"
	"go.uber.org/zap"
)

// Config 인증 서비스 설정 구조체
type Config struct {
	Service Service `yaml:"service"`
	Server  Server  `yaml:"server"`
	GeoLite GeoLite `yaml:"geolite"`
	JWT     JWT     `yaml:"jwt"`
	Log     Log     `yaml:"log"`
	Email   Email   `yaml:"email"`
	Logger  *zap.Logger
}

var (
	// AppConfig는 어플리케이션 전체에서 사용하는 설정 인스턴스입니다.
	AppConfig *Config
)

// Load 설정 파일 로드
func Load() (*Config, error) {
	// pkg/config 패키지를 사용하여 설정 파일 로드
	cfg, err := config.Load("geo")
	if err != nil {
		return nil, err
	}

	// Config 구조체 생성
	appConfig := &Config{}

	// 서비스 정보
	appConfig.Service.Name = cfg.GetString("service.name")
	appConfig.Service.Version = cfg.GetString("service.version")
	appConfig.Service.BaseURL = cfg.GetString("service.base_url")

	// HTTP 서버 설정
	appConfig.Server.HTTP.Port = cfg.GetString("server.port")
	appConfig.Server.HTTP.Timeout = cfg.GetInt("server.timeout")
	appConfig.Server.HTTP.Debug = cfg.GetBool("server.debug")

	// gRPC 서버 설정
	appConfig.Server.GRPC.Port = cfg.GetString("server.grpc.port")
	appConfig.Server.GRPC.Timeout = cfg.GetInt("server.grpc.timeout")

	// GeoLite 설정
	appConfig.GeoLite.DbPath = cfg.GetString("geolite.db_path")

	// JWT 설정
	appConfig.JWT.Secret = cfg.GetString("jwt.secret")
	appConfig.JWT.PrivateKey = cfg.GetString("jwt.private_key")
	appConfig.JWT.PublicKey = cfg.GetString("jwt.public_key")
	appConfig.JWT.AccessTokenExpiry = cfg.GetInt("jwt.access_token_expiry")
	appConfig.JWT.RefreshTokenExpiry = cfg.GetInt("jwt.refresh_token_expiry")

	// 이메일 설정
	appConfig.Email.SenderEmail = cfg.GetString("email.sender_email")
	appConfig.Email.SMTPHost = cfg.GetString("email.smtp_host")
	appConfig.Email.SMTPPort = cfg.GetInt("email.smtp_port")
	appConfig.Email.SMTPUser = cfg.GetString("email.smtp_user")
	appConfig.Email.SMTPPass = cfg.GetString("email.smtp_pass")

	// 로그 설정
	appConfig.Log.Level = cfg.GetString("log.level")
	appConfig.Log.Format = cfg.GetString("log.format")
	appConfig.Log.Output = cfg.GetString("log.output")

	// 로거 설정
	loggerConfig := logger.Config{
		Level:       appConfig.Log.Level,
		Format:      appConfig.Log.Format,
		Output:      appConfig.Log.Output,
		Development: appConfig.Server.HTTP.Debug,
	}

	// 로거 생성
	appConfig.Logger, err = logger.NewZapLogger(loggerConfig)
	if err != nil {
		return nil, err
	}

	// 전역 변수에 설정
	AppConfig = appConfig

	return appConfig, nil
}
