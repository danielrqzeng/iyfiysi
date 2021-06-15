[toc]

# 监控服务快速启动
## 通过docker启动
前提：已有`docker`，`docker-compose`

### 整体服务
* 需要启动prometheus作为数据收集，grafana作为数据展示
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
* 启动的docker-compose文件
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

### prometheus配置
* 由以上docker-compose配置可知，prometheus的布置目录为/data/docker/metrics/prometheus
* 其配置文件添加以下
    ```diff
    - # prometheus/config/prometheus.yml
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

    #Load rules once and periodically evaluate them according to the global     'evaluation_interval'.
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

    +  # iyfiysi scrape config
    +  - job_name: 'iyfiysi'
    +    scrape_interval: 5s
    +    scheme: http
    +    tls_config:
    +      insecure_skip_verify: true
    +    file_sd_configs:
    +    - files:
    +        - /etc/prometheus/file_sd/*.yaml
    +      refresh_interval: 10s
    ```
### grafana
* 由以上docker-compose配置可知，grafana的布置目录为/data/docker/metrics/grafana
* 其配置文件添加以下
    ```diff
    - # grafana/config/grafana.ini
    ...
    [server]
    + domain = 129.28.162.42
    root_url = %(protocol)s://%(domain)s:%(http_port)s/
    [security]
    + admin_user = admin
    + admin_password = a2zone2ten
    ...
    ```
    > 以上可知，grafana的后台管理地址为`http://<out_ip>:3000`
    > 登录的账号密码为`admin/a2zone2ten`
### 启动
* `docker-compose up`
* 配置grafana的数据源为`http://prometheus:9091`