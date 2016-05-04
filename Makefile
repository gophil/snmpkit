.PHONY: all clean
# 被编译的文件
BUILDFILE = main.go

# 编译后的静态链接库文件
TARGETFILE = buildss/snmpdemo

# GOOS为目标主机系统, mac os则为 "darwin", window系列则为 "windows"
GOOS = darwin 
# GOARCH为目标主机CPU架构, 默认为amd64 
GOARCH= amd64

all: format test build clean

test:
	go test -v . 

format:
	gofmt -w .

build:
	mkdir -p builds
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -v -o $(TARGETFILE) $(BUILDFILE)

clean:
	go clean -i