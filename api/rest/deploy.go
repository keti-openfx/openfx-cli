package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/keti-openfx/openfx-cli/cmd/log"
	"net/http"
	"strings"
	"time"

	"github.com/keti-openfx/openfx-cli/api/types"
	"github.com/keti-openfx/openfx-cli/pb"
)

func DeployFunction(c types.DeployConfig) {

	statusCode, deployOutput := deployFunction(c)

	if c.Update == true {
		if statusCode == http.StatusNotFound {
			// Re-run the function with update=false
			c.Update = false
			log.Println("Function Not Found. Deploying Function...")
			_, deployOutput = deployFunction(c)
		} else if statusCode == http.StatusOK {

			log.Println("Function %s already exists, attempting rolling-update.", c.FunctionName)
		}
	}

	fmt.Println(deployOutput)
}

func deployFunction(c types.DeployConfig) (int, string) {

	var deployOutput string
	gateway := strings.TrimRight(c.FxGateway, "/")

	if c.Replace {
		//TODO
		DeleteFunction(gateway, c.FunctionName)
	}

	req := pb.CreateFunctionRequest{
		EnvProcess: c.FxProcess,
		Image:      c.Image,
		Network:    c.Network,
		Service:    c.FunctionName,
		EnvVars:    c.EnvVars,
		Secrets:    c.Secrets, // TODO: allow registry auth to be specified or read from local Docker credentials store
		Labels:     c.Labels,
	}

	if c.Constraints != nil {
		req.Constraints = *c.Constraints
	}

	hasLimits := false
	req.Limits = &pb.FunctionResources{}
	if c.Limits != nil && len(c.Limits.Memory) > 0 {
		hasLimits = true
		req.Limits.Memory = c.Limits.Memory
	}
	if c.Limits != nil && len(c.Limits.CPU) > 0 {
		hasLimits = true
		req.Limits.CPU = c.Limits.CPU
	}
	if c.Limits != nil && len(c.Limits.GPU) > 0 {
		hasLimits = true
		req.Limits.GPU = c.Limits.GPU
	}
	if !hasLimits {
		req.Limits = nil
	}

	hasRequests := false
	req.Requests = &pb.FunctionResources{}
	if c.Requests != nil && len(c.Requests.Memory) > 0 {
		hasRequests = true
		req.Requests.Memory = c.Requests.Memory
	}
	if c.Requests != nil && len(c.Requests.CPU) > 0 {
		hasRequests = true
		req.Requests.CPU = c.Requests.CPU
	}

	if !hasRequests {
		req.Requests = nil
	}

	reqBytes, _ := json.Marshal(&req)
	reader := bytes.NewReader(reqBytes)
	var request *http.Request

	timeout := 60 * time.Second
	client := MakeHTTPClient(&timeout)

	method := http.MethodPost
	// "application/json"
	if c.Update {
		method = http.MethodPut
	}

	var err error
	request, err = http.NewRequest(method, gateway+"/system/functions", reader)

	//TODO
	// set OAUTH, Basic auth ....
	//SetAuth(request, gateway)

	if err != nil {
		deployOutput = fmt.Sprintln(err)
		return http.StatusInternalServerError, deployOutput
	}

	res, err := client.Do(request)
	if err != nil {
		deployOutput = fmt.Sprintf("Make sure you created the function and set the correct gateway url.\n %s\n", err.Error())
		return http.StatusInternalServerError, deployOutput
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	switch res.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		deployOutput += fmt.Sprintf("Deployed. %s.\n", res.Status)

		deployedURL := fmt.Sprintf("URL: %s/function/%s", gateway, c.FunctionName)
		deployOutput += fmt.Sprintln(deployedURL)
	case http.StatusUnauthorized:
		deployOutput += fmt.Sprintln("unauthorized access, run \"fxcli login\" to setup authentication for this server")
		/*
			case http.StatusNotFound:
				if replace && !update {
				deployOutput += fmt.Sprintln("Could not delete-and-replace function because it is not found (404)")
		*/
	default:
		bytesOut, err := ioutil.ReadAll(res.Body)
		if err == nil {
			deployOutput += fmt.Sprintf("Unexpected status: %d, message: %s\n", res.StatusCode, string(bytesOut))
		}
	}

	return res.StatusCode, deployOutput
}
