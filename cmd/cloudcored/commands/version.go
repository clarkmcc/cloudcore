package commands

import (
	"github.com/clarkmcc/cloudcore/pkg/utils"
	"github.com/clarkmcc/cloudcore/pkg/version"
	"github.com/spf13/cobra"
)

var Version = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintStruct(version.Get())
	},
}
