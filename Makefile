

proto:
	protoc --go_out=services --go-grpc_out=services services/registry/registry.proto

.PHONY: proto	