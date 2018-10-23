package function

import (
	"errors"
	"fmt"

	"github.com/keti-openfx/openfx-cli/api/grpc"
	"github.com/keti-openfx/openfx-cli/config"
	"github.com/spf13/cobra"
)

func init() {
	listCmd.Flags().StringVarP(&configFile, "config", "f", "", "Path to YAML config file describing function(s)")
	listCmd.Flags().StringVarP(&gateway, "gateway", "g", "localhost:31113", "Gateway URL to store in YAML config file")
}

var listCmd = &cobra.Command{
	Use:   `list -f <YAML_CONIFIG_FILE>`,
	Short: "Lists OpenFx functions",
	Long: `
	Lists OpenFx function
`,
	Example: `  openfx function list -f config.yml
	openfx funtion list -g localhost:31113
                  `,
	PreRunE: preRunList,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runList(); err != nil {
			fmt.Println(err.Error())
		}
		return
	},
}

func preRunList(cmd *cobra.Command, args []string) error {
	var configURL string
	if cmd.Flag("config").Value.String() != "" {
		if err := parseConfigFile(); err != nil {
			return err
		}
		configURL = fxServices.Openfx.FxGatewayURL
	}
	gateway = config.GetFxGatewayURL(gateway, configURL)

	return nil
}

func runList() error {

	if gateway == "" {
		return errors.New("please provide a gateway url")
	}

	fnList, err := grpc.List(gateway)
	if err != nil {
		return err
	}

	fmt.Printf("%-15s\t%-20s\t%-10s\t%-10s\t%-10s\n", "Function", "Image", "Invocations", "Replicas", "Status")
	for _, fn := range fnList.Functions {
		fnImage := fn.Image
		if len(fnImage) > 30 {
			fnImage = fnImage[0:28] + ".."
		}
		fmt.Printf("%-15s\t%-20s\t%-10d\t%-10d\t", fn.Name, fnImage, fn.InvocationCount, fn.Replicas)

		if fn.AvailableReplicas == 0 {
			fmt.Printf("%-10s\n", "Not Ready")
		} else {
			fmt.Printf("%-10s\n", "Ready")
		}
	}

	return nil
}
