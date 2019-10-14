package function

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/keti-openfx/openfx-cli/cmd/runtime"
	"github.com/keti-openfx/openfx-cli/config"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var (
	runtimeName string
	handlerDir  string
	handlerName string
)

func init() {
	initCmd.Flags().StringVarP(&configFile, "config", "f", "", "Path to YAML config file describing function(s)")
	initCmd.Flags().StringVarP(&runtimeName, "runtime", "r", "", "Runtime(Language) to use")
	initCmd.Flags().StringVarP(&handlerDir, "dir", "d", "", "Directory containing handler file")
	initCmd.MarkFlagRequired("runtime")
}

var initCmd = &cobra.Command{
	Use: `init <FUNCTION_NAME> --runtime <RUNTIME>
  openfx-cli function init <FUNCTION_NAME> -r <RUNTIME> [-f <APPEND_EXISTING_YAML_FILE>] [-g <FX_GATEWAY_ADDRESS>]`,
	Short: "Prepare a OpenFx function",
	Long: `
	The init command creates a new function template based upon hello-world in the given runtime. When user execute init command, config file, runtime directory, and directory with function name are created. Also, in directory with function name, there is handler file and user can modify this file later. 
`,
	Example: `  openfx-cli function init echo --runtime go
  openfx-cli function init read-write -f ./config.yml -r python3
  openfx-cli function init read-write --config ./config.yml --runtime java --gateway localhost:31113
  `,
	PreRunE: preRunInit,
	RunE:    runInit,
}

// validateFunctionName provides least-common-denominator validation - i.e. only allows valid Kubernetes services names
func validateFunctionName(functionName string) error {
	// Regex for RFC-1123 validation:
	// 	k8s.io/kubernetes/pkg/util/validation/validation.go
	var validDNS = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
	matched := validDNS.MatchString(functionName)
	if matched {
		return nil
	}
	return fmt.Errorf(`function name can only contain a-z, 0-9 and dashes`)
}

func preRunInit(cmd *cobra.Command, args []string) error {

	if len(args) < 1 || len(args) > 1 {
		return fmt.Errorf("please provide a name for the function")
	}

	functionName = args[0]
	if err := validateFunctionName(functionName); err != nil {
		return err
	}

	if configFile == "" {
		configFile = config.DefaultConfigFile
	}

	if _, err := os.Stat(configFile); err != nil {
		if os.IsNotExist(err) {
			fxServices = config.NewServices()
		} else {
			return err
		}
	} else {
		if err := parseConfigFile(); err != nil {
			return err
		}
	}

	if len(gateway) > 0 {
		fxServices.Openfx.FxGatewayURL = gateway
	}

	if _, ok := fxServices.Functions[functionName]; ok {
		return fmt.Errorf("Function %s already exists in %s file.", functionName, configFile)
	}

	if !runtime.ExistFileOrDir(config.DefaultRuntimeDir) {
		if err := runtime.DownloadRuntimes(config.DefaultRuntimeDir, config.DefaultRuntimeRepo); err != nil {
			return err
		}
	}

	return nil
}

func runInit(cmd *cobra.Command, args []string) error {

	if _, err := os.Stat(functionName); err == nil {
		return fmt.Errorf("folder: %s already exists", functionName)
	}

	r, err := runtime.GetRuntime(functionName, runtimeName, config.DefaultRuntimeDir)
	if err == nil {
		fmt.Printf("Folder: %s created.\n", functionName)
	} else {
		return fmt.Errorf("folder: could not create %s\n %s", functionName, err)
	}

	if handlerDir == "" {
		handlerDir = "./src"
	}

	fxServices.Functions[functionName] = config.Function{
		Runtime:      runtimeName,
		Handler:      config.Handler{Dir: handlerDir, Name: r.Handler.Name, File: r.Handler.File},
		RegistryURL:  config.DefaultRegistry,
		Image:        config.DefaultRegistry + "/" + functionName,
		BuildArgs:    r.BuildArgs,
		BuildOptions: r.BuildPackages,
	}

	confYaml, err := yaml.Marshal(&fxServices)
	if err != nil {
		return err
	}

	fmt.Printf("Function handler created in folder: %s\n", functionName+"/src")
	fmt.Printf("Rewrite the function handler code in %s folder\n", functionName+"/src")

	confWriteErr := ioutil.WriteFile("./"+functionName+"/"+configFile, confYaml, 0600)
	if confWriteErr != nil {
		return fmt.Errorf("error writing config file %s", confWriteErr)
	}

	fmt.Printf("Config file written: %s\n", configFile)

	return nil
}
