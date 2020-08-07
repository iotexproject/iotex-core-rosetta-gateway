#!/usr/bin/env bash
set -o nounset -o pipefail -o errexit

# Kill all dangling processes on exit.
cleanup() {
	printf "${OFF}"
	pkill -P $$ || true
}
trap "cleanup" EXIT

# ANSI escape codes to brighten up the output.
GRN=$'\e[32;1m'
OFF=$'\e[0m'


cd tests
printf "${GRN}### Starting the iotex server...${OFF}\n"
GW="iotex-server -config-path=config_testnet.yaml -genesis-path=genesis_testnet.yaml -plugin=gateway"
${GW} &
sleep 3

printf "${GRN}### Starting the Rosetta gateway...${OFF}\n"
GW="iotex-core-rosetta-gateway"
${GW} &
sleep 3

cd ../rosetta-cli-config
printf "${GRN}### Run rosetta-cli check...${OFF}\n"
rosetta-cli check:data --configuration-file testing/iotex-testing.json &

cd ../tests/inject
printf "${GRN}### Inject some actions...${OFF}\n"
go test

sleep 10 #wait for the last candidate action

cd ../../rosetta-cli-config
printf "${GRN}### Run rosetta-cli view:account and view:block...${OFF}\n"
rosetta-cli view:account '{"address":"io1ph0u2psnd7muq5xv9623rmxdsxc4uapxhzpg02"}' --configuration-file testing/iotex-testing.json
rosetta-cli view:block 10 --configuration-file testing/iotex-testing.json

printf "${GRN}### Tests finished.${OFF}\n"
