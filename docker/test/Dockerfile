FROM golang:1.14-buster
ENV GO111MODULE=on

ARG ROSETTA_CLI_RELEASE="0.7.1"
RUN mkdir -p /usr/local/build-rosetta-cli && \
    wget --quiet --show-progress --progress=bar:force:noscroll -O /usr/local/build-rosetta-cli/rosetta-cli-${ROSETTA_CLI_RELEASE}.tar.gz https://github.com/coinbase/rosetta-cli/archive/v${ROSETTA_CLI_RELEASE}.tar.gz && \
    cd /usr/local/build-rosetta-cli && \
    tar -xf rosetta-cli-${ROSETTA_CLI_RELEASE}.tar.gz && \
    cd rosetta-cli-${ROSETTA_CLI_RELEASE} && go build && \
    cp rosetta-cli /usr/local/bin/rosetta-cli && \
    cd / && rm -rf /usr/local/build-rosetta-cli

ARG CORE_VERSION="r1.6"
RUN git clone --single-branch --branch $CORE_VERSION https://github.com/iotexproject/iotex-core /usr/local/build-core && \
    cd /usr/local/build-core && \
    make clean build && \
    cp ./bin/server /usr/local/bin/iotex-server  && \
    cd / && rm -rf /usr/local/build-core

WORKDIR apps/iotex-core-rosetta-gateway
COPY . .
ENV ConfigPath=/etc/iotex/gateway_config.yaml
RUN mkdir /etc/iotex && cp tests/gateway_config.yaml /etc/iotex && \
    go build -o /usr/local/bin/iotex-core-rosetta-gateway . && \
    chmod +x tests/test.sh
CMD ./tests/test.sh
