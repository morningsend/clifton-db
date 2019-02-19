gen-all: gen-client-proto gen-raft-proto
	echo "generating"

gen-client-proto:
	protoc -I=./kv-client-api -I=${GOPATH}/src -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf --gofast_out=plugins=grpc:./kv-client-api/ kv-client-api/api.proto

gen-raft-proto:
	protoc -I=./pkg/raft/proto -I=${GOPATH}/src -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf --gofast_out=plugins=grpc:./pkg/raft/proto pkg/raft/proto/raft.proto
