package runtime

import (
	"fmt"

	"github.com/keti-openfx/openfx-cli/config"
	"github.com/spf13/cobra"
)

func init() {
}

var listCmd = &cobra.Command{
	Use:   `list`,
	Short: "list a OpenFx runtime",
	Long:  `lists of Supported Runtimes which can operate with your function`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !ExistFileOrDir(config.DefaultRuntimeDir) {
			if err := DownloadRuntimes(config.DefaultRuntimeDir, config.DefaultRuntimeRepo); err != nil {
				return err
			}
		}
		return nil
	},
	RunE: runList,
}

func runList(cmd *cobra.Command, args []string) error {

	dir := config.DefaultRuntimeDir

	if runtimeDir != "" {
		dir = runtimeDir
	}

	runtimes, err := ReadRuntimeList(dir)
	if err != nil {
		return err
	}

	fmt.Printf("Supported Runtimes are:\n")

	for name, runtime := range runtimes.Runtimes {
		runtime.Name = name
		fmt.Printf("- %s\n", name)
	}
	fmt.Println()

	return nil
}
