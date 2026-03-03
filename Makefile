# Makefile 用于简化常用命令
# 使用方法: make <命令名>

.PHONY: run build clean deps test help

# 默认目标：显示帮助信息
help:
	@echo "可用命令："
	@echo "  make run     - 运行应用（开发模式）"
	@echo "  make build   - 编译应用到 bin/server"
	@echo "  make clean   - 清理编译产物"
	@echo "  make deps    - 安装/更新依赖"
	@echo "  make test    - 运行测试"
	@echo "  make help    - 显示此帮助信息"

# 运行应用（开发模式）
# 直接运行 main.go，不需要编译
run:
	@echo "🚀 启动应用..."
	go run cmd/server/main.go

# 编译应用
# 编译后的可执行文件保存在 bin/server
build:
	@echo "🔨 编译应用..."
	@mkdir -p bin
	go build -o bin/server cmd/server/main.go
	@echo "✅ 编译完成: bin/server"

# 清理生成的文件
# 删除编译产物
clean:
	@echo "🧹 清理编译产物..."
	rm -rf bin/
	@echo "✅ 清理完成"

# 安装/更新依赖
# 下载并整理 go.mod 中的依赖
deps:
	@echo "📦 安装依赖..."
	go mod download
	go mod tidy
	@echo "✅ 依赖安装完成"

# 运行测试
# 执行所有测试用例
test:
	@echo "🧪 运行测试..."
	go test -v ./...
	@echo "✅ 测试完成"

# 格式化代码
# 使用 gofmt 格式化所有 Go 代码
fmt:
	@echo "✨ 格式化代码..."
	go fmt ./...
	@echo "✅ 格式化完成"

# 代码检查
# 使用 go vet 检查代码问题
vet:
	@echo "🔍 检查代码..."
	go vet ./...
	@echo "✅ 检查完成"

# 数据库迁移
# 运行数据库迁移（需要先启动应用）
migrate:
	@echo "📊 数据库迁移..."
	@echo "提示: 数据库迁移会在应用启动时自动执行"
	@echo "如需手动迁移，请运行应用"
