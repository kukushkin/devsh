#

GO = go
PROJECT_NAME = devsh
BUILD_DIR = build
MAIN = main.go

BUILD_FLAGS_DEBUG = -ldflags "-X main.debug=true"
BUILD_FLAGS_RELEASE = -ldflags "-X main.debug=false"

run:
	$(GO) run $(MAIN)

build: build-debug
build-debug:
	$(GO) build -o build/$(PROJECT_NAME) $(BUILD_FLAGS_DEBUG)

build-release:
	$(GO) build -o build/$(PROJECT_NAME) $(BUILD_FLAGS_RELEASE)