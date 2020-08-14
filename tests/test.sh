#!/usr/bin/env bash
set -o nounset -o pipefail -o errexit

# Kill all dangling processes on exit.
cleanup() {
  # cat rosetta-cli.log
  printf "${OFF}"
  pkill -P $$ || true
}
trap "cleanup" EXIT

# ANSI escape codes to brighten up the output.
GRN=$'\e[32;1m'
OFF=$'\e[0m'

ROSETTA_PATH=$(pwd)

cd tests
printf "${GRN}### Starting the iotex server...${OFF}\n"
GW="iotex-server -config-path=config_testnet.yaml -genesis-path=genesis_testnet.yaml -plugin=gateway"
${GW} &
sleep 3

printf "${GRN}### Starting the Rosetta gateway...${OFF}\n"
GW="iotex-core-rosetta-gateway"
${GW} &
sleep 3

cd $ROSETTA_PATH/rosetta-cli-config
printf "${GRN}### Run rosetta-cli check:construction${OFF}\n"
rosetta-cli check:construction --configuration-file testing/iotex-testing.json >rosetta-cli.log 2>&1 &
sleep 1
SEND_TO=$(grep -o "Waiting for funds on \w\+" rosetta-cli.log | rev | cut -d' ' -f 1 | rev)
cd $ROSETTA_PATH/tests/inject
printf "${GRN}### Starting transfer, send to: ${SEND_TO}${OFF}\n"
ROSETTA_SEND_TO=$SEND_TO go test -test.run TestInjectTransfer10IOTX
printf "${GRN}### Finished transfer${OFF}\n"
# check transfer result
cd $ROSETTA_PATH/rosetta-cli-config
COUNT=$(grep -c "\[STATS\] Transactions Confirmed: 1 (Created: 1, In Progress: 0, Stale: 0, Failed: 0) Addresses Created: 2" rosetta-cli.log)
printf "${GRN}### count: ${COUNT}${OFF}\n"
if [ $COUNT -lt 1 ]; then
  printf "${GRN}rosetta-cli check:construction test failed${OFF}\n"
  exit 1
else
  printf "${GRN}### Run rosetta-cli check:construction succeeded${OFF}\n"
fi

cd $ROSETTA_PATH/rosetta-cli-config
printf "${GRN}### Run rosetta-cli check:data${OFF}\n"
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
