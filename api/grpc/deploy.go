package grpc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/keti-openfx/openfx-cli/config"
	"github.com/keti-openfx/openfx-cli/pb"
	grpcgo "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeployConfig struct {
	FxGateway string
	FxProcess string

	FunctionName string
	Image        string
	EnvVars      map[string]string
	Labels       map[string]string
	Annotations  map[string]string

	Network     string
	Constraints *[]string
	Secrets     []string

	Limits   *config.FunctionResources
	Requests *config.FunctionResources

	Update  bool
	Replace bool
}

func Deploy(c DeployConfig) error {

	gateway := strings.TrimRight(c.FxGateway, "/")

	conn, err := grpcgo.Dial(gateway, grpcgo.WithInsecure())
	if err != nil {
		return errors.New("did not connect: " + err.Error())
	}
	client := pb.NewFxGatewayClient(conn)

	if c.Replace {
		_, statusErr := client.Delete(context.Background(), &pb.DeleteFunctionRequest{FunctionName: c.FunctionName})
		if statusErr != nil {
			return errors.New("did not delete: " + statusErr.Error())
		}
	}

	req := pb.CreateFunctionRequest{
		Image:       c.Image,
		Service:     c.FunctionName,
		EnvVars:     c.EnvVars,
		Secrets:     c.Secrets, // TODO: allow registry auth to be specified or read from local Docker credentials store
		Labels:      c.Labels,
		Annotations: c.Annotations,
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

	if c.Update {
		_, statusErr := client.Update(context.Background(), &req)
		st, ok := status.FromError(statusErr)
		if !ok {
			return errors.New("Invaild status error.")
		}
		if st.Code() == codes.NotFound {
			fmt.Println("Attempting update... but Function Not Found. Deploying Function...")
		} else if st.Code() == codes.OK {
			fmt.Printf("Function %s already exists, attempting rolling-update.\n", c.FunctionName)
			return nil
		}
	}

	_, statusErr := client.Deploy(context.Background(), &req)
	st, ok := status.FromError(statusErr)
	if !ok {
		return errors.New("Invaild status error.")
	}
	if st.Code() == codes.AlreadyExists {
		fmt.Printf("Function %s already exists. failed deploying.\n", c.FunctionName)
	}
	if statusErr != nil {
		return statusErr
	}

	return nil
}
