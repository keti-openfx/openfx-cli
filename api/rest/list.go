package rest

import (
	"encoding/json"

	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/keti-openfx/openfx-cli/pb"
)

// ListFunctions list deployed functions
func ListFunctions(gateway string) ([]pb.Function, error) {
	var results []pb.Function

	gateway = strings.TrimRight(gateway, "/")

	timeout := 60 * time.Second
	client := MakeHTTPClient(&timeout)

	getRequest, err := http.NewRequest(http.MethodGet, gateway+"/system/functions", nil)
	//TODO
	//SetAuth(getRequest, gateway)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to OpenFx on URL: %s", gateway)
	}

	res, err := client.Do(getRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to OpenFx on URL: %s", gateway)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	switch res.StatusCode {
	case http.StatusOK:

		bytesOut, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("cannot read result from OpenFx on URL: %s", gateway)
		}
		jsonErr := json.Unmarshal(bytesOut, &results)
		if jsonErr != nil {
			return nil, fmt.Errorf("cannot parse result from OpenFx on URL: %s\n%s", gateway, jsonErr.Error())
		}
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("unauthorized access, run \"fxcli login\" to setup authentication for this server")
	default:
		bytesOut, err := ioutil.ReadAll(res.Body)
		if err == nil {
			return nil, fmt.Errorf("server returned unexpected status code: %d - %s", res.StatusCode, string(bytesOut))
		}
	}
	return results, nil
}
