FROM golang:1.13.5-stretch
ENV GO111MODULE=on

RUN apt-get install -y --no-install-recommends make

ARG CORE_VERSION="v1.0.0"
ARG GATEWAY_VERSION="master"
ARG POLL_META_URI="https://storage.googleapis.com/blockchain-archive/poll.mainnet.tar.gz"

RUN git clone --single-branch --branch $CORE_VERSION https://github.com/iotexproject/iotex-core /usr/local/build-core && \
    cd /usr/local/build-core && \
    make clean build && \
    cp ./bin/server /usr/local/bin/iotex-server  && \
    cd / && rm -rf /usr/local/build-core

RUN mkdir -p /etc/iotex/meta && \
    curl -L $POLL_META_URI > /etc/iotex/meta/poll.tar.gz && \
    cd /etc/iotex/meta && \
    tar -xzf poll.tar.gz

RUN git clone --single-branch --branch $GATEWAY_VERSION https://github.com/iotexproject/iotex-core-rosetta-gateway /usr/local/build-gateway && \
    cd /usr/local/build-gateway && \
    make clean build && \
    cp ./iotex-core-rosetta-gateway /usr/local/bin/iotex-core-rosetta-gateway  && \
    cd / && rm -rf /usr/local/build-gateway

VOLUME /data
WORKDIR /data

ENV ConfigPath=/data/etc/iotex-rosetta/config.yaml
CMD iotex-server -config-path=/data/etc/iotex/config_override.yaml -genesis-path=/data/etc/iotex/genesis.yaml -plugin=gateway & iotex-core-rosetta-gateway

EXPOSE 8080/tcp
EXPOSE 14014/tcp
EXPOSE 4689/tcp
