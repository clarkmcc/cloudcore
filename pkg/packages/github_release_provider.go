package packages

import (
	"context"
	"errors"
	"github.com/google/go-github/v56/github"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

const agentBinaryName = "cloudcored"

var _ Provider = (*GithubReleaseProvider)(nil)

type GithubReleaseProviderConfig struct {
	Owner string
	Repo  string
}

type GithubReleaseProvider struct {
	*github.Client
	logger *zap.Logger

	owner string
	repo  string
}

func (g *GithubReleaseProvider) GetLatestPackages(ctx context.Context) ([]Package, error) {
	rel, res, err := g.Repositories.GetLatestRelease(ctx, g.owner, g.repo)
	if res.StatusCode == http.StatusNotFound {
		return nil, ErrNoPackages
	}
	if err != nil {
		return nil, err
	}
	var packages []Package
	for _, ga := range rel.Assets {
		if ga.BrowserDownloadURL == nil {
			continue
		}
		a, err := parseAssetName(ga.GetName())
		if err != nil {
			g.logger.Warn("parsing asset name", zap.Error(err), zap.String("name", ga.GetName()))
			continue
		}
		if a.Name != agentBinaryName {
			continue
		}
		packages = append(packages, Package{
			GOOS:        a.GOOS,
			GOARCH:      a.GOARCH,
			GOARM:       a.GOARM,
			Version:     a.Version,
			DownloadURL: *ga.BrowserDownloadURL,
		})
	}
	return packages, nil
}

func (g *GithubReleaseProvider) FindLatestPackage(ctx context.Context, goos, goarch, goarm string) (*Package, error) {
	packages, err := g.GetLatestPackages(ctx)
	if err != nil {
		return nil, err
	}
	for _, p := range packages {
		if p.GOOS == goos && p.GOARCH == goarch && p.GOARM == goarm {
			return &p, nil
		}
	}
	return nil, ErrNoPackages
}

func NewGithubReleaseProvider(config GithubReleaseProviderConfig, logger *zap.Logger) *GithubReleaseProvider {
	return &GithubReleaseProvider{
		Client: github.NewClient(nil),
		owner:  config.Owner,
		repo:   config.Repo,
		logger: logger,
	}
}

type asset struct {
	Name    string
	Version string
	GOOS    string
	GOARCH  string
	GOARM   string
}

func parseAssetName(name string) (*asset, error) {
	parts := strings.Split(trimSuffix(name), "_")
	if len(parts) < 4 {
		return nil, errors.New("invalid asset name")
	}
	var goarm string
	if len(parts) >= 5 {
		goarm = parts[4]
	}
	return &asset{
		Name:    parts[0],
		Version: parts[1],
		GOOS:    parts[2],
		GOARCH:  parts[3],
		GOARM:   goarm,
	}, nil
}

func trimSuffix(name string) string {
	name = strings.TrimSuffix(name, ".deb")
	name = strings.TrimSuffix(name, ".apk")
	name = strings.TrimSuffix(name, ".rpm")
	name = strings.TrimSuffix(name, ".exe")
	return name
}
