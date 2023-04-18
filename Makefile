# set the binary names
SERVER_BINARY_NAME=server
CLIENT_BINARY_NAME=client

# set the build flags
BUILD_FLAGS=-ldflags="-s -w"

# set the docker image names and tags
SERVER_IMAGE_NAME=server
SERVER_IMAGE_TAG=latest
CLIENT_IMAGE_NAME=client
CLIENT_IMAGE_TAG=latest

#build the binaries
build: build-server build-client
# build the binaries and docker
all: build-server build-client scratch docker-server docker-client

build-server:
	@echo "Building server..."
	@go build ${BUILD_FLAGS} -o ${SERVER_BINARY_NAME} server.go
	@upx ${SERVER_BINARY_NAME}

build-client:
	@echo "Building client..."
	@go build ${BUILD_FLAGS} -o ${CLIENT_BINARY_NAME} client.go
	@upx ${CLIENT_BINARY_NAME}

# import scratch
scratch:
	@echo "Building scratch docker image..."
	tar cv --files-from /dev/null | docker import - scratch

# docker
docker-server:
	@echo "Building server docker image..."
	@docker build -t ${SERVER_IMAGE_NAME}:${SERVER_IMAGE_TAG} -f ./Dockerfile/Dockerfile.server .

docker-client:
	@echo "Building client docker image..."
	@docker build -t ${CLIENT_IMAGE_NAME}:${CLIENT_IMAGE_TAG} -f ./Dockerfile/Dockerfile.client .

# cleanup
clean:
	@echo "Cleaning up..."
	@rm -f ${SERVER_BINARY_NAME} ${CLIENT_BINARY_NAME}
clean-docker:
	# remove any existing images
	@echo "Removing existing images..."
	@docker rmi ${SERVER_IMAGE_NAME}:${SERVER_IMAGE_TAG} ${CLIENT_IMAGE_NAME}:${CLIENT_IMAGE_TAG}
clean-all: clean clean-docker