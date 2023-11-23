package agent

import (
	"errors"
	"fmt"
	"github.com/clarkmcc/cloudcore/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
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

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	viper.AddConfigPath(cwd)

	err = viper.ReadInConfig()
	var e viper.ConfigFileNotFoundError
	if err != nil && !errors.As(err, &e) {
		return nil, err
	}

	psk, err := getPsk(cmd, cwd)
	if err != nil {
		return nil, fmt.Errorf("getting psk: %w", err)
	}
	if len(psk) > 0 {
		viper.SetDefault("preSharedKey", psk)
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	utils.PrintStruct(cfg)
	return &cfg, nil
}

// getPsk attempts to load the psk from any of the locations where it
// could be provided, specifically
//  1. The command line flags -- in which case we save it to a file so that
//     it doesn't need to be provided every time.
//  2. The psk file (created previously).
//
// On linux we save this to the /etc/cloudcored/psk file.
// On darwin we save this to the ~/.cloudcored/psk file.
// Windows not supported for PSK files yet.
func getPsk(cmd *cobra.Command, cwd string) (string, error) {
	psk := cmd.Flag("psk").Value.String()
	if len(psk) > 0 {
		return psk, writePskToFile(psk)
	}
	pskBytes, err := getPskFromFile(cwd)
	if err != nil {
		return "", err
	}
	return string(pskBytes), nil
}

func getPskFromFile(_ string) ([]byte, error) {
	var filename string
	switch runtime.GOOS {
	case "linux":
		filename = "/etc/cloudcored/psk"
	case "darwin":
		dir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("getting home dir: %w", err)
		}
		filename = filepath.Join(dir, ".cloudcored", "psk")
	default:
		return nil, fmt.Errorf("reading psk from file not supported on %s", runtime.GOOS)
	}
	return os.ReadFile(filename)
}

func writePskToFile(psk string) error {
	var filename string
	switch runtime.GOOS {
	case "linux":
		filename = "/etc/cloudcored/psk"
	case "darwin":
		dir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("getting home dir: %w", err)
		}
		filename = filepath.Join(dir, ".cloudcored", "psk")
	default:
		return fmt.Errorf("saving psk to file not supported on %s", runtime.GOOS)
	}
	err := os.MkdirAll(filepath.Dir(filename), 0600)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, []byte(psk), 0600)
}
