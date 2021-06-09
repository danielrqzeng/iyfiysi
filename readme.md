# iyfiysi
## quick start
### 1.安装protoc
#### win平台
* 下载pb编译工具protoc [win版本](https://repo1.maven.org/maven2/com/google/protobuf/protoc/3.9.0/protoc-3.9.0-windows-x86_64.exe)
* 可以移到${GO_PATH}/bin目录下，然后使用goland的控制台来使用命令
#### linux平台
* 下载pb编译工具protoc [linux版本](https://repo1.maven.org/maven2/com/google/protobuf/protoc/3.9.0/protoc-3.9.0-linux-x86_64.exe)
* 改名为protoc,添加到bin路径：`mv protoc-3.9.0-linux-x86_64.exe protoc`,`cp protoc /usr/bin/`
### 2.安装iyfiysi和依赖
* `go get github.com/RQZeng/iyfiysi`
* `go mod download`
* 在linux系统中安装`sh build.sh`
### 项目生成逻辑
![](https://i0.hdslb.com/bfs/album/e5af3a3404ea5b34b182a6d27aa887750eb73d6c.png)
### 服务生成逻辑
![](https://www.hualigs.cn/image/60c0aeb81bc99.jpg)
### 3.创建项目
* 项目名称由组织和项目名称构成，此处这名称非常重要，iyfiysi产生的很多[标识]()都是基于此名称
* `cd /data/project`
* 项目初始化：`iyfiysi new -n test.com -a test`
    * 组织名称为：test.com
    * 项目名称为：test
    * 项目生成于`/data/project/test.com/test`
* 启动etcd
* 编译项目`cd /data/project/test.com/test`;`sh build.sh`
* 运行
    * `./test_server`
    * `./test_gateway`

### 服务架构
![](https://www.hualigs.cn/image/60c0ae1f09bf5.jpg)
---

### 服务治理

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