ROOT := $(shell pwd)

all: build

SOURCEDIR=./
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')
BINARY_NAME=docker-ecr-registry-plugin
LOCAL_BINARY=bin/local/$(BINARY_NAME)
LINUX_AMD64_BINARY=bin/linux-amd64/$(BINARY_NAME)
DARWIN_AMD64_BINARY=bin/darwin-amd64/$(BINARY_NAME)

.PHONY: build
build: $(LOCAL_BINARY)

$(LOCAL_BINARY): $(SOURCES) GITCOMMIT_SHA
	. ./scripts/shared_env && ./scripts/build_binary.sh ./bin/local $(VERSION) $(shell cat GITCOMMIT_SHA)
	@echo "Built $(BINARY_NAME)"

.PHONY: all-variants
all-variants: linux-amd64 darwin-amd64 windows-amd64

.PHONY: linux-amd64
linux-amd64: $(LINUX_AMD64_BINARY)
$(LINUX_AMD64_BINARY): $(SOURCES) GITCOMMIT_SHA
	./scripts/build_variant.sh linux amd64 $(VERSION) $(shell cat GITCOMMIT_SHA)

.PHONY: linux-amd64-image
linux-amd64-image: docker/Dockerfile.linux-amd64 linux-amd64
	./scripts/build_image.sh docker/Dockerfile.linux-amd64 drone-ecr-registry-plugin:$(shell cat GITCOMMIT_SHA)

.PHONY: darwin-amd64
darwin-amd64: $(DARWIN_AMD64_BINARY)
$(DARWIN_AMD64_BINARY): $(SOURCES) GITCOMMIT_SHA
	./scripts/build_variant.sh darwin amd64 $(VERSION) $(shell cat GITCOMMIT_SHA)

GITCOMMIT_SHA: $(GITFILES)
	git rev-parse --short=7 HEAD > GITCOMMIT_SHA
