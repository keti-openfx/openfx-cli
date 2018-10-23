package grpc

import (
	"context"
	"errors"
	"strings"

	"github.com/keti-openfx/openfx-cli/pb"
	grpcgo "google.golang.org/grpc"
)

func GetLog(functionName, dcfGateway string) (string, error) {

	gateway := strings.TrimRight(dcfGateway, "/")

	conn, err := grpcgo.Dial(gateway, grpcgo.WithInsecure())
	if err != nil {
		return "", errors.New("did not connect: " + err.Error())
	}
	client := pb.NewFxGatewayClient(conn)

	reply, statusErr := client.GetLog(context.Background(), &pb.FunctionRequest{FunctionName: functionName})
	if statusErr != nil {
		return "", errors.New("did not get log: " + statusErr.Error())
	}

	return reply.Msg, nil
}
