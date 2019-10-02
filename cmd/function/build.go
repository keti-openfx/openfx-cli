package function

import (
	"errors"
	"fmt"
	"strings"
	"io/ioutil"

	"github.com/keti-openfx/openfx-cli/builder"
	"github.com/keti-openfx/openfx-cli/cmd/log"
	"github.com/keti-openfx/openfx-cli/config"
	"github.com/spf13/cobra"
)

var (
	noCache      bool
	buildVerbose bool
)

func init() {
	buildCmd.Flags().StringVarP(&configFile, "config", "f", "", "Path to YAML config file describing function(s)")
	buildCmd.Flags().BoolVar(&noCache, "nocache", false, "Do not use cache when building runtime image")
	buildCmd.Flags().BoolVarP(&buildVerbose, "buildverbose", "v", false, "Print function build log")
}

var buildCmd = &cobra.Command{
	Use:   `build -f <YAML_CONFIG_FILE>`,
	Short: "Build OpenFx function Image",
	Long: `
	Build OpenFx function Image via the supplied YAML config 
	`,
	Example: `
	openfx-cli function build 
	openfx-cli function build -f ./config.yaml
	openfx-cli function build -v
	`,
	PreRunE: preRunBuild,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runBuild(); err != nil {
			fmt.Println(err.Error())
		}
		return
	},
}

func preRunBuild(cmd *cobra.Command, args []string) error {
	if configFile == "" {
		files, err := ioutil.ReadDir("./")
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			if strings.Contains(f.Name(), "yaml") {
				configFile = f.Name()
			}
		}

		if err := parseConfigFile(); err != nil { 
			return err
		}

	} else {
		if err := parseConfigFile(); err != nil {
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

	result, err := builder.BuildImage(function.Image, function.Handler.Dir, function.Name, function.RegistryURL, function.Runtime, nocache, buildArgs, function.BuildOptions, verbose)
	if err != nil {
		log.Print(result)
		return err
	}

	return nil
}

func runBuild() error {
	if len(fxServices.Functions) <= 0 {
		return errors.New("")
	}

	for name, function := range fxServices.Functions {
		function.Name = name

		log.Info("Building function (%s) image...\n", function.Name)
		if err := build(noCache, buildVerbose, function); err != nil {
			return err
		}
	}

	return nil
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
