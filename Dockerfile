FROM ubuntu:18.04 as builder

ENV APP_DIR=/go/src
ENV GO_VERSION=1.13.3
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io

COPY . ${APP_DIR}

RUN sed -i 's/archive.ubuntu.com/mirrors.aliyun.com/g' /etc/apt/sources.list \
    && sed -i 's/security.ubuntu.com/mirrors.aliyun.com/g'  /etc/apt/sources.list \
    && apt-get update && apt-get install -y wget make g++ \
    && wget -nv https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz &&  tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz

WORKDIR ${APP_DIR}
RUN  export PATH=$PATH:/usr/local/go/bin && make all


FROM ubuntu:18.04

ENV TZ=Asia/Shanghai
ENV LANG=en_US.UTF-8

RUN sed -i 's/archive.ubuntu.com/mirrors.aliyun.com/g' /etc/apt/sources.list \
    && sed -i 's/security.ubuntu.com/mirrors.aliyun.com/g'  /etc/apt/sources.list \
    && apt-get update \
    && ln -snf /usr/share/zoneinfo/${TZ} /etc/localtime && echo ${TZ} > /etc/timezone \
    && apt-get install -y locales tzdata ffmpeg \
    && localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8 \
    && apt-get clean && apt-get autoclean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ENV APP_DIR=/usr/local/hkapi
ENV GIN_MODE=release
ENV LOG_LEVEL=info
ENV NVR_SAVE_PATH=/srv/hkapi/stream
ENV SQLITE_PATH=/var/hkapi/database
ENV DB_TYPE=sqlite3
ENV DB_DSN=${SQLITE_PATH}/sqlite.db

COPY --from=builder /go/src/lib ${APP_DIR}/lib
COPY --from=builder /go/src/hkapi.out ${APP_DIR}
COPY --from=builder /go/src/run.sh ${APP_DIR}

EXPOSE 3000

WORKDIR ${APP_DIR}
RUN mkdir -p ${NVR_SAVE_PATH} && mkdir -p ${SQLITE_PATH} && chmod +x run.sh

ENTRYPOINT [ "./run.sh" ]