EXECUTABLES = git go find pwd
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

DOCKER_IMAGE=mnbbrown/flagger
BINARY=flagctl
CMD_PATH=./cmd/flagctl
VERSION=0.1.0
BUILD=`git rev-parse HEAD`
PLATFORMS=darwin linux
ARCHITECTURES=386 amd64

default: build

all: clean build_all

setup:
	mkdir -p __dist

run: build
	./__dist/$(BINARY) serve

build: setup
	go build ${LDFLAGS} -o __dist/${BINARY} ${CMD_PATH}

build_all: setup
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -v -o __dist/$(BINARY)-$(GOOS)-$(GOARCH) $(CMD_PATH))))

# Remove only what we've created
clean:
	find ${ROOT_DIR}/__dist -name '${BINARY}[-?][a-zA-Z0-9]*[-?][a-zA-Z0-9]*' -delete

docker: build_all build_web
	docker build -f Dockerfile -t $(DOCKER_IMAGE):latest __dist/


clean_web:
	rm -rf __dist/ui

build_web: clean_web
	cd ui && PUBLIC_URL=/ui NODE_ENV=production npm run build
	cp -R ui/build __dist/ui

test:
	go test ./pkg ./cmd/flagctl ./client

.PHONY: check clean install build_all all
