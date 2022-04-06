GOPATH:=$(shell go env GOPATH)

.PHONY: init
init:
	@go get -u google.golang.org/protobuf/proto
	@go install github.com/golang/protobuf/protoc-gen-go@latest
	@go install github.com/asim/go-micro/cmd/protoc-gen-micro/v4@latest

.PHONY: proto
proto:
	@protoc --proto_path=./proto/pb --micro_out=./protobuf/pb --gofast_out=:./protobuf/pb proto/pb/*.proto

.PHONY: update
update:
	@go get -u

.PHONY: tidy
tidy:
	@go env -w GOPROXY=https://goproxy.cn,direct
	@go env -w GO111MODULE=on
	@go mod tidy

.PHONY: build
build:
	@go build -o shopproduct *.go

.PHONY: test
test:
	@go test -v ./... -cover

.PHONY: docker
docker:
	@docker build -t shopproduct:latest .
