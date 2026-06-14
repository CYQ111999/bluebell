.PHONY: all build run test clean gotool help frontend-build dev start-all stop-all

# 变量定义
BINARY="bluebell_study"
CONFIG_FILE="conf/config.yaml"
FRONTEND_DIR="frontend"

# 默认目标：格式化代码并编译
all: gotool build

# 编译 Go 项目
build:
	@echo "Building..."
	go build -o ./bin/$(BINARY)

# 开发环境运行（带配置文件）
run:
	@echo "Running application..."
	go run main.go $(CONFIG_FILE)

# 生产环境编译（Linux amd64，优化体积）
build-prod:
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/$(BINARY)

# 运行测试
test:
	@echo "Running tests..."
	go test -v ./...

# 代码格式化和检查
gotool:
	@echo "Running go fmt..."
	go fmt ./
	@echo "Running go vet..."
	go vet ./

# 清理编译产物
clean:
	@echo "Cleaning..."
	@if [ -f bin/$(BINARY) ]; then rm -rf bin; fi
	@go clean

# 构建前端项目
frontend-build:
	@echo "Building frontend..."
	cd $(FRONTEND_DIR) && npm install && npm run build

# 启动前端开发服务器
frontend-dev:
	@echo " Starting frontend dev server..."
	cd $(FRONTEND_DIR) && npm run dev

# 同时启动前后端（Windows PowerShell）- 后端使用 Air
start-all:
	@echo "========================================="
	@echo "  🚀 启动前后端开发服务"
	@echo "========================================="
	@echo ""
	@powershell -Command "Start-Process powershell -ArgumentList '-NoExit', '-Command', 'cd $(FRONTEND_DIR); npm run dev' -WindowStyle Normal"
	@echo "✅ 前端服务已在新窗口启动 (http://localhost:3000)"
	@echo "⏳ 5秒后启动后端服务（支持热重载）..."
	@powershell -Command "Start-Sleep -Seconds 5"
	@echo "🔥 启动后端服务（Air）..."
	@air

# 停止所有服务（需要手动 Ctrl+C）
stop-all:
	@echo "请手动关闭各个终端窗口，或按 Ctrl+C 停止服务"

# 同时构建前后端
build-all: gotool frontend-build build

# 查看帮助信息
help:
	@echo "可用命令："
	@echo "  make          - 格式化代码并编译"
	@echo "  make build    - 编译生成二进制文件"
	@echo "  make run      - 运行后端应用"
	@echo "  make frontend-dev   - 启动前端开发服务器"
	@echo "  make start-all      - 同时启动前后端（后端使用 Air）"
	@echo "  make build-prod - 编译生产环境版本（Linux）"
	@echo "  make test     - 运行测试"
	@echo "  make gotool   - 代码格式化与检查"
	@echo "  make clean    - 清理编译产物"
	@echo "  make frontend-build - 构建前端项目"
	@echo "  make build-all - 构建前后端全部"
