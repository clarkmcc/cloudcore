package config

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func init() {
	viper.SetDefault("server.endpoint", "127.0.0.1:10000")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.debug", true)
	viper.SetDefault("database.flavor", AgentDatabaseFlavorMemory)
}

type AgentDatabaseFlavor string

const (
	AgentDatabaseFlavorMemory AgentDatabaseFlavor = "memory"
)

type Config struct {
	Server       server
	Logging      Logging
	Database     database
	PreSharedKey string `yaml:"preSharedKey"`
}

type database struct {
	Flavor AgentDatabaseFlavor
}

type server struct {
	Endpoint string
}

type Logging struct {
	Level string `json:"level"`
	Debug bool   `json:"debug"`
}

func New() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if cwd, err := os.Getwd(); err == nil {
		viper.AddConfigPath(cwd)
	}
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	printStructure(cfg)
	return &cfg, nil
}

func printStructure(v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}