# jaeger快速启动
## 通过docker启动
前提：已有`docker`，`docker-compose`
* Jaeger是Uber开源的分布式跟踪系统，在本项目中作为链路追踪的组件
* 使用改项目需要先准备好etcd服务，或者按照以下步骤启动一个etcd服务
* 假设在目录`/data/docker/jaeger`中启动一个用于体验测试的jaeger镜像
    > 此版本其数据都是放内存里面，因此只是作为体验测试版本，生产环境要做正式部署才是
* 启动脚本
```sh
# /data/docker/jaeger/setup.sh
docker run \                                              
    -d \
    -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
    --rm \
    --name jaeger \
    -p 6831:6831/udp \
    -p 9411:9411 \
    -p 16686:16686 \
    jaegertracing/all-in-one:latest
```
* 端口说明

    |端口号	|协议	|组件	   |功能|
    |--	|--	|--	   |--|
    |5775	|UDP	|agent	|通过thrift的compact协议接收zipkin.thrift数据|
    |6831	|UDP	|agent	|通过thrift的compact协议接收jaeger.thrift数据|
    |6832	|UDP	|agent	|通过thrift的binary协议接收jaeger.thrift数据|
    |5778	|HTTP	|agent	|服务配置接口|
    |16686	|HTTP	|web	|Jaeger Web UI的端口|
    |9411	|HTTP	|collector	|兼容zipkin的http端点|
    > 本项目中，主要使用
    > * 6831端口做数据上报
    > * 16686端口做webui