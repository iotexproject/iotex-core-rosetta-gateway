#!/usr/bin/env bash
set -o nounset -o pipefail -o errexit
trap "exit 1" INT
# Kill all dangling processes on exit.
cleanup() {
	printf "${OFF}"
	pkill -P $$ || true
	wait || true
}
trap "cleanup" EXIT

# ANSI escape codes to brighten up the output.
GRN=$'\e[32;1m'
OFF=$'\e[0m'
GW="./run.sh"
printf "${GRN}### Starting the Rosetta gateway...${OFF}\n"
cd ..
${GW} &

sleep 3

printf "${GRN}### Validating Rosetta gateway implementation...${OFF}\n"
cd tests
./rosetta-cli check --lookup-balance-by-block false --end 10 --bootstrap-balances ./bootstrap_balances.json --block-concurrency 4
./rosetta-cli view:account '{"address":"io10t7juxazfteqzjsd6qjk7tkgmngj2tm7n4fvrd"}'
./rosetta-cli view:block 4034780
rm -rf /tmp/rosetta-cli*

# Clean up after a successful run.
rm -rf ./test/rosetta*

printf "${GRN}### Tests finished.${OFF}\n"