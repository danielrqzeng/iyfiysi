# iyfiysi
iyfiysi是一个生成一个简单易用的分布式框架工具。
通过iyfiysi生成的是一个依赖少，易于快速扩展，提供api服务的框架。其基于[grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)，集成了服务治理，配置管理，鉴权，rpc服务，链路追踪，监控告警等特性于一体的高可用，高可靠，可扩展的api服务框架



## quick start
---
### 1.安装protoc
#### win平台
* 下载pb编译工具`protoc`
  * [win版本](https://repo1.maven.org/maven2/com/google/protobuf/protoc/3.9.0/protoc-3.9.0-windows-x86_64.exe)
  * [linux版本](https://repo1.maven.org/maven2/com/google/protobuf/protoc/3.9.0/protoc-3.9.0-linux-x86_64.exe)
  * [osx版本](https://repo1.maven.org/maven2/com/google/protobuf/protoc/3.9.0/protoc-3.9.0-osx-x86_64.exe)
* 放置于${GO_PATH}/bin目录(window平台，务必保证其为`protoc.exe`,linux或者mac平台保证为`protoc`)
---
### 2.安装iyfiysi和依赖
要求go版本>=1.13,安装`iyfiysi`
* 进入到GO_PATH目录中，`cd GO_PATH`
* `go get github.com/RQZeng/iyfiysi`
* `cd github.com/RQZeng/iyfiysi`
* 在linux|mac中安装[sh build.sh]()
  ```sh
  #效果如下
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
* 在window系统中安装[build.bat]()
  ```
  # 效果如下
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
  * app名称:项目名称，比如test
  > 此标识是项目里面非常重要的，需要做成唯一可识
* 此处假设我们启动项目于目录`/data/project`,使用组织标识为`iyfiysi.com`,app为`test`,`cd /data/project`
* 新建项目:[iyfiysi new -o iyfiysi.com -a test]()
* 新建完成，可以看到生成的项目布局如下：
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
  |   |-- confd # confd的配置和配置生成
  |   |   |-- conf.d
  |   |   `-- templates
  |   `-- grafana # grafana的dashbord 配置
  |       |-- iyfiysi.json
  |       |-- node.json
  |       `-- process.json
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
### 4.etcd服务
* etcd在项目中，作用是配置中心和服务治理，是项目中必不可少的依赖
* 使用改项目需要先准备好etcd服务，或者按照以下步骤启动一个etcd服务
* 假设在目录`/data/docker/etcd`中启动一个etcd的镜像服务（单节点）
* `mkdir /data/docker/etcd/etcd-data`
* 启动配置
```yaml
#/data/docker/etcd/docker-compose.yml
etcd:
    image: 'quay.io/coreos/etcd:v3.1.7'
    restart: always
    ports:
        - '2379:2379'
        - '2380:2380'
        - '4001:4001'
    environment:
        - TZ=CST-8
        - LANG=zh_CN.UTF-8
    command:
        /usr/local/bin/etcd
        -name etcd0
        -data-dir /etcd-data
        -advertise-client-urls http://172.30.0.14:2379,http://172.30.0.14:4001
        -listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001
        -initial-advertise-peer-urls http://172.30.0.14:2380
        -listen-peer-urls http://0.0.0.0:2380
        -initial-cluster-token docker-etcd
        -initial-cluster etcd0=http://172.30.0.14:2380
        -initial-cluster-state new
    volumes:
        - /data/docker/etcd/etcd-data:/etcd-data
```
> 以上可知
> 服务端口为http://172.30.0.14:2379，http://172.30.0.14:4001
> 集群端口为http://172.30.0.14:2380

---
### 5.项目生成
* 项目标识：项目使用组织和app名称来标识一个项目
  * 组织：一个域名，比如 iyfiysi.com
  * app名称:项目名称，比如test
  > 此标识是项目里面非常重要的，需要做成唯一可识
* 此处假设我们启动项目于目录`/data/project`,使用组织标识为`iyfiysi.com`,app为`test`
* 新建项目:`cd /data/project`;`iyfiysi new -o iyfiysi.com -a test`
* 新建完成，可以看到生成的项目布局如下：
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
  |   |-- confd # confd的配置和配置生成
  |   |   |-- conf.d
  |   |   `-- templates
  |   `-- grafana # grafana的dashbord 配置
  |       |-- iyfiysi.json
  |       |-- node.json
  |       `-- process.json
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

* 编译：`cd /data/project/iyfiysi.com/test;sh build.sh`
* 编辑生成三个bin文件，分别为
  * test_conf：配置中心，其功能是将配置上传到etcd，配置文件对应的是`conf/app.yaml`
  * test_gateway：网关服务器
  * test_server：业务服务器
---
### 6.业务实现
通过在pb文件定义我们对外的服务，并且在业务服务器中实现该逻辑，以下举例怎么一步步实现一个[加法]()业务
* 编辑pb文件如下
  ```js
  //vim proto/service.proto

  //...

  //@add sum req&rsp
  message SumRequest {
      uint64 val1 = 1;
      uint64 val2 = 2;
  }         

  message SumResponse {
      uint64 sum = 1;
  }  
  //...

  service testService {
      // ...   

      //@add sum rpc
      rpc Sum(SumRequest) returns (SumResponse) {
          option (google.api.http) = {
              post: "/v1/sum"
              body: "*"                                                         
          };
      }     
      //...
  }         
  ```
* 构建：`sh build.sh`，构建之后，业务逻辑会生成在`internal/app/server/service`中
* 做业务逻辑的实现
  ```go
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
      rsp.Sum = req.Val1 + req.Val2 //实现业务逻辑-加法
      return
  }
  ```
* 编译：`sh build.sh`
* 启动服务
  * 配置中心将配置文件上传到etcd:`./test_conf`
  * 启动gateway：`./test_gateway`,其端口可以在配置`conf/app.yaml`.gateway中找到，目前默认是8000~8050
  * 启动server：`./test_server`
* 测试：`curl --location --request POST 'http://172.30.0.14:8000/v1/sum' --header 'Content-Type: applicati/json' --data '{"val1":100,"val2":200}'`
---

## 架构解析
由以上我们可以知道，除了初期的搭建环境比较繁琐之外，业务逻辑实现起来是非常简易的，只需要定义pb和实现业务逻辑即可。下面将介绍一下iyfiysi生成的项目是如何运作的
### 服务架构
![](https://www.hualigs.cn/image/60c0ae1f09bf5.jpg)
* 图中虚框都是可选的服务，主要是服务监控和链路追踪部分
* 实框中主要有etcd服务，提供配置中心和服务治理的功能
* 项目编译出来三个业务进程
  * conf:配置管理，其作用是将配置上传到etcd，服务首次启动时候需要上传或者配置变更了也需要上传配置文件
  * gateway：网关，为外部请求提供入口，其集成服务发现，频率限制，链路追踪，监控等等功能
  * server：业务实现，根据pb定义接口，做业务的实现，其集成服务注册链路追踪，监控等等功能

### 配置中心
* etcd本质上是一个kv数据库，还带有版本控制的功能，是以其是合适做配置中心，存放配置信息的
* 本项目使用etcd作为配置中心，通过进程conf业务配置到etcd，以供gateway和server使用
* gateway和server进程，启动时只需要指定etcd和配置信息，即可启动，其不会读取本地的配置信息
* 配置信息什么时候上传？
  * 首次启动时候，必须先通过conf进程上传配置信息`conf/app.yaml`到etcd
  * 配置有变动时候，也可以通过conf上传
> 配置信息上传之后，在gateway和server是即时生效的
> 配置生效不同于配置对应的逻辑生效
> 比如侦听了一个服务端口在8090，此时配置信息修改为8091，虽然配置生效了，但是无法更改配置对应的侦听端口
> 比如一个限制最大次数的值，每次都是从viper读取，若是此配置更改了，此配置对应的逻辑也会生效

### 服务治理
* etcd提供前缀+租赁方式，可以很便利地实现一个服务治理
* 本项目中，server服务注册+租赁提供服务，gateway服务发现使用server提供的服务
* 

### 链路追踪
* 路径：`pkg/trace`
* 使用的是jeager组件，可以关闭，默认是开启
* 链路追踪包含
    * 请求路径
    * 请求耗时
    * 其他组件追踪（函数追踪，db追踪，cache追踪)


### 中间件（拦截器）
* 路径：`pkg/intercepor`
* 使用的组件有两个方面
    * grpcgateway官方提供的常用组件
    * 自定义组件

* 使用的是jeager组件，可以关闭，默认是开启
* 链路追踪包含
    * 请求路径
    * 请求耗时
    * 其他组件追踪（函数追踪，db追踪，cache追踪)

---
### [监控]()
* 监控使用的是[prometheus]()作为数据收集和[grafana]()作为数据展示和管理
* 监控是可选的，通过控制开关来控制是否开启,默认是关闭
* 监控分为三个维度，分别是：
    * [机器监控]()：对机器的指标进行监控，cpu,mem,io等
    * [进程监控]()：通过进程名，对某些进程进行监控，cpu,mem,io等指标
    * [业务监控]()：对业务进行监控，比如业务的qps，耗时等等
#### 基础概念
* 指标源：产生和提供指标的服务进程，一般指各种exporter，各种业务服务器等
* 指标收集服务器：此处是prometheus
* prometheus是通过定时向指标源拉取的方式获取指标，并且保存展示

指标系统构成
    ![指标系统构成](https://www.hualigs.cn/image/60c09a9aa5bfc.jpg)


#### 监控安装
* 启动**prometheus**和**grafana**服务
    * 此处简单演示下怎么启动prometheus服务，是以用一个docker启动起来（需要先准备好docker和docker-compose）,若是已经有了这两个服务，请忽略此步骤，准备好入口即可
    * 假设我们部署服务在机器`129.28.162.42/172.30.0.14`,目录`/data/docker/metrics`(使用者替换成自己的ip以使用之)
    * 最终部署文件目录
        ```sh
        [root@VM_0_14_centos metrics]# tree -L 3
        .
        |-- docker-compose.yml
        |-- grafana
        |   |-- config
        |   |   `-- grafana.ini # grafana的启动文件，此文件使用默认即可，然后修改下登录的账号密码
        |   |-- data
        |   `-- plugins
        `-- prometheus
            |-- config
            |   |-- file_sd #此目录是指标目标，包含机器，进程，业务等提供指标的实例
            |   `-- prometheus.yml # prometheus的启动文件
            `-- data
        ```
    * docker-compose文件
        ```yaml
        # /data/docker/metrics/docker-compose.yml
        version: '2'

        networks:
          monitor:
            driver: bridge

        services:
          prometheus:
            image: prom/prometheus:latest
            container_name: prometheus
            hostname: prometheus
            restart: always
            volumes:
              - /data/docker/metrics/prometheus/config:/etc/prometheus
              - /data/docker/metrics/prometheus/data:/prometheus
            ports:
              - "9091:9091"
            expose:
              - "8086"
            command:
              - '--config.file=/etc/prometheus/prometheus.yml' #docker中的配置文件
              - '--log.level=info'
              - '--web.listen-address=0.0.0.0:9091' #服务接口
              - '--storage.tsdb.path=/prometheus'
              - '--storage.tsdb.retention=15d' #保存15天
              - '--query.max-concurrency=50'
            networks:
              - monitor

          grafana:
            image: grafana/grafana:7.5.3
            container_name: grafana
            restart: always
            volumes:
              - /data/docker/metrics/grafana/config/grafana.ini:/etc/grafana/grafana.ini
            ports:
              - "3000:3000"
              - "25:25"
            networks:
              - monitor
            depends_on:
              - prometheus
        ```
        > 由上面配置我们可以得知
        > * prometheus的地址为：`${out_ip}:9091`
        > * grafana的数据源为:`http://prometheus:9091`
        > * grafana的管理后台为:`${out_ip}:3000`
    * grafana的配置文件
        ```ini
        # grafana/config/grafana.ini
        ...
        [server]
        domain = 129.28.162.42
        root_url = %(protocol)s://%(domain)s:%(http_port)s/
        [security]
        admin_user = admin
        admin_password = a2zone2ten
        ...
        ```
        > 以上可知，grafana的后台管理地址为`129.28.162.42:3000`(129.28.162.42为部署机器的外网ip，根据自己机器变更之)
        > 登录的账号密码为`admin/a2zone2ten`
    * prometheus的配置文件
        ```yaml
        # prometheus/config/prometheus.yml
        #my global config
        global:
          scrape_interval:     15s 
          evaluation_interval: 15s 

        #Alertmanager configuration
        alerting:
          alertmanagers:
          - static_configs:
            - targets:
              #- alertmanager:9093

        #Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
        rule_files:
          #- "first_rules.yml"
          #- "second_rules.yml"

        #A scrape configuration containing exactly one endpoint to scrape:
        #Here it's Prometheus itself.
        scrape_configs:
          #The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
          - job_name: 'prometheus'
            static_configs:
            - targets: ['localhost:9091']

          # iyfiysi scrape config
          - job_name: 'iyfiysi'
            scrape_interval: 5s
            scheme: http
            tls_config:
              insecure_skip_verify: true
            file_sd_configs:
            - files:
                - /etc/prometheus/file_sd/*.yaml
              refresh_interval: 10s
        ```
    * 启动：`docker-compose up`
* 启动[机器监控]()
    * 下面实例中，是运行在linux机器的，[其他环境下载地址](https://github.com/prometheus/node_exporter/releases/)
    * 机器监控，使用的是[node_exporter](https://github.com/prometheus/node_exporter/releases/download/v0.18.1/node_exporter-0.18.1.linux-amd64.tar.gz)
    * 解压后启动启动：`nohup ./node_exporter --web.listen-address=":9100" >/dev/null 2>&1 &`
        > 以上启动命令，可知此指标源的地址是`XXX.XXX.XXX.XXX:9100`
    * 添加指标源配置至prometheus的指标源文件目录
        ```yaml
        # /data/docker/metrics/prometheus/config/file_sd/node.yaml
        # gen by iyfiysi at 2021 Jun 03
        - labels:
            project: "/qq.com/test" #修改为org/project的形式
            role: "node"                         
          targets:
            - "172.30.0.14:9100" #修改为node_exporter侦听地址
        ```
* 启动[进程监控]()
    * 下面实例中，是运行在linux机器的，[其他环境下载地址](https://github.com/ncabatoff/process-exporter/releases)
    * 机器监控，使用的是[process-exporter](https://github.com/ncabatoff/process-exporter/releases/download/v0.7.5/process-exporter-0.7.5.linux-amd64.tar.gz)
    * 启动配置
        ```yaml
        process_names:
            - comm:
              - test_gateway
              - test_server
        ```
    * 解压后启动启动：`nohup ./process-exporter -config.path process.yml -web.listen-address=":9256" >/dev/null 2>&1 & `
        > 以上启动命令，可知此指标源的地址是`XXX.XXX.XXX.XXX:9256`
    * 添加指标源配置至prometheus的指标源文件目录
        ```yaml
        # /data/docker/metrics/prometheus/config/file_sd/process.yaml
        # gen by iyfiysi at 2021 Jun 03
        - labels:
            project: "/qq.com/test" #修改为org/project的形式
            role: "process"                         
          targets:
            - "172.30.0.14:9256" #修改为node_exporter侦听地址
        ```
* 启动[业务监控]()
    * 由于业务服务器启动多少是根据业务而定，是以其势必会有业务服务器乱启动，乱关停
    * 通过服务发现，将当前运行业务服务器找到，并且告知prometheus这些指标源
    * 使用[confd](https://github.com/kelseyhightower/confd/releases/download/v0.16.0/confd-0.16.0-linux-amd64)，通过获取ectd中注册的业务服务器，来生成业务服务器的指标源
    * 启动配置已经由iyfiysi新建项目时候生成
        ```sh
        [root@VM_0_14_centos test]# tree metric/confd/
        metric/confd/
        |-- conf.d
        |   `-- confd.toml
        |-- s.sh
        `-- templates
            `-- confd_tmpl.tmpl
        ```
    * 将生成目录修正为prometheus的指标源文件目录
        ```yaml
        # gen by iyfiysi at 2021 Jun  03
        # confd config file for qq.com/test

        [template]
        src = "confd_tmpl.tmpl" # templates/
        #dest = "${PATH_TO_PROMETHEUS}/config/file_sd/test.yaml" # @TODO 此处修改为prometheus的指标源文件目录
        dest = "/data/docker/metrics/prometheus/config/file_sd/test.yaml"
        keys = [
            "/test/metric",
        ]
        ```
    * 启动服务`sh s.sh`
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
## 高端制造
### 项目生成逻辑
![](https://i0.hdslb.com/bfs/album/e5af3a3404ea5b34b182a6d27aa887750eb73d6c.png)
### 服务生成逻辑
![](https://www.hualigs.cn/image/60c0aeb81bc99.jpg)

---
