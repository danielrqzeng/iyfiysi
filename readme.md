# iyfiysi
## quick start
### 安装protoc
#### win平台
* 下载pb编译工具protoc [win版本](https://repo1.maven.org/maven2/com/google/protobuf/protoc/3.9.0/protoc-3.9.0-windows-x86_64.exe)
* 可以移到${GO_PATH}/bin目录下，然后使用goland的控制台来使用命令
#### linux平台
* 下载pb编译工具protoc [linux版本](https://repo1.maven.org/maven2/com/google/protobuf/protoc/3.9.0/protoc-3.9.0-linux-x86_64.exe)
* 改名为protoc,添加到bin路径：`mv protoc-3.9.0-linux-x86_64.exe protoc`,`cp protoc /usr/bin/`
### 安装iyfiysi和依赖
* `go get github.com/RQZeng/iyfiysi`
* `go mod download`
* 在linux系统中安装`sh build.sh`
### 创建项目
* `cd /data/project`
* `iyfiysi new -n test.com -a test`:项目将会构建在目录`/data/project/test.com/test`中
* 启动etcd
* 编译项目`cd /data/project/test.com/test`;`sh build.sh`
* 运行
    * `./test_server`
    * `./test_gateway`