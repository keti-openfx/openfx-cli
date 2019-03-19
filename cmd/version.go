package cmd

import (
	"fmt"

	"github.com/keti-openfx/openfx-cli/config"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "display the Openfx CLI version information",
	Long: `
	display the Openfx CLI version information
	`,
	Example: `openfx-cli version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", config.FxCliVersion)
	},
}
