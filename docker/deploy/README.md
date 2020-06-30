## Running IoTeX node and Rosetta gateway in a Docker

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
To run the docker image:
```bash
docker run -v {YOUR_LOCAL_DATA_PATH}:/data -p 8080:8080 -p 4689:4689 -it iotex/iotex-core-rosetta
```
