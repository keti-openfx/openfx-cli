package function

import (
	"os"

	"github.com/keti-openfx/openfx-cli/config"
	"github.com/spf13/cobra"
)

var (
	configFile   string
	fxServices   *config.Services
	gateway      string
	functionName string
)

var FunctionCmd = &cobra.Command{
	Use:     "function SUBCOMMAND",
	Aliases: []string{"fn"},
	Short:   "function specific operations",
	Long: `
	function command allows user to init, list, deploy, edit, delete functions running on Openfx
	`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	FunctionCmd.AddCommand(initCmd)
	FunctionCmd.AddCommand(buildCmd)
	//FunctionCmd.AddCommand(runCmd)
	FunctionCmd.AddCommand(deployCmd)
	FunctionCmd.AddCommand(deleteCmd)
	FunctionCmd.AddCommand(listCmd)
	FunctionCmd.AddCommand(callCmd)
	FunctionCmd.AddCommand(infoCmd)
	FunctionCmd.AddCommand(logCmd)

	FunctionCmd.PersistentFlags().StringVarP(&gateway, "gateway", "g", "", "Set gateway URL")
}

func parseConfigFile() error {

	var err error
	if _, err := os.Stat(configFile); err != nil {
		if os.IsNotExist(err) {
		}
		return err
	}
	if fxServices, err = config.ParseConfigFile(configFile); err != nil {
		return err
	}

	return nil
}
