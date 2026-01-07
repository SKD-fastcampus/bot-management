package config

type JWT struct {
	Secret             string `yaml:"secret"`
	PrivateKey         string `yaml:"private_key"`
	PublicKey          string `yaml:"public_key"`
	AccessTokenExpiry  int    `yaml:"access_token_expiry"`
	RefreshTokenExpiry int    `yaml:"refresh_token_expiry"`
}
