package function

import (
	"errors"
	"fmt"

	"github.com/keti-openfx/openfx-cli/api/grpc"
	"github.com/keti-openfx/openfx-cli/cmd/log"
	"github.com/keti-openfx/openfx-cli/config"
	"github.com/spf13/cobra"
)

func init() {
	deleteCmd.Flags().StringVarP(&configFile, "config", "f", "", "Path to YAML config file describing function(s)")
	deleteCmd.Flags().StringVarP(&gateway, "gateway", "g", "localhost:31113", "Gateway URL to store in YAML config file")
}

var deleteCmd = &cobra.Command{
	Use:     `delete -f <YAML_CONIFIG_FILE>`,
	Aliases: []string{"remove", "rm"},
	Short:   "Delete OpenFx functions",
	Long: `
	Delete OpenFx function via the supplied YAML config using
the "-f" flag (which may contain multiple function definitions)
`,
	Example: `  openfx-cli function delete -f config.yml
                  `,
	PreRunE: preRunDelete,
	Run: func(cmd *cobra.Command, args []string) {

		if err := runDelete(); err != nil {
			fmt.Println(err.Error())
		}
		return
	},
}

func preRunDelete(cmd *cobra.Command, args []string) error {
	fxServices = config.NewServices()

	var configURL string
	if cmd.Flag("config").Value.String() != "" {
		if err := parseConfigFile(); err != nil {
			return err
		}
		configURL = fxServices.Openfx.FxGatewayURL
	}

	gateway = config.GetFxGatewayURL(gateway, configURL)

	if len(args) > 0 {
		fxServices.Functions = make(map[string]config.Function, 0)
		fxServices.Functions[args[0]] = config.Function{}
	}

	return nil
}

func runDelete() error {
	if len(fxServices.Functions) <= 0 {
		return errors.New("")
	}

	for name, function := range fxServices.Functions {
		function.Name = name
		if err := grpc.Delete(gateway, function.Name); err != nil {
			return err
		}

		log.Print("Deleted: %s\n", function.Name)
	}

	return nil
}
