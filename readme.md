# 字节后端训练营-红包雨项目

### 功能实现

* [x] 流水线部署（火山引擎, 推送k8s集群）
* [x] 反爬机制（反爬中间件）
* [x] 压力测试（ab, webbench，python多线程并发测试）
* [x] 缓存（redis） 
* [x] 雪花算法生成分布式id
* [x] 随机红包金额生成
* [x] 缓存红包队列，动态红包添加
* [x] 全局异常捕获中间件
* [ ] 性能优化（sql explain、profiler、火焰图）
* [x] 消息队列数据库存储削峰（RocketMQ）
* [x] 滚动日志记录（logrus+lumberjack）
* [ ] 熔断机制（hystrix）
* [ ] 优雅重启
* [x] 流量控制（Sentinel匀速排队）
* [x] 解决跨域问题（跨域中间件）

### 环境配置

#### 本地调试

放到$GOPATH/src/red_envelope下

- 升级npm
    - npm -g install npm@6.14.10
- 安装vue
    - sudo npm install -g vue
    - sudo npm install -g @vue/cli
    - sudo npm install @vue/cli-service -g
- 前端安装
    - cd frontend && sudo npm install
    - 降低eslint版本：sudo npm install --save-dev eslint@5
- 安装docker
    - sudo apt -y install docker.io
    - systemctl start docker
- 安装运行mysql容器
    - 拉取：sudo docker pull mysql
    - 启动：sudo docker container run -d --net=host -v /var/lib/mysql:/var/lib/mysql -p 3306:3306 --name mysql --env MYSQL_ROOT_PASSWORD=root mysql
    - 恢复：cat 1.sql | sudo docker exec -i mysql /usr/bin/mysql -u root --password=root
    - 进入：sudo docker container exec -it mysql bash
    - 备份：sudo docker exec mysql /usr/bin/mysqldump -u root --password=root red_envelope > 1.sql
- 安装运行redis容器
    - 拉取：sudo docker pull redis
    - 启动：sudo docker container run -d --net=host -v /home/redis:/home/redis --name redis redis
    - 进入：sudo docker container exec -it redis bash
- 安装运行rocketmq容器
    - 拉取nameserver：sudo docker pull foxiswho/rocketmq:server-4.5.1
    - 拉取broker：sudo docker pull foxiswho/rocketmq:broker-4.5.1
    - 启动nameserver：sudo docker container run -d -p 9876:9876 --name rmq_server foxiswho/rocketmq:server-4.5.1
    - 启动broker：sudo docker container run -d -p 10911:10911 -p 10909:10909 --name rmq_broker --link rmq_server:namesrv -e "NAMESRV_ADDR=namesrv:9876" -e "JAVA_OPTS=-Duser.home=/opt" -e "JAVA_OPT_EXT=-server -Xms128m -Xmx128m" foxiswho/rocketmq:broker-4.5.1
    - 拉取可视化控制台：sudo docker pull styletang/rocketmq-console-ng
    - 启动控制台：sudo docker container run -d --name rmq_console -p 8180:8080 --link rmq_server:namesrv -e "JAVA_OPTS=-Drocketmq.namesrv.addr=namesrv:9876 -Dcom.rocketmq.sendMessageWithVIPChannel=false" -t styletang/rocketmq-console-ng
    - 浏览器访问：http://localhost:8180
- 从docker构建运行
    - cd backend && sudo docker build -f Dockerfile2 -t my_app .
    - sudo docker run --net=host --name my_app --rm my_app
- 不从docker构建运行
    - 下载安装go：cd ~ && curl -O https://dl.google.com/go/go1.15.2.linux-amd64.tar.gz && tar xvf go1.15.2.linux-amd64.tar.gz && sudo chown -R root:root ./go && sudo mv go /usr/local
    - 设置环境变量：sudo vim ~/.profile 写入：
        - export GOPATH=$HOME/go 
        - export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
    - 刷新：source ~/.profile 
    - 设置开启go mod：go env -w GO111MODULE=auto
    - 设置go代理：go env -w GOPROXY=https://goproxy.cn,direct
    - 初始化mod：go mod init
    - 增加缺失的包，移除没用的包：go mod tidy
    - 编译运行：go build main.go && ./main
- 压力测试工具安装
    - ab:sudo apt update && sudo apt install -y apache2-utils
    - webbench:cd ./test && ./webbench.out
    - python多线程并发测试: cd ./text && python test.py
- 生成profiler分析结果火焰图
    - 设置configure中的UseProfiler=true
    - 下载go-torch：go get -u github.com/google/pprof
    - 下载FlameGraph：mkdir -p $GOPATH/src/github.com/uber/go-torch && cd $GOPATH/src/github.com/uber/go-torch && git clone https://github.com/brendangregg/FlameGraph.git
    - 生成火焰图：go-torch -u http://localhost:6060

#### 部署

- 流水线编译构建： go env -w GO111MODULE=auto && go env -w GOPROXY=https://goproxy.cn,direct && go build main.go
- kubectl安装连接节点：
    - curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl
    - chmod +x ./kubectl
    - sudo mv ./kubectl /usr/local/bin/kubectl
    - kubectl version --client
    - 根据火山引擎的连接信息配置公网连接
    - 查看节点：kubectl get node
    - 查看pod：kubectl get pods
    - 进入pods：kubectl exec -it NAME  [-c  containerName] [-n namespace] -- bash [command]
    - 退出：exit
