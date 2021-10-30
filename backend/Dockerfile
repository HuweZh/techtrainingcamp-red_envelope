FROM golang

# 改变工作目录，即执行go命令的目录
WORKDIR $GOPATH/src/red_envelope

# 将本地内容添加到镜像指定目录
ADD . $GOPATH/src/red_envelope

# 设置开启go mod
RUN go env -w GO111MODULE=auto

# 设置go代理
RUN go env -w GOPROXY=https://goproxy.cn,direct

#会在当前目录生成一个go.mod文件用于包管理
RUN go mod init
#增加缺失的包，移除没用的包
RUN go mod tidy

#进入src编译
WORKDIR $GOPATH/src/red_envelope/src
RUN go build -i -o ../bin/red_envelope

# 指定镜像内部服务监听的端口
EXPOSE 8080


# 镜像默认入口命令，即go编译后的可执行文件
ENTRYPOINT ["../bin/red_envelope"]

#CMD ["go","run","app.go"]
