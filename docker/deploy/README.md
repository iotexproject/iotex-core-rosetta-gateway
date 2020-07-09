## Running IoTeX node and Rosetta gateway in a Docker

This setup is for running IoTeX mainnet node with Rosetta only.

Please put the `etc` in this directory under your local data path:
```bash
cp -rf ./etc {YOUR_LOCAL_DATA_PATH}
```

Port `8080` is rosetta gateway port, and port `4689` is IoTeX node p2p port.

If you are running the node behind a reverse proxy, it is suggested to set the reverse proxy external IP address in `etc/iotex/config_override.yaml` `externalHost` field. If you like to use your own private key instead of randomly assigned one, you can also set it in `etc/iotex/config_override.yaml` `producerPrivKey` field.

To build the docker image:
```bash
docker build . -t iotex-core-rosetta

```

You can also find the built image here: `iotex/iotex-core-rosetta:latest`

To run the docker image:
```bash
docker run -v {YOUR_LOCAL_DATA_PATH}:/data -p 8080:8080 -p 4689:4689 -it iotex/iotex-core-rosetta
```

Once your node sync to tip height, you can check with rosetta-cli with following command:
```bash
rosetta-cli check --lookup-balance-by-block=false --bootstrap-balances=./rosetta-cli-config/bootstrap_balances_mainnet.json --exempt-accounts=./rosetta-cli-config/exempt_accounts_mainnet.json --block-concurrency=100
```
Notice that we have one address in the exempt list, that address is our staking protocol address. Currently, we don't have a method to retrieve the total staking token amount of our voters through API. We are working on adding such a method. For now, we exempt this address from reconciliation.
