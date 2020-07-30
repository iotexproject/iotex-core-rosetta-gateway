FROM golang:1.14-buster
ENV GO111MODULE=on
WORKDIR apps/iotex-core-rosetta-gateway
COPY . .
RUN chmod +x tests/test.sh
CMD ./tests/test.sh
