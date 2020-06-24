FROM golang:1.13.5-stretch
ENV GO111MODULE=on
ENV ConfigPath=/etc/iotex/config.yaml
WORKDIR apps/iotex-core-rosetta-gateway
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN mkdir /etc/iotex
COPY config.yaml /etc/iotex
RUN go build -o /usr/local/bin/iotex-core-rosetta-gateway .
CMD [ "iotex-core-rosetta-gateway"]