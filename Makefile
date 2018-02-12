GHQ_ROOT = $(shell ghq root)
PROTOBUF_DIR = ${GHQ_ROOT}/github.com/google/protobuf

stat.pb.go: stat.proto
	protoc --go_out=plugins=grpc:./ $^