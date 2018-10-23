package runtime

import (
	"github.com/spf13/cobra"
)

var (
	overwrite bool
)

func init() {
	pullCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite existing runtime")
}

var pullCmd = &cobra.Command{
	Use:   `pull <REPOSITORY_URL>`,
	Short: "Downloads runtime from the specified github repo",
	Long: `
	Downloads the compressed github repo specified by [URL], and extracts the 'runtime'
	directory from the root of the repo, if it exists.
	`,
	Example: `openfx runtime pull https://github.com/openfx/openfx/runtimes
	openfx runtime pull http://`,
	RunE: runPull,
}

func runPull(cmd *cobra.Command, args []string) error {
	return nil
}
