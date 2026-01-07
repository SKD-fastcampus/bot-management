package config

type Service struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	BaseURL string `yaml:"base_url"`
}
