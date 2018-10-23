all: fxgateway

fxgateway:
	protoc -I . \
		-I${GOPATH}/src/github.com/keti-openfx/openfx/pb \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:. \
		${GOPATH}/src/github.com/keti-openfx/openfx/pb/fxgateway.proto
