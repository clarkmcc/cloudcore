package config

import "github.com/spf13/viper"

func init() {
	viper.MustBindEnv("port", "PORT")
	viper.SetDefault("port", 10000)
	_ = viper.BindEnv("logging.level", "LOGGING_LEVEL")
	_ = viper.BindEnv("logging.debug", "LOGGING_DEBUG")
	viper.MustBindEnv("auth.signingSecret", "AUTH_TOKEN_SIGNING_SECRET")
	viper.SetDefault("database.type", "cockroachdb")
	viper.SetDefault("database.name", "cloudcore")
	viper.MustBindEnv("database.connectionString", "DATABASE_CONNECTION_STRING")
}

type ServerConfig struct {
	Port     int
	Logging  Logging        `json:"logging"`
	Auth     serverAuth     `json:"auth"`
	Database serverDatabase `json:"database"`
}

type serverAuth struct {
	TokenSigningSecret string `json:"tokenSigningSecret"`
}

type DatabaseType string

const (
	DatabaseTypeMemory      DatabaseType = "memory"
	DatabaseTypeCockroachDB DatabaseType = "cockroachdb"
)

type serverDatabase struct {
	Type             DatabaseType `json:"type"`
	Name             string       `json:"name"`
	ConnectionString string       `json:"connectionString"`
}

func NewServerConfig() (*ServerConfig, error) {
	var cfg ServerConfig
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	printStructure(cfg)
	return &cfg, nil
}
