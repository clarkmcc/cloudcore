package appbackend

import (
	"fmt"
	"github.com/graphql-go/graphql"
)

var packageType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Package",
	Fields: graphql.Fields{
		"goos":    &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"goarch":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"version": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"goarm":   &graphql.Field{Type: graphql.String},
	},
})

var listPackages = &graphql.Field{
	Type: graphql.NewList(packageType),
	Args: graphql.FieldConfigArgument{},
	Resolve: wrapper(func(rctx resolveContext[any]) (any, error) {
		return rctx.packages.GetLatestPackages(rctx)
	}),
}

var buildDeployAgentCommand = &graphql.Field{
	Type: graphql.NewNonNull(graphql.String),
	Args: graphql.FieldConfigArgument{
		"projectId":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		"goos":        &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		"goarch":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		"generatePsk": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Boolean)},
	},
	Resolve: wrapper(func(rctx resolveContext[any]) (string, error) {
		goos := rctx.getStringArg("goos")
		projectId := rctx.getStringArg("projectId")
		err := rctx.canAccessProject(projectId)
		if err != nil {
			return "", err
		}

		// Optionally generate a pre-shared key for the agent to use
		var psk string
		if rctx.getBoolArg("generatePsk") {
			psk, err = rctx.db.GeneratePreSharedKey(rctx, projectId)
			if err != nil {
				return "", fmt.Errorf("generating psk: %w", err)
			}
		}

		switch goos {
		case "linux":
			// We can use the same command for all linux architectures
			// and the installation script will handle the nuances and
			// platform detection processes that need to happen.
			return buildLinuxInstallCommand(psk), nil
		default:
			return "", fmt.Errorf("unsupported goos: %s", goos)
		}
	}),
}

func buildLinuxInstallCommand(psk string) string {
	if len(psk) == 0 {
		return fmt.Sprintf("curl -fsSL https://raw.githubusercontent.com/clarkmcc/cloudcore/main/scripts/linux/install.sh | sh")
	} else {
		return fmt.Sprintf("curl -fsSL https://raw.githubusercontent.com/clarkmcc/cloudcore/main/scripts/linux/install.sh | sh -s -- --psk %s", psk)
	}
}
