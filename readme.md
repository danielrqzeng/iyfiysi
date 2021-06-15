# iyfiysi教程文档
* [iyfiysi](#iyfiysi)
   * [快速开始](#快速开始)
      * [<a href="#protoc">1.安装protoc</a>](#1安装protoc)
      * [2.安装iyfiysi和依赖](#2安装iyfiysi和依赖)
      * [3.项目生成](#3项目生成)
      * [4.业务实现](#4业务实现)
      * [5.服务启动](#5服务启动)
   * [架构解析](#架构解析)
      * [服务架构](#服务架构)
   * [<a href="">框架关键技术说明</a>](#框架关键技术说明)
      * [配置中心](#配置中心)
      * [服务治理](#服务治理)
      * [链路追踪](#链路追踪)
         * [链路追踪配置](#链路追踪配置)
         * [链路追踪内容](#链路追踪内容)
      * [API文档](#api文档)
      * [中间件（拦截器）](#中间件拦截器)
         * [网关（gateway）](#网关gateway)
         * [服务器（server）](#服务器server)
      * [<a href="">监控</a>](#监控)
         * [基础概念](#基础概念)
         * [监控安装](#监控安装)
   * [高端制造](#高端制造)
      * [项目生成逻辑](#项目生成逻辑)
      * [服务生成逻辑](#服务生成逻辑)
---
# iyfiysi

**iyfiysi**是一个生成一个简单易用的分布式框架工具。
通过iyfiysi生成的是一个依赖少，易于快速扩展，提供api服务的框架。其基于[grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)，集成了服务治理，配置管理，鉴权，限流，rpc服务，链路追踪，监控告警等特性于一体的高可用，高可靠，可扩展的api服务框架

iyfiysi生成的框架优点在于
* 多端开发，框架集成工具可以使得其在主流操作系统都可以直接一致地运行
* 业务框架是进程式的框架，暂时不启用docker方式部署，大大减少所需要的知识储备
* 适用于特别需要扩展性的场景，其扩展起来只需要不停加进程即可
* devops式地集成了监控，链路追踪等逻辑，利于对业务进行监控和评估

iyfiysi框架适用的场景
* api接口
* 需要不定时不定量地扩展
* 对接口数据进行监控管理

## 快速开始
---
### [1.安装protoc](#protoc)
protoc是一个由proto文件生成各种语言数据接口的工具，此项目使用的是3.9.0的版本的protoc
* 下载`protoc`
  * [win版本](https://repo1.maven.org/maven2/com/google/protobuf/protoc/3.9.0/protoc-3.9.0-windows-x86_64.exe)
  * [linux版本](https://repo1.maven.org/maven2/com/google/protobuf/protoc/3.9.0/protoc-3.9.0-linux-x86_64.exe)
  * [osx版本](https://repo1.maven.org/maven2/com/google/protobuf/protoc/3.9.0/protoc-3.9.0-osx-x86_64.exe)
* 放置于${GO_PATH}/bin目录(window平台，务必保证其为`protoc.exe`,linux或者mac平台保证为`protoc`)
  > 需添加${GO_PATH}/bin为系统执行路径

---
### 2.安装iyfiysi和依赖
要求go版本>=1.13
* 进入到GO_PATH目录中，`cd GO_PATH`
* `go get github.com/RQZeng/iyfiysi`
* `cd github.com/RQZeng/iyfiysi`
* 在linux|mac中安装[sh install.sh]()
  ```sh
  # 安装完毕，效果如下
  [root@VM_0_14_centos bin]# ll protoc*
  ...
  -rwxr-xr-x 1 root root 15520427 Jun 11 17:43 /data/go_path/bin/iyfiysi
  -rwxr-xr-x 1 root root  6347332 Jun 11 17:43 /data/go_path/bin/protoc
  -rwxr-xr-x 1 root root  6347332 Jun 11 17:43 /data/go_path/bin/protoc-gen-go
  -rwxr-xr-x 1 root root  7464659 Jun 11 17:43 /data/go_path/bin/protoc-gen-grpc-gateway
  -rwxr-xr-x 1 root root 16758590 Jun 11 17:43 /data/go_path/bin/protoc-gen-iyfiysi
  -rwxr-xr-x 1 root root  7174301 Jun 11 17:43 /data/go_path/bin/protoc-gen-swagger
  ...
  ```
* 在window系统中安装[install.bat]()
  ```sh
  # 安装完毕，效果如下
  PS F:\go_path\bin> ls
  ...
  -a----        2021/6/11     17:40        6501376 protoc-gen-go.exe
  -a----        2021/6/11     17:40        7615488 protoc-gen-grpc-gateway.exe
  -a----        2021/6/11     17:40       16616448 protoc-gen-iyfiysi.exe
  -a----        2021/6/11     17:40        7955968 protoc-gen-swagger.exe
  -a----        2021/6/11     17:40        7955968 protoc.exe
  -a----        2021/6/11     17:40       15400448 iyfiysi.exe
  ...
  ```
---

### 3.项目生成
* 项目标识：项目使用组织和app名称来标识一个项目
  * 组织：一个域名，比如 iyfiysi.com
  * app名称:项目名称，比如[test]()
  > 此标识是项目里面非常重要的，需要做成唯一可识
* 此处假设我们启动项目于目录`/data/project`,使用组织标识为`iyfiysi.com`,app为`test`
  > 对于创建目录的要求，其不在${GO_PATH}/src中即可
* 进入目录中`cd /data/project`
* 使用命令，新建项目:[iyfiysi new -o iyfiysi.com -a test]()
* 新建完成，可以在目录`/data/project/iyfiysi.com/test/`看到生成的项目布局如下：
  ```sh
  [root@VM_0_14_centos test]# tree -L 3
  .
  |-- build.sh #构建脚本
  |-- cmd
  |   |-- conf #配置进程
  |   |   `-- main.go
  |   |-- gateway #gateway进程
  |   |   `-- main.go
  |   `-- server #server进程
  |       `-- main.go
  |-- conf #项目配置
  |   `-- app.yaml
  |-- go.mod
  |-- go.sum
  |-- internal
  |   |-- app
  |   |   |-- gateway #cmd.gateway-->app.gateway
  |   |   `-- server #cmd.server->app.server
  |   |       `-- service #业务实现
  |   `-- pkg
  |       |-- conf
  |       |-- data
  |       |-- db
  |       |-- governance #服务治理
  |       |-- interceptor #中间件
  |       |-- logger
  |       |-- trace
  |       `-- utils
  |-- keystore #http2的自签名证书
  |   |-- ca.crt
  |   |-- ca.key
  |   |-- grpc.crt
  |   |-- grpc.csr
  |   `-- grpc.key
  |-- LICENSE
  |-- logs
  |-- metric #监控相关
  |   |-- confd  # confd的配置和启动脚本
  |   |   |-- conf.d # 配置
  |   |   |-- once.sh
  |   |   |-- templates # 生成模版
  |   |   `-- watch.sh
  |   |-- grafana # grafana的dashbord 配置
  |   |   |-- iyfiysi.json # 项目监控
  |   |   |-- node.json #节点监控
  |   |   `-- process.json # 进程监控
  |   |-- process-exporter # 进程监控的配置和启动脚本
  |   |   |-- process.yml
  |   |   `-- setup.sh
  |   `-- prometheus # prometheus的启动配置和监控源（机器和进程）的配置
  |       |-- node.yaml
  |       |-- process.yaml
  |       `-- prometheus.yml
  |-- proto
  |   |-- gen.sh
  |   |-- google
  |   |   |-- api
  |   |   |-- protobuf
  |   |   `-- rpc
  |   `-- service.proto #对外服务的定义的pb文件
  |-- test_conf
  |-- test_gateway
  `-- test_server
  ```

* 编译：`cd iyfiysi.com/test;`
  * linux|mac:[sh build.sh]()
  * window:[build.bat]()
* 编辑生成三个bin文件，分别为
  * test_conf：配置中心，其功能是将配置上传到etcd，配置文件对应的是`conf/app.yaml`
  * test_gateway：网关服务器
  * test_server：业务服务器

---

### 4.业务实现
通过在pb文件定义我们对外的服务，并且在业务服务器中实现该逻辑，以下举例怎么一步步实现一个[加法]()业务
* 编辑pb文件如下
  ```diff
  //vim proto/service.proto

  + //SumRequest ...
  + message SumRequest {
  +     uint64 val1 = 1;
  +     uint64 val2 = 2;
  + }         

  + //SumResponse ...
  + message SumResponse {
  +     uint64 sum = 1;
  + }

  service testService {
  +    //@add sum rpc
  +    rpc Sum(SumRequest) returns (SumResponse) {
  +        option (google.api.http) = {
  +            post: "/v1/sum"
  +            body: "*"                                                         
  +         }
  +     }
  }
  ```
* 构建：`sh build.sh`，构建之后，业务逻辑会生成在`internal/app/server/service`中
* 做业务逻辑的实现
  ```diff
  //vim internal/app/server/service/test_service_sum.go
  
  // Code generated by protoc-gen-iyfiysi at 2021 Jun 10
  // DO AS YOU WANT

  package service

  import(
      "context"
      pb "iyfiysi.com/test/proto"
  )
  
  // Sum ... @TODO
  func (s *TestServiceImpl) Sum(
    ctx context.Context,req *pb.SumRequest)(//request param                                                   
      rsp *pb.SumResponse, err error) { //response
      
      rsp =&pb.SumResponse{}
  +    rsp.Sum = req.Val1 + req.Val2 //@add 实现业务逻辑-加法
      return
  }
  ```
* 编译：`sh build.sh`

### 5.服务启动
* 准备好etcd服务器，比如其在`http://127.0.0.1:2379`
  > 若是没有etcd，本文也准备了一个[简单启动etcd的教程](http://github.com/iyfiysi/blob/master/etcd.md)，可以现搭起来一个
* 修改conf/app.yaml
  ```diff
  ...
  # etcd,不支持即改即生效
  etcd:
    enable: true #是否开启etcd服务，目前只能开启
    metricKey: "/iyfiysi.com/test/metric" #服务监控的key
    serviceKey: "/iyfiysi.com/test/service" #注册服务的key
    swaggerKey: "/iyfiysi.com/test/swagger" #文档服务的key
    etcdServer:
  +    - "http://127.0.0.1:2379" # @modify 修改为etcd的服务接口

  ...
  ```
* 启动服务
  * 配置中心将配置文件上传到etcd:`./test_conf`
  * 启动gateway：`./test_gateway`
  * 启动server：`./test_server`
  * 启动完成之后，查看gateway侦听了那些接口，在nginx配置个反向代理，即可对外提供服务了
* 端口服务说明
  <center>
  
  ![Untitled Diagram-port_1_.png](https://www.hualigs.cn/image/60c821bb6b945.jpg)
  </center>

  * 8080~8085:swagger接口，默认是不开启，需要在配置中开启才提供
  * 8000~8050:gateway对外服务将会在这些端口中选择可用端口，提供服务
  * 30000~30500:server向gateway提供的接口，启动时候选择可用的
  * 41000~41500:gateway向prometheus提供的监控接口
  * 42000~42500:server向prometheus提供的监控接口
* 测试：`curl --location --request POST 'http://172.30.0.14:8000/v1/sum' --header 'Content-Type: applicati/json' --data '{"val1":100,"val2":200}'`
* 当然也可以开启了swagger接口服务之后，通过swagger查看
  ![image.png](https://i.loli.net/2021/06/13/q8WSLrhGoXpBnYT.png)
---


## 架构解析
由以上我们可以知道，5个步骤即可将框架部署启动完毕，业务逻辑实现起来是非常简易的，只需要定义pb和实现业务逻辑即可。下面将介绍一下iyfiysi生成的项目是如何运作的
### 服务架构
![](https://www.hualigs.cn/image/60c82e0fa3a63.jpg)
* 图中虚框都是可选的服务，主要是服务监控和链路追踪部分
* 实框中主要有etcd服务，提供配置中心和服务治理的功能
* 项目编译出来三个业务进程
  * conf:配置管理，其作用是将配置上传到etcd，服务首次启动时候需要上传或者配置变更了也需要上传配置文件
  * gateway：网关，为外部请求提供入口，其集成服务发现，频率限制，链路追踪，监控等等功能
  * server：业务实现，根据pb定义接口，做业务的实现，其集成服务注册链路追踪，监控等等功能


---
## [框架关键技术说明]()
### 配置中心
* etcd本质上是一个kv数据库，带有保活租赁，前缀侦听等功能，是以其是合适做配置中心，存放配置信息的
* 本项目使用etcd作为配置中心，通过进程conf业务配置到etcd，以供gateway和server使用
* gateway和server进程，启动时指定etcd服务器和配置key，即可启动，其将会读取etcd的远程配置信息启动程序
* 配置信息什么时候上传？
  * 首次启动时候，必须先通过conf进程上传配置信息`conf/app.yaml`到etcd
  * 配置有变动时候，也可以通过conf上传
  > 配置信息上传之后，在gateway和server是即时生效的，同时也会有事件通知其变动情况
  > 配置生效不同于配置对应的逻辑生效
  > 比如侦听了一个服务端口在8090，后配置信息修改为8091，虽然配置生效了，但是无法更改配置对应的侦听端口
  > 比如一个限制最大次数的值，使用时从viper读取，若是此配置更改了，此配置对应的逻辑也会生效

进程怎么使用
* 配置是以一个kv的形式，保存在etcd中
* 三个进程（conf,server,gateway）,启动的命令行参数，可以指明etcd服务器，key
* 若是未指明，则使用默认值，理论上只需要指明etcd服务器即可，key是项目创建时用户指定的组织和app名组成

---
### 服务治理
服务治理的调用关系说明
* 服务提供者在启动时，向注册中心注册自己提供的服务
* 服务消费者在启动时，向注册中心订阅自己所需的服务
* 注册中心返回服务提供者地址列表给消费者，如果有变更，注册中心将基于长连接推送变更数据给消费者
* 服务消费者，从提供者地址列表中，基于负载均衡算法，选一台可用给消费者进行调用

本框架的服务治理详情如下：
* etcd提供前缀+租赁保活的方式，通过这种方式，可以实现一个即插即用的容易scale的服务集群
* 本框架中，使用etcd作为服务治理的服务治理中心，框架的服务治理主要有三个部分组成
  * **对外业务服务**：由业务服务器提供服务（服务注册），网关对服务进行发现和使用（服务订阅和发现）
  * **对内grpc服务**：[@TODO]()，业务需求决定是否需要其他的grpc的服务
  * **监控服务**：有业务服务器和网关服务器提供服务，confd对服务进行发现和使用
  * **文档服务**：swagger服务，此服务是可选的，由开关控制。开启开关后，网关服务器提供服务

本框架的服务治理配置信息
```yaml
# vim conf/app.yaml

# ...
# etcd,不支持即改即生效
etcd:
  enable: true #是否开启etcd服务，目前只能开启
  metricKey: "/iyfiysi.com/test/metric" #服务监控的key,iyfiysi.com/test为项目生成时候，用户传进来的组织和app名称
  serviceKey: "/iyfiysi.com/test/service" #注册服务的key
  swaggerKey: "/iyfiysi.com/test/swagger" #文档服务的key                       
  etcdServer:
    - "http://127.0.0.1:2379"
# ...

```
---
### 链路追踪
基于`jaeger`的链路追踪，使得请求一目了然。在项目中，每个请求都被trace记录，并且上报jeager服务后台，使用jaeger的服务后台，即可查看请求链路情况
> 需要预先准备好jaeger服务后台，若是没有，本项目也提供了一个快速搭建[jaeger服务后台](http://github.com/iyfiysi/blob/master/jaeger.md)的方式
#### 链路追踪配置
```yaml
# vim conf/app.yaml

#...
# jaeger,不支持配置即改即生效
jaeger:
  enable: true
  jaegerServer:
    - "localhost:6831"
#...
```
* 链路追踪默认是关闭的，因为其有性能损耗，需要用户根据业务情况，自个开启
#### 链路追踪内容
* 项目中，默认记录了以下的span
  * http的span
  * grpc的span
  > 以上span记录了名称，路径，耗时等等信息
* 项目中还提供了以下的span以供具体业务使用
  * mysql的span
  * redis的span
* 另外，项目中，暂时并没有实现函数级别的span，其已经列为TODO计划中，目前追求是在编译侧直接搞定，不需要人工做任何开发，是以使用该框架的用户，无需自行做函数级别的span
* 监控的效果图如下：
![](https://www.hualigs.cn/image/60c856004bbad.jpg)
---

### API文档
基于pb接口定义，生成了swagger的接口文档，以供开发者更好地对接
* 文档服务的配置信息如下
```yaml
# vim conf/app.yaml

#...
# swagger服务
swagger:  
  enable: true
  minPort: 8080
  maxPort: 8085
  ignoreIP:
    *ignoreIPRef
  potentialIP:
    *potentialIPRef
  path: "/swagger/"
#...
```
* 默认情况下，文档服务是关闭的，只有在开发对接情况下才开放，生产环境务必关闭之
* 以上可知，其服务接口为8080~8085,启动后可以通过注册中心查看
* 其web地址为：http://<your_server_addr>:8080/swagger/

---
### 中间件（拦截器）
拦截器是一种共性控制类的功能，在实际业务处理之前，对请求进行验证。拦截器的代码放置于`pkg/intercepor`中
<center>


![](https://www.hualigs.cn/image/60c872e948005.jpg)
</center>

#### 网关（gateway）
* 日志拦截器：记录请求的审计日志
* 监控拦截器：对请求进行监控统计，并且将这些数据保存，以便prometheus拉取
* 重试拦截器：是否重试
* 限流拦截器：[动态限流](https://fredal.xin/netflix-concuurency-limits),基本原理是基于tcp的拥塞控制
* 链路追踪拦截器：链路追踪
#### 服务器（server）
* 日志拦截器：记录请求的审计日志
* 监控拦截器：对请求进行监控统计，并且将这些数据保存，以便prometheus拉取
* 认证拦截器：验证是否请求端（gateway）是否有权限调用服务，使用的是token的方式校验
* 恢复拦截器：异常恢复
* 链路追踪拦截器：链路追踪

---
### [监控]()
* 监控使用的是[prometheus]()作为数据收集和[grafana]()作为数据展示和管理
* 监控是可选的，通过控制开关来控制是否开启,默认是关闭
* 监控分为三个维度，分别是：
    * [机器监控]()：对机器的指标进行监控，cpu,mem,io等
    * [进程监控]()：通过进程名，对某些进程进行监控，cpu,mem,io等指标
    * [业务监控]()：对业务进行监控，比如业务的qps，耗时等等
* 在项目中的配置如下
```yaml
# conf/app.yaml
#...
metrics:
  enable: true # 是否开启监控
  gateway:
    path: "/metrics"
    minPort: 41000
    maxPort: 41500
    ignoreIP:
      *ignoreIPRef
    potentialIP:
      *potentialIPRef
  server:
    path: "/metrics"
    minPort: 42000
    maxPort: 42500
    ignoreIP:
      *ignoreIPRef
    potentialIP:
      *potentialIPRef
#...
```
#### 基础概念
* 指标源：产生和提供指标的服务进程，一般指各种exporter，各种业务服务器等
* 指标收集服务器：此处是prometheus
* prometheus是通过定时向指标源拉取的方式获取指标，并且保存展示

指标系统构成
    ![指标系统构成](https://www.hualigs.cn/image/60c09a9aa5bfc.jpg)


#### 监控安装
* 需要预先安装**prometheus**和**grafana**服务，本项目也提供了一个快速搭建[prometheus服务后台](http://github.com/iyfiysi/blob/master/prometheus.md)的方式
* 此处假设安装后
  * prometheus的`file_sd`目录为`/data/docker/metrics/prometheus/config/file_sd`
  * grafana的web地址为`http://<out_ip>:3000
  ---
* 启动[机器监控]()
    * 下面实例中，是运行在linux机器的，[其他环境下载地址](https://github.com/prometheus/node_exporter/releases/)
    * 机器监控，使用的是[node_exporter](https://github.com/prometheus/node_exporter/releases/download/v0.18.1/node_exporter-0.18.1.linux-amd64.tar.gz)
    * `cd metric/node_exporter`，下载&解压
    * 启动：`sh setup.sh`
      > 以上启动命令，可知此指标源的地址是`XXX.XXX.XXX.XXX:9100`
    * 监控发现：
    框架生成机器监控的配置
      ```diff
      - metric/prometheus/node.yaml 
      # 1.set up node_exporter and get the listen addr
      # 2.modify $targets to the node_exporter listen addr                            
      # 3.put me to ${PATH_TO_PROMETHEUS}/config/file_sd
      # gen by iyfiysi at 2021 Jun 15
      - labels:
          project: "/iyfiysi.com/test"
          role: "node"
        targets:
      +    - "172.30.0.14:9100" #修改为node_exporter侦听地址,监控多少个机器，次第添加即可
      ```
    * 将上面修改后的配置，放置于prometheus的文件服务发现目录`/data/docker/metrics/prometheus/config/file_sd`
    ---
* 启动[进程监控]()
    * 下面实例中，是运行在linux机器的，[其他环境下载地址](https://github.com/ncabatoff/process-exporter/releases)
    * 进程监控，使用的是[process-exporter](https://github.com/ncabatoff/process-exporter/releases/download/v0.7.5/process-exporter-0.7.5.linux-amd64.tar.gz)
    * `cd metric/process-exporter`，下载&解压
    * 启动：`sh setup.sh`
        > 以上启动命令，可知此指标源的地址是`XXX.XXX.XXX.XXX:9256`
    * 监控发现：
    框架生成进程监控的配置
      ```diff
      - metric/prometheus/process.yaml 
      # 1.set up process-exporter and get the listen addr                             
      # 2.modify $targets to the process-exporter listen addr
      # 3.put me to ${PATH_TO_PROMETHEUS}/config/file_sd
      # gen by iyfiysi at 2021 Jun 15
      - labels:
          project: "/iyfiysi.com/test"
          role: "process"
        targets:
      +    - "172.30.0.14:9256" #修改为process-exporter侦听地址
      ``` 
    * 将上面修改后的配置，放置于prometheus的文件服务发现目录`/data/docker/metrics/prometheus/config/file_sd`
  ---
* 启动[业务监控]()
    * 由于业务服务器启动多少是根据业务而定，是以其势必会有业务服务器乱启动，乱关停
    * 通过服务发现，将当前运行业务服务器找到，并且告知prometheus这些指标源
    * 使用[confd](https://github.com/kelseyhightower/confd/releases/download/v0.16.0/confd-0.16.0-linux-amd64)，通过获取在ectd中注册的监控服务实例，来生成监控服务器的指标源
    * 启动配置已经由iyfiysi新建项目时候生成
        ```sh
        metric/confd/
        |-- conf.d
        |   `-- confd.toml
        |-- templates
        |   `-- confd_yaml.tmpl
        |-- once.sh # 单次生成监控服务的实例地址
        `-- watch.sh #watch的模式生成监控服务的实例地址
        ```
    * 将生成目录修正为prometheus的指标源文件目录
        ```diff
        - metric/confd/conf.d/confd.toml
        # gen by iyfiysi at 2021 Jun 15
        # confd config file for iyfiysi.com/test

        [template]
        src = "confd_yaml.tmpl" # templates/confd_yaml.tmpl
        + dest = "/data/docker/metrics/prometheus/config/file_sd/test.yaml"
        keys = [
            "/iyfiysi.com/test/metric",
        ]
        ```
    * 启动服务`sh watch.sh`
* 指标展示grafana
    * 导入数据源`http://prometheus:9091`
    * 导入dashbord配置文件，其是有iyfiysi生成于`metric/grafana`
        * node.json：机器监控dashboard，[来源](https://grafana.com/grafana/dashboards/8919)
        * process.json: 进程监控dashboard，[来源](https://grafana.com/grafana/dashboards/249)
        * iyfiysi.json: 业务监控dashbord
    * 效果展示之-机器监控
    ![效果图](https://i.loli.net/2021/04/14/xP7D9bS16fFopkn.png)
    * 效果展示之-进程监控
    ![效果图](https://i.loli.net/2021/04/15/YfN3r4JdV1ceX9o.png)
    * 效果展示之-业务监控
    ![效果图](https://www.hualigs.cn/image/60c0ad8104fe3.jpg)


---
## 框架样例
* [短网址服务](https://github.com/RQZeng/short_url)
  * 一个基于iyfiysi框架开发的短网址服务，提供短域名编码，短语编码，短域名跳转，禁用等api服务
  * [业务体验](https://surl4.me/),下图为业务UI
  ![](https://www.hualigs.cn/image/60c8a321b874b.jpg)

---
## 高端制造
### 项目生成逻辑
![](https://i0.hdslb.com/bfs/album/e5af3a3404ea5b34b182a6d27aa887750eb73d6c.png)
### 服务生成逻辑
![](https://www.hualigs.cn/image/60c0aeb81bc99.jpg)

---
