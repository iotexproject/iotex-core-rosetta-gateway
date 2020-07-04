FROM golang:1.13.5-stretch
ENV GO111MODULE=on
WORKDIR apps/iotex-core-rosetta-gateway
COPY . .
RUN chmod +x tests/test.sh
CMD ./tests/test.sh