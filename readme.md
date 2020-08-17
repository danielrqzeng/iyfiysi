# iyfiysi
## pb依赖
### 安装protoc
#### win平台
* 下载pb编译工具protoc[win版本](https://repo1.maven.org/maven2/com/google/protobuf/protoc/3.9.0/protoc-3.9.0-windows-x86_64.exe)
* 可以移到${GO_PATH}/bin目录下，然后使用goland的控制台来使用命令
#### linux平台
* 下载pb编译工具protoc[linux版本](https://repo1.maven.org/maven2/com/google/protobuf/protoc/3.9.0/protoc-3.9.0-linux-x86_64.exe)
* 改名为protoc,添加到bin路径：`mv protoc-3.9.0-linux-x86_64.exe protoc`,`cp protoc /usr/bin/`
### 安装下载依赖go.mod
* `go mod download`
* 下载protoc-gen-go&运行依赖proto
	> `github.com/golang/protobuf/protoc-gen-go`，protoc-gen-go产生proto协议对应的go源码，以便使用
	> `github.com/golang/protobuf/proto`,此为产生的go源码需要的
	* `go get -d -u github.com/golang/protobuf`
	* `git checkout v1.3.2`
		> 由于太新版本会有兼容性问题，是以此处选择1.3.2版本
* 安装`protoc-gen-go`:`go install github.com/golang/protobuf/protoc-gen-go`