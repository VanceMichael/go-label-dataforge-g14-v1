# BENZHI_README

这是一个 Go 后端服务，用于DataForge 是面向资源提供方、审核员、数据产品开发者和场景运营者的纯后端平台，覆盖资源存证登记、公示赋码发证，以及目录授权、沙箱租约、模型测试和产品发布两条相互依赖的生命周期。

## 标准构建、运行和测试命令

进入容器后执行：

```bash
# 编译
cd '/app' && GOTOOLCHAIN=local go build ./...

# 启动
cd '/app' && GOTOOLCHAIN=local go run ./cmd/server

# 测试
cd '/app' && GOTOOLCHAIN=local go test ./...
```

## Docker 构建和进入容器

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh benzhi-task-417-amd64 linux/amd64
./build_benzhi_docker.sh benzhi-task-417-arm64 linux/arm64
docker run -it benzhi-task-417-amd64:latest
docker run -it --platform linux/arm64 benzhi-task-417-arm64:latest
```
