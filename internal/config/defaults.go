package config

import "github.com/spf13/viper"

func init() {
	viper.SetDefault("server.endpoint", "127.0.0.1:10000")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.debug", true)
	viper.SetDefault("database.flavor", AgentDatabaseFlavorMemory)
}
