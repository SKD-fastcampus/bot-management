package config

type Server struct {
	// HTTP 서버 설정
	HTTP struct {
		Port    string `yaml:"port"`
		Timeout int    `yaml:"timeout"`
		Debug   bool   `yaml:"debug"`
	} `yaml:"http"`

	// gRPC 서버 설정
	GRPC struct {
		Port    string `yaml:"port"`
		Timeout int    `yaml:"timeout"`
	} `yaml:"grpc"`
}
