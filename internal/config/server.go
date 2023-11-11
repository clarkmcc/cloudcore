package config

import "github.com/spf13/viper"

func init() {
	viper.MustBindEnv("port", "PORT")
	viper.SetDefault("port", 10000)
	_ = viper.BindEnv("logging.level", "LOGGING_LEVEL")
	_ = viper.BindEnv("logging.debug", "LOGGING_DEBUG")
	viper.MustBindEnv("auth.signingSecret", "AUTH_TOKEN_SIGNING_SECRET")
}

type ServerConfig struct {
	Port    int
	Logging Logging `json:"logging"`
	Auth    serverAuth
}

type serverAuth struct {
	TokenSigningSecret string `json:"tokenSigningSecret"`
}

func NewServerConfig() (*ServerConfig, error) {
	var cfg ServerConfig
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
