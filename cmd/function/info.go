package function

import (
	"fmt"

	"github.com/keti-openfx/openfx-cli/api/grpc"
	"github.com/keti-openfx/openfx-cli/config"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

type FunctionInfo struct {
	Name              string            `yaml:"name"`
	Image             string            `yaml:"image"`
	InvocationCount   uint64            `yaml:"invocationcount"`
	Replicas          uint64            `yaml:"replicas"`
	Annotations       map[string]string `yaml:"annotations"`
	AvailableReplicas uint64            `yaml:"availablereplicas"`
	Labels            map[string]string `yaml:"labels"`
}

func init() {
	infoCmd.Flags().StringVarP(&configFile, "config", "f", "", "Path to YAML config file describing function(s)")
}

var infoCmd = &cobra.Command{
	Use:   `info <FUNCTION_NAME>`,
	Short: "Display OpenFx function information",
	Long: `
	Display OpenFX function information
`,
	Example: `  openfx-cli function info -f config.yml
	openfx-cli function info -g localhost:31113
                  `,
	PreRunE: preRunInfo,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("please provide a name for the function")
		}

		functionName = args[0]

		if err := runInfo(); err != nil {
			return err
		}
		return nil
	},
}

func preRunInfo(cmd *cobra.Command, args []string) error {
	gateway = config.GetFxGatewayURL(gateway, "")

	return nil
}

func runInfo() error {
	fnInfo, err := grpc.GetMeta(functionName, gateway)
	if err != nil {
		return err
	}

	yamlInfo, err := yaml.Marshal(&fnInfo)
	if err != nil {
		return err
	}

	var fi FunctionInfo
	err = yaml.Unmarshal(yamlInfo, &fi)
	if err != nil {
		return err
	}

	yamlInfo, err = yaml.Marshal(&fi)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", string(yamlInfo))

	return nil
}
