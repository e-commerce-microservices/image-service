.PHONY: protogen

protogen:
	protoc --proto_path=proto proto/image_service.proto proto/general.proto \
	--go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative

.PHONY: build
build:
	docker build -t ngoctd/ecommerce-image:latest .

.PHONY: rebuild
rebuild:
	docker build -t ngoctd/ecommerce-image:latest . && \
	docker push ngoctd/ecommerce-image

.PHONY: redeploy
redeploy:
	kubectl rollout restart deployment depl-image