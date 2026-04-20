# --------- 阶段 1: 构建阶段 (builder) ---------
FROM golang:1.24 AS builder

LABEL maintainer="Tencent BlueKing"

# 安装编译铜锁 (Tongsuo) 所需的基础工具 (Debian 环境)
RUN apt-get update && apt-get install -y gcc make perl wget git

# 配置私有仓库访问
ARG GIT_TOKEN
ENV GOPRIVATE=code.cwoa.net
RUN git config --global url."https://oauth2:${GIT_TOKEN}@code.cwoa.net/".insteadOf "https://code.cwoa.net/"

# 下载并编译铜锁
ENV TONGSUO_PATH=/usr/local/tongsuo
RUN wget --no-check-certificate https://github.com/Tongsuo-Project/Tongsuo/archive/refs/tags/8.3.2.tar.gz && \
    tar zxvf 8.3.2.tar.gz && \
    cd Tongsuo-8.3.2 && \
    ./config --prefix=$TONGSUO_PATH -static -fPIC && \
    make -j$(nproc) && \
    make install

# 设置 CGO 编译参数，启用 CGO 并链接铜锁
# 参考 Makefile 中的参数设置
ENV CGO_ENABLED=1 \
    CGO_CFLAGS="-I${TONGSUO_PATH}/include -Wno-deprecated-declarations" \
    CGO_LDFLAGS="-L${TONGSUO_PATH}/lib -lssl -lcrypto -ldl -lpthread" \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /build

# 拷贝源码。构建上下文应为仓库根目录
COPY bcs-services/bcs-user-manager/ .

# 编译参数（由外部传入，通过 -X 注入版本信息等）
ARG LDFLAGS=""

# 执行 go mod tidy 整理依赖并编译
RUN go mod tidy && \
    go build -ldflags "${LDFLAGS}" -o bcs-user-manager main.go

# --------- 阶段 2: 运行阶段 (final) ---------
FROM tencentos/tencentos4-minimal

RUN yum update -y && \
    yum install -y ca-certificates bash vim curl gettext glibc tzdata && \
    yum clean all && \
    rm -rf /var/cache/yum/*

RUN mkdir -p /data/bcs/logs/bcs /data/bcs/cert
RUN mkdir -p /data/bcs/bcs-user-manager

# 从 builder 阶段拷贝编译好的二进制文件
COPY --from=builder /build/bcs-user-manager /data/bcs/bcs-user-manager/

# 拷贝配置文件模板和启动脚本
COPY install/conf/bcs-services/bcs-user-manager/bcs-user-manager.json.template /data/bcs/bcs-user-manager/
COPY install/conf/bcs-services/bcs-user-manager/container-start.sh /data/bcs/bcs-user-manager/

RUN chmod +x /data/bcs/bcs-user-manager/container-start.sh

ENV TZ="Asia/Shanghai"
RUN ln -fs /usr/share/zoneinfo/${TZ} /etc/localtime && echo ${TZ} > /etc/timezone

WORKDIR /data/bcs/bcs-user-manager/
CMD ["/data/bcs/bcs-user-manager/container-start.sh"]