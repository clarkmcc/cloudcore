package main

import (
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/config"
	"github.com/clarkmcc/cloudcore/pkg/packages"
)

// componentConfigs is a fx.Provider function that returns component-specific configurations
// that are part of the global configuration.
func componentConfigs(config *config.Config) packages.GithubReleaseProviderConfig {
	return config.PackageManagement.GithubRelease
}
