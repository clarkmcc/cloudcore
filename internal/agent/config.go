package agent

import (
	"github.com/clarkmcc/cloudcore/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func init() {
	viper.SetDefault("server.endpoint", "127.0.0.1:10000")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.debug", true)
	viper.SetDefault("database.flavor", databaseFlavorMemory)
}

type databaseFlavor string

const (
	databaseFlavorMemory databaseFlavor = "memory"
)

type Config struct {
	Server       server
	Logging      Logging
	Database     database
	PreSharedKey string `yaml:"preSharedKey"`
}

type database struct {
	Flavor databaseFlavor
}

type server struct {
	Endpoint string
}

type Logging struct {
	Level string `json:"level"`
	Debug bool   `json:"debug"`
}

func NewConfig(cmd *cobra.Command) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if cwd, err := os.Getwd(); err == nil {
		viper.AddConfigPath(cwd)
	}
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	psk := cmd.Flag("psk").Value.String()
	if len(psk) > 0 {
		viper.Set("preSharedKey", psk)
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	utils.PrintStruct(cfg)
	return &cfg, nil
}
