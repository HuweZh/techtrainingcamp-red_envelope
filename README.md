# 字节后端训练营-红包雨项目

### 功能实现

* [ ] 流水线部署（火山引擎，github，docker）
* [ ] 反爬机制 
* [ ] 压力测试（webbench）
* [ ] 缓存（redis） 

### 环境配置

#### 本地调试

- 安装vue
    - sudo npm install -g vue
    - sudo npm install -g @vue/cli
    - sudo npm install @vue/cli-service -g
- 前端安装
    - cd frontend && sudo npm install
    - 降低eslint版本：sudo npm install --save-dev eslint@5
- 本地运行
    - 设置开启go mod：go env -w GO111MODULE=auto
    - 设置go代理：go env -w GOPROXY=https://goproxy.cn,direct
    - 初始化mod：go mod init
    - 增加缺失的包，移除没用的包：go mod tidy
    - 编译运行：go build main.go && ./main
