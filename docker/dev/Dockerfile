FROM golang:1.14.4-stretch as build
ENV GO111MODULE=on
WORKDIR apps/iotex-core-rosetta-gateway
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /usr/local/bin/iotex-core-rosetta-gateway .


FROM debian:stretch
RUN mkdir /etc/iotex
COPY config.yaml /etc/iotex
ENV ConfigPath=/etc/iotex/config.yaml
COPY --from=build /usr/local/bin/iotex-core-rosetta-gateway /usr/local/bin
CMD [ "iotex-core-rosetta-gateway"]
