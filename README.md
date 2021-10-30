# techtrainingcamp-red_envelope
The fifth group of red envelope grabbing back-end function implementation


### 功能实现

* [ ] 容器部署（docker）
* [ ] 反爬机制 
* [ ] 压力测试（webbench）
* [ ] 缓存（redis） 

### 环境配置

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
- 从docker构建
    - cd backend && sudo docker build -t my_app .
- 运行
    - sudo docker run -p 8080:8080 --name my_app --rm my_app