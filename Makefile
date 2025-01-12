.PHONY: build clean test

# 项目名称
BINARY_NAME=aicm

# 构建目标
build:
	go build -o $(BINARY_NAME) ./cmd/aicm

# 清理构建文件
clean:
	rm -f $(BINARY_NAME)

# 运行测试
test:
	go test ./...

# 安装到系统路径
install: build
	mv $(BINARY_NAME) /usr/local/bin/
# install to go bin
i: 
	go install ./...
# 格式化代码
fmt:
	go fmt ./...

# 检查代码风格
lint:
	golangci-lint run

# 运行程序
run: build
	./$(BINARY_NAME)
