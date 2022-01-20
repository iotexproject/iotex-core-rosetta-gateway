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

ROSETTA_PATH=$(pwd)

function constructionCheckTest() {
  cd $ROSETTA_PATH/rosetta-cli-config
  printf "${GRN}### Run rosetta-cli check:construction${OFF}\n"
  rosetta-cli check:construction --configuration-file testing/iotex-testing.json >rosetta-cli.log 2>&1 &
  constructionCheckPID=$!
  sleep 1

  ## TODO change this to sub process, sleep 1s, may not be right
  #  SEND_TO=$(grep -o "Did you forget to fund? \[\w\+\]" rosetta-cli.log | rev | cut -d ']' -f2 | cut -d '[' -f1 | rev | head -n1)

  SEND_TO=$(grep -o "Please fund the address \[\w\+\]" rosetta-cli.log | rev | cut -d ']' -f2 | cut -d '[' -f1 | rev | awk 'END{print $1}')
  cd $ROSETTA_PATH/tests/inject
  printf "${GRN}### Starting transfer, send to: ${SEND_TO}${OFF}\n"
  ROSETTA_SEND_TO=$SEND_TO go test -test.run TestInjectTransfer10IOTX
  printf "${GRN}### Finished transfer funds${OFF}\n"

  sleep 30
  cd $ROSETTA_PATH/rosetta-cli-config
  ## TODO change this grep to a sub process, fail this grep in x sec should fail the test
  COUNT=$(grep -c "Transactions Confirmed: 1" rosetta-cli.log)
  printf "${GRN}### Finished check transfer, count: ${COUNT}${OFF}\n"
  ps -p $constructionCheckPID >/dev/null
  printf "${GRN}### Run rosetta-cli check:construction succeeded${OFF}\n"
}

function dataCheckTest() {
  cd $ROSETTA_PATH/rosetta-cli-config
  printf "${GRN}### Run rosetta-cli check:data${OFF}\n"
  rosetta-cli check:data --configuration-file testing/iotex-testing.json &
  dataCheckPID=$!

  cd $ROSETTA_PATH/tests/inject
  printf "${GRN}### Inject some actions...${OFF}\n"
  go test

  sleep 10 #wait for the last candidate action
  ps -p $dataCheckPID >/dev/null
  printf "${GRN}### Run rosetta-cli check:data succeeded${OFF}\n"
}

function viewTest(){
  cd $ROSETTA_PATH/rosetta-cli-config
  printf "${GRN}### Run rosetta-cli view:balance and view:block...${OFF}\n"
  rosetta-cli view:balance '{"address":"io1ph0u2psnd7muq5xv9623rmxdsxc4uapxhzpg02"}' --configuration-file testing/iotex-testing.json
  rosetta-cli view:block 10 --configuration-file testing/iotex-testing.json
  printf "${GRN}### Run rosetta-cli view succeeded${OFF}\n"
}

function startServer(){
  cd $ROSETTA_PATH
  printf "${GRN}### Starting the iotex server...${OFF}\n"
  GW="iotex-server -config-path=./tests/config_test.yaml -genesis-path=./tests/genesis_test.yaml -plugin=gateway"
  ${GW} &
  sleep 3

  printf "${GRN}### Starting the Rosetta gateway...${OFF}\n"
  ConfigPath=$ROSETTA_PATH/tests/gateway_config.yaml
  GW="iotex-core-rosetta-gateway"
  ${GW} &
  sleep 3
}


printf "${GRN}### Start testing${OFF}\n"
startServer

constructionCheckTest &
# constructionCheckTestPID=$!

dataCheckTest

viewTest

#wait $constructionCheckTestPID

printf "${GRN}### Tests finished.${OFF}\n"
