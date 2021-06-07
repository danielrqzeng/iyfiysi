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
### 3.创建项目
* `cd /data/project`
* `iyfiysi new -n test.com -a test`:项目将会构建在目录`/data/project/test.com/test`中
* 启动etcd
* 编译项目`cd /data/project/test.com/test`;`sh build.sh`
* 运行
    * `./test_server`
    * `./test_gateway`

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

### 监控
监控主要分为三个部分(此处假设已经启动了grafana和prometheus)
#### 机器监控
* 使用node_exporter对机器进行监控
* 此处使用版本为[node_exporter `0.18.1`](https://github.com/prometheus/node_exporter/releases/download/v0.18.1/node_exporter-0.18.1.linux-amd64.tar.gz)
* 配置prometheus中的pull配置如下
    ```yaml
    scrape_configs:
        ...

        - job_name: "node"
          static_configs:
          - targets: ["172.30.0.14:9100"]

        ...
    ```
* 其中的targets为我们将要启动的nod-exportor的服务端口，其默认提供的metric服务为`http://172.30.0.14:9100/metrics`
* 将node_exportor启动即可`nohup ./node_exporter >/dev/null 2>&1 &`
* 在grafana中导入dashbord,导入以下路径即可`https://grafana.com/grafana/dashboards/8919`
    ![效果图](https://i.loli.net/2021/04/14/xP7D9bS16fFopkn.png)
#### 进程监控
* 使用process_exporter对某些进程进行监控
* 此处使用版本为[process_exporter `0.7.5`](https://github.com/ncabatoff/process-exporter/releases/download/v0.7.5/process-exporter-0.7.5.linux-amd64.tar.gz)
* 配置prometheus中的pull配置如下
    ```yaml
    scrape_configs:
        ...

        - job_name: "process"
          static_configs:
          - targets: ["172.30.0.14:9256"]

        ...
    ```
* 其中的targets为我们将要启动的process_exportor的服务端口，其默认提供的metric服务为`http://172.30.0.14:9256/metrics`
* process-exportor本身的配置（即指明要监控那些进程，此处假设要监控的进程名称为voice_gateway和voice_server）
    ```yaml
    # process.yml
    process_names:
    - comm:
      - voice_gateway
      - voice_server    
    ```
* 将process-exportor启动即可`nohup ./process-exporter -config.path process.yml >/dev/null 2>&1 &`
* 在grafana中导入dashbord,导入以下路径即可`https://grafana.com/grafana/dashboards/249`
    ![效果图](https://i.loli.net/2021/04/15/YfN3r4JdV1ceX9o.png)
#### 服务监控
##### go本身监控
* 监控go_goroutines的数量变化：`go_goroutines{job=~"iyfiysi.*"}`
* 在用内存变化:`process_resident_memory_bytes{job=~"iyfiysi.*"}`
```json
go_gc_duration_seconds：持续时间秒
go_gc_duration_seconds_sum：gc-持续时间-秒数-总和
go_memstats_alloc_bytes：Go内存统计分配字节
go_memstats_alloc_bytes_total：Go内存统计分配字节总数
go_memstats_buck_hash_sys_bytes：用于剖析桶散列表的堆空间字节
go_memstats_frees_total：内存释放统计
go_memstats_gc_cpu_fraction：垃圾回收占用服务CPU工作的时间总和
go_memstats_gc_sys_bytes：圾回收标记元信息使用的内存字节
go_memstats_heap_alloc_bytes：服务分配的堆内存字节数
go_memstats_heap_idle_bytes：申请但是未分配的堆内存或者回收了的堆内存（空闲）字节数
go_memstats_heap_inuse_bytes：正在使用的堆内存字节数
go_memstats_heap_objects：堆内存块申请的量
go_memstats_heap_released_bytes：返回给OS的堆内存
go_memstats_heap_sys_bytes：系统分配的作为运行栈的内存
go_memstats_last_gc_time_seconds：垃圾回收器最后一次执行时间
go_memstats_lookups_total：被runtime监视的指针数
go_memstats_mallocs_total：服务malloc的次数
go_memstats_mcache_inuse_bytes：mcache结构体申请的字节数(不会被视为垃圾回收)
go_memstats_mcache_sys_bytes：操作系统申请的堆空间用于mcache的字节数
go_memstats_mspan_inuse_bytes：用于测试用的结构体使用的字节数
go_memstats_mspan_sys_bytes：系统为测试用的结构体分配的字节数
go_memstats_next_gc_bytes：垃圾回收器检视的内存大小
go_memstats_other_sys_bytes：golang系统架构占用的额外空间
go_memstats_stack_inuse_bytes：正在使用的栈字节数
go_memstats_stack_sys_bytes：系统分配的作为运行栈的内存
go_memstats_sys_bytes：服务现在系统使用的内存
go_threads：线程
```
##### 请求
* 分每个server(服务器，不是服务），统计某段时间平均请求数量:`sum(rate(grpc_server_handled_total{job=~"iyfiysi.*"}[$__interval])) by (instance)`

* 对每个方法，统计其请求排行:(X)
`sum(rate(grpc_server_handled_total{job=~"iyfiysi.*"}[$__interval])) by (grpc_method)`

`sum(rate(grpc_server_started_total{job=~"iyfiysi.*"}[1m])) by (grpc_method)`

`topk(10, sort_desc(sum(grpc_server_handled_total) by (grpc_method)))`
* 分每个方法，统计其qps
    * 总qps：`sum(rate(grpc_server_handled_total{job=~"iyfiysi.*"}[$__interval])) by (grpc_method)`
    * 时均错误请求：`sum(rate(grpc_client_handled_total{job=~"iyfiysi.*",grpc_code!="OK"}[$__interval])) by (grpc_method)`
    * 时均成功请求：`sum(rate(grpc_client_handled_total{job=~"iyfiysi.*",grpc_code="OK"}[$__interval])) by (grpc_method)`
* 分每个方法，统计其耗时（tp耗时）
    * histogram_quantile(0.99,sum(rate(grpc_server_handling_seconds_bucket{job=~"iyfiysi.*",grpc_type="unary"}[$__interval])) by (grpc_method,le))
    * histogram_quantile(0.90,sum(rate(grpc_server_handling_seconds_bucket{job=~"iyfiysi.*",grpc_type="unary"}[$__interval])) by (grpc_method,le))
    * histogram_quantile(0.75,sum(rate(grpc_server_handling_seconds_bucket{job=~"iyfiysi.*",grpc_type="unary"}[$__interval])) by (grpc_method,le))
    * histogram_quantile(0.50,sum(rate(grpc_server_handling_seconds_bucket{job=~"iyfiysi.*",grpc_type="unary"}[$__interval])) by (grpc_method,le))
    > lagend: {{grpc_method}}-tp50
* 分每个方法，统计成功率
sum(irate(grpc_client_handled_total{job=~"iyfiysi.*",grpc_type="unary",grpc_code!="OK"}[$__interval])) by (grpc_method)

错误率：
sum(rate(grpc_client_handled_total{job=~"iyfiysi.*",grpc_type="unary",grpc_code!="OK"}[$__interval])) by (grpc_method)
 / 
sum(rate(grpc_client_handled_total{job=~"iyfiysi.*",grpc_type="unary"}[$__interval])) by (grpc_method)
 * 100.0
* 对错误码，统计数量
sum(rate(grpc_client_handled_total{job=~"iyfiysi.*",grpc_type="unary",grpc_code!="OK"}[$__interval])) by (grpc_method,grpc_code)