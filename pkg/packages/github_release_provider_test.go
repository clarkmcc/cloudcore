package packages

import (
	"context"
	"github.com/clarkmcc/cloudcore/pkg/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestGithubReleaseProvider_FindLatestPackage(t *testing.T) {
	provider := NewGithubReleaseProvider(GithubReleaseProviderConfig{
		Owner: "clarkmcc",
		Repo:  "cloudcore",
	}, zap.NewNop())

	packages, err := provider.GetLatestPackages(context.Background())
	assert.NoError(t, err)
	assert.Greater(t, len(packages), 0)
	utils.PrintStruct(packages)

	pack, err := provider.FindLatestPackage(context.Background(), "linux", "amd64", "")
	assert.NoError(t, err)
	assert.NotNil(t, pack)
	utils.PrintStruct(pack)
}
