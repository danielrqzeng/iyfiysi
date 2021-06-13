# etcd快速启动
## 通过docker启动
前提：已有`docker`，`docker-compose`
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
> 服务端口为http://172.30.0.14:2379，http://127.0.0.1:2379，http://172.30.0.14:4001
> 集群端口为http://172.30.0.14:2380
