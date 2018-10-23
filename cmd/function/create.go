package function

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/keti-openfx/openfx-cli/api/grpc"
	"github.com/keti-openfx/openfx-cli/builder"
	"github.com/keti-openfx/openfx-cli/cmd/log"
	"github.com/keti-openfx/openfx-cli/cmd/runtime"
	"github.com/keti-openfx/openfx-cli/config"
	"github.com/spf13/cobra"
)

var (
	replace bool
	update  bool
	noCache bool
	verbose bool
)

func init() {
	createCmd.Flags().StringVarP(&configFile, "config", "f", "", "Path to YAML config file describing function(s)")

	createCmd.Flags().BoolVar(&noCache, "nocache", false, "Do not use cache when building runtime image")
	createCmd.Flags().BoolVar(&replace, "replace", false, "Remove and re-create existing function(s)")
	createCmd.Flags().BoolVar(&update, "update", false, "Perform rolling update on existing function(s)")
	createCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print function build log")
	createCmd.MarkFlagRequired("config")
}

var createCmd = &cobra.Command{
	Use:   `create -f <YAML_CONIFIG_FILE>`,
	Short: "Create OpenFx functions",
	Long: `
	Build OpenFx function Image & Deploys OpenFx function containers via the supplied YAML config using the "-f" flag (which may contain multiple function definitions)
`,
	Example: `  openfx function create -f config.yml
  openfx function create -f ./config.yml --replace=false --update=true
  openfx function create -f ./functions.yaml --replace=true --update=false
                  `,
	PreRunE: preRunCreate,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runCreate(); err != nil {
			fmt.Println(err.Error())
		}
		return
	},
}

func preRunCreate(cmd *cobra.Command, args []string) error {

	if update && replace {
		return errors.New(`one of "--update" flag or "--replace" flag must be false\n`)
	}

	var configURL string
	if configFile == "" {
		e := fmt.Sprintf("please provide a '-f' flag for function creation\n")
		return errors.New(e)
	} else {
		if err := parseConfigFile(); err != nil {
			return err
		}
		configURL = fxServices.Openfx.FxGatewayURL
	}
	gateway = config.GetFxGatewayURL(gateway, configURL)

	if !runtime.ExistFileOrDir(config.DefaultRuntimeDir) {
		if err := runtime.DownloadRuntimes(config.DefaultRuntimeDir, config.DefaultRuntimeRepo); err != nil {
			return err
		}
	}

	return nil
}

func build(nocache, verbose bool, function config.Function) error {
	var err error
	buildArgs, err := parseBuildArgs(function.BuildArgs)
	if err != nil {
		return err
	}

	if function.Handler.File != "" {
		buildArgs["handler_file"] = function.Handler.File
	}
	if function.Handler.Name != "" {
		buildArgs["handler_name"] = function.Handler.Name
	}

	result, err := builder.BuildImage(function.Image, function.Handler.Dir, function.Name, function.Runtime, nocache, buildArgs, function.BuildOptions, verbose)
	if err != nil {
		log.Print(result)
		return err
	}
	return nil
}

func deploy(gw string, function config.Function, update, replace bool) error {
	//function.Secrets
	//sendRegistryAuth
	//EnvVar
	fileEnvironment, err := readFiles(function.EnvironmentFile)
	if err != nil {
		return err
	}
	allEnvironment := mergeMap(function.Environment, fileEnvironment)

	//Labels
	labelMap := map[string]string{}
	if function.Labels != nil {
		labelMap = *function.Labels
	}

	//Annotations
	AnnoMap := map[string]string{}
	if function.Maintainer != "" {
		AnnoMap["maintainer"] = function.Maintainer
	}
	if function.Description != "" {
		AnnoMap["desc"] = function.Description
	}

	// Get FxProcess to use from the ?
	deployConfig := grpc.DeployConfig{
		FxGateway:    gw,
		FunctionName: function.Name,
		Image:        function.Image,
		EnvVars:      allEnvironment,
		Labels:       labelMap,
		Annotations:  AnnoMap,
		Constraints:  function.Constraints,
		Secrets:      append(function.Secrets, "regcred"),
		Limits:       function.Limits,
		Requests:     function.Requests,

		Update:  update,
		Replace: replace,
	}
	if err := grpc.Deploy(deployConfig); err != nil {
		return err
	}
	return nil
}

func runCreate() error {
	if len(fxServices.Functions) <= 0 {
		return errors.New("")
	}

	for name, function := range fxServices.Functions {
		function.Name = name

		//BUILD
		if function.SkipBuild {
			log.Print("Skipping build: %s\n", function.Name)
		} else {
			log.Info("Building: %s, Image:%s\n", function.Name, function.Image)
			if err := build(noCache, verbose, function); err != nil {
				return err
			}
		}

		//PUSH
		if function.SkipBuild {
			log.Print("Skipping push: %s\n", function.Name)
		} else {
			log.Info("Pushing: %s, Image:%s\n", function.Name, function.Image)
			if verbose {
				err := builder.ExecCommandPipe("./", []string{"docker", "push", function.Image}, os.Stdout, os.Stderr)
				if err != nil {
					return err
				}
			} else {
				_, err := builder.ExecCommand("./", []string{"docker", "push", function.Image})
				if err != nil {
					return err
				}
			}
		}

		log.Info("Deploying: %s\n", function.Name)
		//DEPLOY
		if err := deploy(gateway, function, update, replace); err != nil {
			return err
		}

		log.Info("http trigger url: http://%s/function/%s \n", gateway, function.Name)

	}

	return nil
}

func readFiles(files []string) (map[string]string, error) {
	envs := make(map[string]string)

	for _, file := range files {
		bytesOut, readErr := ioutil.ReadFile(file)
		if readErr != nil {
			return nil, readErr
		}

		envFile := config.EnvironmentFile{}
		unmarshalErr := yaml.Unmarshal(bytesOut, &envFile)
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}
		for k, v := range envFile.Environment {
			envs[k] = v
		}
	}
	return envs, nil
}

func mergeMap(i map[string]string, j map[string]string) map[string]string {
	merged := make(map[string]string)

	for k, v := range i {
		merged[k] = v
	}
	for k, v := range j {
		merged[k] = v
	}
	return merged
}

func parseBuildArgs(args []string) (map[string]string, error) {
	mapped := make(map[string]string)

	for _, kvp := range args {
		index := strings.Index(kvp, "=")
		if index == -1 {
			return nil, fmt.Errorf("each build-arg must take the form key=value")
		}

		values := []string{kvp[0:index], kvp[index+1:]}

		k := strings.TrimSpace(values[0])
		v := strings.TrimSpace(values[1])

		if len(k) == 0 {
			return nil, fmt.Errorf("build-arg must have a non-empty key")
		}
		if len(v) == 0 {
			return nil, fmt.Errorf("build-arg must have a non-empty value")
		}

		mapped[k] = v
	}

	return mapped, nil

}
