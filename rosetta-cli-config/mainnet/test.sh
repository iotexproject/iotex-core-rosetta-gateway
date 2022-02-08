GRN=$'\e[32;1m'
RED=$'\e[31;1m'
OFF=$'\e[0m'

printf "${GRN}### Run rosetta-cli check:data${OFF}\n"
rosetta-cli-0.7.2 check:data --configuration-file ./iotex.json >/dev/null 2>&1 &
dataCheckPID=$!
wait $dataCheckPID
if [ "$?" == 0 ]
then
    printf "${GRN}✓ Run rosetta-cli check:data succeeded${OFF}\n"
else
    printf "${RED}x Run rosetta-cli check:data failed${OFF}\n"
fi
