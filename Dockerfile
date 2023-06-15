FROM golang:1.20-alpine as builder

COPY . /app
WORKDIR /app
ENV CGO_ENABLED=0
RUN go env -w GOPROXY=https://goproxy.cn,direct \
    && go build -o gin-rtsp


FROM ubuntu:18.04

ENV TZ=Asia/Shanghai
ENV LANG=en_US.UTF-8
ENV LOG_LEVEL=info

RUN sed -i 's/archive.ubuntu.com/mirrors.aliyun.com/g' /etc/apt/sources.list \
    && sed -i 's/security.ubuntu.com/mirrors.aliyun.com/g'  /etc/apt/sources.list \
    && apt-get update \
    && ln -snf /usr/share/zoneinfo/${TZ} /etc/localtime && echo ${TZ} > /etc/timezone \
    && apt-get install -y locales tzdata ffmpeg \
    && localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8 \
    && apt-get clean && apt-get autoclean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*


COPY --from=builder /app/gin-rtsp /usr/local/bin/gin-rtsp

EXPOSE 3000

ENTRYPOINT [ "/usr/local/bin/gin-rtsp" ]
