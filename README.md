# 1、如何生成API文档
点击./main.go文件内的//go:generate

# 2、如何构建
make build
make build:linux

# 3、构建后如何通过命令启动
./bin/orderin admin-api -c ./config/config.yaml
./bin/orderin app-api  -c ./config/config.yaml


# 3、如何构建镜像
make docker_build name=admin-api env=test
make docker_build name=app-api env=test

# 4、容器内部排查
docker run --rm -it orderin-admin-api:1.0.0 sh