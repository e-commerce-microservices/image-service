.PHONY: protogen

protogen:
	protoc --proto_path=proto proto/image_service.proto proto/general.proto \
	--go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative