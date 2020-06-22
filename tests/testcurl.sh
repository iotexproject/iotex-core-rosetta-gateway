#!/usr/bin/env bash
curl -X POST --data '{
    "network_identifier": {
        "blockchain": "IoTex",
        "network": "testnet"},
    "account_identifier": {
        "address": "io1vdtfpzkwpyngzvx7u2mauepnzja7kd5rryp0sg"
    }}' http://127.0.0.1:8080/account/balance

#response:
#{"block_identifier":{"index":3986321,"hash":"931345a809f68dd454716f75c3a08350232be071a56212fac7fb666fc4e608c5"},"balances":[{"value":"12000000000000000000","currency":{"symbol":"IoTex","decimals":18}}],"metadata":{"nonce":0}}


curl -X POST --data '{"metadata": {}}' http://127.0.0.1:8080/network/list
#response:
#{"network_identifiers":[{"blockchain":"IoTex","network":"testnet"}]}

curl -X POST --data '{
    "network_identifier": {
        "blockchain": "IoTex",
        "network": "testnet"},"metadata": {}}' http://127.0.0.1:8080/network/options
#response:
#{"version":{"rosetta_version":"1.3.5","node_version":"v1.0.0"},"allow":{"operation_statuses":[{"status":"OK","successful":true}],"operation_types":["transfer"],"errors":[{"code":1,"message":"unable to get chain ID","retriable":true},{"code":2,"message":"invalid blockchain specified in network identifier","retriable":false},{"code":3,"message":"invalid sub-network identifier","retriable":false},{"code":4,"message":"invalid network specified in network identifier","retriable":false},{"code":5,"message":"network identifier is missing","retriable":false},{"code":6,"message":"unable to get latest block","retriable":true},{"code":7,"message":"unable to get genesis block","retriable":true},{"code":8,"message":"unable to get account","retriable":true},{"code":9,"message":"blocks must be queried by index and not hash","retriable":false},{"code":10,"message":"invalid account address","retriable":false},{"code":11,"message":"a valid subaccount must be specified ('general' or 'escrow')","retriable":false},{"code":12,"message":"unable to get block","retriable":true},{"code":13,"message":"operation not implemented","retriable":false},{"code":14,"message":"unable to get transactions","retriable":true},{"code":15,"message":"unable to submit transaction","retriable":false},{"code":16,"message":"unable to get next nonce","retriable":true},{"code":17,"message":"malformed value","retriable":false},{"code":18,"message":"unable to get node status","retriable":true}]}}

curl -X POST --data '{
    "network_identifier": {
        "blockchain": "IoTex",
        "network": "testnet"},"metadata": {}}' http://127.0.0.1:8080/network/status
#response:
#{"current_block_identifier":{"index":3998161,"hash":"03f2ed3dd20912f6636f43d3eb9bb8423b340ae68185c3959dc1d1f1e83549a5"},"current_block_timestamp":1592618000000,"genesis_block_identifier":{"index":1,"hash":""},"peers":null}

curl -X POST --data '{
    "network_identifier": {
        "blockchain": "IoTex",
        "network": "testnet"
    },
    "options": {"id":"io1vdtfpzkwpyngzvx7u2mauepnzja7kd5rryp0sg"}
}' http://127.0.0.1:8080/construction/metadata
#response:
#{"metadata":{"nonce":{"Nonce":0,"Balance":"12000000000000000000"}}}

curl -X POST --data '{
    "network_identifier": {
        "blockchain": "IoTex",
        "network": "testnet"},
    "signed_transaction": "0a470801100118c0843d220d31303030303030303030303030522e0a01311229696f316c397661716d616e776a3437746c7270763665746633707771307330736e73713476786b6532124104ea8046cf8dc5bc9cda5f2e83e5d2d61932ad7e0e402b4f4cb65b58e9618891f54cba5cfcda873351ad9da1f5a819f54bba9e8343f2edd1ad34dcf7f35de552f31a41d53b8aa4b0165326dcf2eddf4da1fcba8e864f805426c73ee7e73748713c48774bf117f7e78f18459645386ecbb644ca3cca89069920b20ff405768d3d1d6bb301"
}' http://127.0.0.1:8080/construction/submit
#response:
#{"transaction_identifier":{"hash":"292cda920534be56c78d6f13686dc7dbb94b77714b93abefb9f1e18679e2ae27"}}

# transfer action
curl -X POST --data '{
    "network_identifier": {
        "blockchain": "IoTex",
        "network": "testnet"},
    "block_identifier": {"index": 390873}}' http://127.0.0.1:8080/block
#response:
#{"block":{"block_identifier":{"index":390873,"hash":"5c084459315fcf0839ed9f2d8b89ca8fb039695a56007a071e5ce9d3c8908d95"},"parent_block_identifier":{"index":390872,"hash":"3ae76de97535f4908d7dd6b2d5f232543b1e5a9fe80a0e9d8f91fdd27d9363eb"},"timestamp":1573620900000,"transactions":[{"transaction_identifier":{"hash":"b37d5db44bd3dc182617b56744e12cab94486808eae1dc401599b611ed388164"},"operations":[{"operation_identifier":{"index":0},"type":"fee","status":"success","account":{"address":"io1ph0u2psnd7muq5xv9623rmxdsxc4uapxhzpg02"},"amount":{"value":"-10000000000000000","currency":{"symbol":"IoTex","decimals":18}}},{"operation_identifier":{"index":1},"type":"transfer","status":"success","account":{"address":"io1ph0u2psnd7muq5xv9623rmxdsxc4uapxhzpg02"},"amount":{"value":"-10000000000000000000","currency":{"symbol":"IoTex","decimals":18}}},{"operation_identifier":{"index":2},"type":"transfer","status":"success","account":{"address":"io1vdtfpzkwpyngzvx7u2mauepnzja7kd5rryp0sg"},"amount":{"value":"10000000000000000000","currency":{"symbol":"IoTex","decimals":18}}}]}]}}

# Execution multisend
curl -X POST --data '{
    "network_identifier": {
        "blockchain": "IoTex",
        "network": "testnet"},
    "block_identifier": {"index": 4032647}}' http://127.0.0.1:8080/block

# stakeCreate action
curl -X POST --data '{
    "network_identifier": {
        "blockchain": "IoTex",
        "network": "testnet"},
    "block_identifier": {"index": 4034780}}' http://127.0.0.1:8080/block
#response:
#{"block":{"block_identifier":{"index":4034780,"hash":"bc1ad74d423f84e553602798e86019254b70d4499f1738a11c285ab9e31ea3b2"},"parent_block_identifier":{"index":4034779,"hash":"ac25b97cb7c9743b496cf45586d442ae4777753c77ca08d539bb96f30bca08c6"},"timestamp":1592801095000,"transactions":[{"transaction_identifier":{"hash":"9f261c47ad6611388c8e4569d2db378d2a7d98607c4259f5f9819ae6703742e6"},"operations":[{"operation_identifier":{"index":0},"type":"fee","status":"success","account":{"address":"io1mflp9m6hcgm2qcghchsdqj3z3eccrnekx9p0ms"},"amount":{"value":"-10000000000000000","currency":{"symbol":"Iotx","decimals":18}}},{"operation_identifier":{"index":1},"type":"stakeCreate","status":"success","account":{"address":"io1mflp9m6hcgm2qcghchsdqj3z3eccrnekx9p0ms"},"amount":{"value":"-100000000000000000000","currency":{"symbol":"Iotx","decimals":18}}}]}]}}