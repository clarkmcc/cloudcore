package packages

import (
	"context"
	"errors"
)

type Package struct {
	GOOS        string `json:"goos"`
	GOARCH      string `json:"goarch"`
	GOARM       string `json:"goarm"`
	Version     string `json:"version"`
	DownloadURL string `json:"downloadUrl"`
}

var ErrNoPackages = errors.New("no packages found")

type Provider interface {
	// GetLatestPackages returns all the available packages for the latest release and
	// returns ErrNoPackages if no packages can be found.
	GetLatestPackages(ctx context.Context) ([]Package, error)

	// FindLatestPackage returns the latest package for the given GOOS, GOARCH, and GOARM
	// and returns ErrNoPackages if no package can be found.
	FindLatestPackage(ctx context.Context, goos, goarch, goarm string) (*Package, error)
}
