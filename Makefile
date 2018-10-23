all: proto build

proto:
	cd pb; $(MAKE)
build:
	go build
	go install
