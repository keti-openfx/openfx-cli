package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

func NewServices() *Services {
	var services Services
	services.Openfx.FxGatewayURL = DefaultGatewayURL

	services.Functions = make(map[string]Function, 0)

	return &services
}

func ParseConfigFile(file string) (*Services, error) {
	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var services Services
	err = yaml.Unmarshal(fileData, &services)
	if err != nil {
		fmt.Printf("Error with YAML Config file\n")
		return nil, err
	}

	return &services, nil
}

func GetFxGatewayURL(argumnetURL, configURL string) string {
	var url string

	envURL := os.Getenv(GatewayEnvVarKey)

	if len(argumnetURL) > 0 && argumnetURL != DefaultGatewayURL {
		url = argumnetURL
	} else if len(configURL) > 0 && configURL != DefaultGatewayURL {
		url = configURL
	} else if len(envURL) > 0 {
		url = envURL
	} else {
		url = DefaultGatewayURL
	}

	url = strings.TrimRight(url, "/")

	return url
}
