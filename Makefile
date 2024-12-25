SRC_DIR=proto
DST_DIR=$(SRC_DIR)/gen

build-proto:
	@echo "build proto"

	protoc -I=${SRC_DIR} \
		--go_out=${DST_DIR} --go_opt=paths=source_relative \
		--go-grpc_out=${DST_DIR} --go-grpc_opt=paths=source_relative \
		${SRC_DIR}/**/*.proto

docker-build:
	docker build --platform linux/amd64 -t grpc-hello-world .

docker-push:
	docker tag grpc-hello-world:latest newinh/grpc-hello-world:latest
	docker push newinh/grpc-hello-world:latest
