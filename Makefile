PACKAGES = $(shell go list ./src/...)

.PHONY: build doc fmt lint run test vet test-cover-html test-cover-func collect-cover-data AUTHORS

# Used to populate version variable in main package.
VERSION=v1.0.0
BUILD_TIME=$(shell date -u +%Y-%m-%d:%H-%M-%S)
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD_TIME}"

# set image name and tags
REGISTRY=hub.easystack.io/captain/
IMAGE_NAME=salmon
TAGS=${VERSION}
NAME=${REGISTRY}${IMAGE_NAME}:${TAGS}

# set go build opts
BUILD_OPTS=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
BUILD_OPTS_MAC=CGO_ENABLED=0 GOOS=darwin

default: build

build-image: build
	docker build -t ${NAME} .

build: fmt
	@echo "build $@"
	 ${BUILD_OPTS} go build ${LDFLAGS} -mod vendor -v -o ./bin/salmon ./

build-mac: fmt
	@echo "build $@"
	 ${BUILD_OPTS_MAC} go build ${LDFLAGS} -mod vendor -v -o ./bin/salmon ./

fmt:
	@echo "fmt $@"
	gofmt -l -w ./

clean:
	@echo "clean $@"
	rm -rf ./bin

