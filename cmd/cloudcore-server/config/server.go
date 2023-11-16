package config

import "github.com/spf13/viper"

func init() {
	viper.MustBindEnv("agentServer.port", "AGENT_SERVER_PORT")
	viper.MustBindEnv("appServer.port", "APP_SERVER_PORT")
	viper.SetDefault("agentServer.port", 10000)
	viper.SetDefault("appServer.port", 10001)
	_ = viper.BindEnv("logging.level", "LOGGING_LEVEL")
	_ = viper.BindEnv("logging.debug", "LOGGING_DEBUG")
	viper.MustBindEnv("auth.signingSecret", "AUTH_TOKEN_SIGNING_SECRET")
	viper.SetDefault("database.type", "cockroachdb")
	viper.SetDefault("database.name", "cloudcore")
	viper.MustBindEnv("database.connectionString", "DATABASE_CONNECTION_STRING")
}

type Config struct {
	AgentServer AgentServer    `json:"agentServer"`
	AppServer   AppServer      `json:"appServer"`
	Logging     Logging        `json:"logging"`
	Auth        serverAuth     `json:"auth"`
	Database    serverDatabase `json:"database"`
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

type Logging struct {
	Level string `json:"level"`
	Debug bool   `json:"debug"`
}

type AgentServer struct {
	Port int `json:"port"`
}

type AppServer struct {
	Port int `json:"port"`
}

func New() (*Config, error) {
	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
