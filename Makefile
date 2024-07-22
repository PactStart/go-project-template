# 设置变量
APP_NAME :=$(name)
ENV := $(env)
BIN_DIR := ./bin
ROOT_DIR := .
DOCKER_IMAGE_NAME_PREFIX=orderin-

# 导入版本文件
VERSION := $(shell cat version.txt)

# 默认目标
default: build

# 构建目标
build:
	@echo "Building..."
	@go build -o $(BIN_DIR)/orderin $(ROOT_DIR)/main.go

build_linux:
	@echo "Building..."
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o $(BIN_DIR)/xxxjz $(ROOT_DIR)/main.go

# 构建Docker镜像目标
docker_build:build_linux
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE_NAME_PREFIX)$(APP_NAME):$(VERSION) --build-arg APP_NAME=$(APP_NAME) --build-arg PROFILES=$(ENV) .

docker_build_all: build_linux
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE_NAME_PREFIX)admin-api:$(VERSION) --build-arg APP_NAME=admin-api --build-arg PROFILES=$(ENV) .
	@docker build -t $(DOCKER_IMAGE_NAME_PREFIX)app-api:$(VERSION) --build-arg APP_NAME=app-api --build-arg PROFILES=$(ENV) .

# 清理目标
clean:
	@echo "Cleaning..."
	@rm -rf $(BIN_DIR)

# 运行目标
run:
	@echo "Running..."
	@go run $(ROOT_DIR)/main.go $(APP_NAME)

# 帮助目标
help:
	@echo "Available targets:"
	@echo "  build     - Build the specified app"
	@echo "  docker_build  - Build Docker image"
	@echo "  clean     - Clean the build directory"
	@echo "  run       - Run the specified app"
	@echo "  help      - Show available targets"

.PHONY: default build clean package install run help