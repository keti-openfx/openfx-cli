package grpc

import (
	"context"
	"errors"
	"strings"

	"github.com/keti-openfx/openfx-cli/pb"
	grpcgo "google.golang.org/grpc"
)

func Invoke(functionName, fxGateway string, input []byte) (string, error) {

	gateway := strings.TrimRight(fxGateway, "/")

	conn, err := grpcgo.Dial(gateway, grpcgo.WithInsecure())
	if err != nil {
		return "", errors.New("did not connect: " + err.Error())
	}
	client := pb.NewFxGatewayClient(conn)

	message, statusErr := client.Invoke(context.Background(), &pb.InvokeServiceRequest{Service: functionName, Input: input})
	if statusErr != nil {
		return "", errors.New("did not invoke: " + statusErr.Error())
	}

	return message.Msg, nil
}
