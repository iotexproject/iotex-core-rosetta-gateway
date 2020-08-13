## Construction Workflow

#### `/construction/derive`

Request
```json
{
	"network_identifier": {
		"blockchain": "IoTeX",
		"network": "testnet"
	},
	"public_key": {
		"hex_bytes": "036a5affe4a290076da7dfc0bd3b3adc722e3310e369b0224feac8ee2a50727443",
		"curve_type": "secp256k1"
	}
}
```

Response
```json
{
	"address": "io1lacjkkue6k20x6dhldx5ksnaz7qxgup6jecd62"
}
```

#### `/construction/preprocess`

Request
```json
{
   "network_identifier":{
      "blockchain":"IoTeX",
      "network":"testnet"
   },
   "operations":[
      {
         "operation_identifier":{
            "index":0
         },
         "type":"NATIVE_TRANSFER",
         "status":"",
         "account":{
            "address":"io1lacjkkue6k20x6dhldx5ksnaz7qxgup6jecd62"
         },
         "amount":{
            "value":"-1654892007125194610",
            "currency":{
               "symbol":"IOTX",
               "decimals":18
            }
         }
      },
      {
         "operation_identifier":{
            "index":1
         },
         "related_operations":[
            {
               "index":0
            }
         ],
         "type":"NATIVE_TRANSFER",
         "status":"",
         "account":{
            "address":"io14yt2xwk953wps5d264gvqq53lvvg9hfzu4mzsc"
         },
         "amount":{
            "value":"1654892007125194610",
            "currency":{
               "symbol":"IOTX",
               "decimals":18
            }
         }
      }
   ]
}
```

Response
```json
{
	"options": {
		"amount": "1654892007125194610",
		"decimals": 18,
		"recipient": "io14yt2xwk953wps5d264gvqq53lvvg9hfzu4mzsc",
		"sender": "io1lacjkkue6k20x6dhldx5ksnaz7qxgup6jecd62",
		"symbol": "IOTX"
	}
}
```

#### `/construction/metadata`

Request
```json
{
	"network_identifier": {
		"blockchain": "IoTeX",
		"network": "testnet"
	},
	"options": {
		"amount": "1654892007125194610",
		"decimals": 18,
		"recipient": "io14yt2xwk953wps5d264gvqq53lvvg9hfzu4mzsc",
		"sender": "io1lacjkkue6k20x6dhldx5ksnaz7qxgup6jecd62",
		"symbol": "IOTX"
	}
}
```

Response
```json
{
	"metadata": {
		"gasLimit": 10000,
		"gasPrice": 1000000000000,
		"nonce": 1
	},
	"suggested_fee": [{
		"value": "10000000000000000",
		"currency": {
			"symbol": "IOTX",
			"decimals": 18
		}
	}]
}
```

#### `/construction/payload`

Request
```json
{
	"network_identifier": {
		"blockchain": "IoTeX",
		"network": "testnet"
	},
	"operations": [{
		"operation_identifier": {
			"index": 0
		},
		"type": "NATIVE_TRANSFER",
		"status": "",
		"account": {
			"address": "io1lacjkkue6k20x6dhldx5ksnaz7qxgup6jecd62"
		},
		"amount": {
			"value": "-1654892007125194610",
			"currency": {
				"symbol": "IOTX",
				"decimals": 18
			}
		}
	}, {
		"operation_identifier": {
			"index": 1
		},
		"related_operations": [{
			"index": 0
		}],
		"type": "NATIVE_TRANSFER",
		"status": "",
		"account": {
			"address": "io14yt2xwk953wps5d264gvqq53lvvg9hfzu4mzsc"
		},
		"amount": {
			"value": "1654892007125194610",
			"currency": {
				"symbol": "IOTX",
				"decimals": 18
			}
		}
	}],
	"metadata": {
		"gasLimit": 10000,
		"gasPrice": 1000000000000,
		"nonce": 1
	}
}
```

Response
```json
{
	"unsigned_transaction": "0aa901100118904e220d313030303030303030303030305292010a13313635343839323030373132353139343631301229696f313479743278776b39353377707335643236346776717135336c7676673968667a75346d7a73631a507b2253656e64657241646472657373223a22696f316c61636a6b6b7565366b3230783664686c6478356b736e617a377178677570366a6563643632222c225265616c5061796c6f6164223a6e756c6c7d",
	"payloads": [{
		"hex_bytes": "861e40e717c21c7ae80fc9639f26b2726a66cd3a5e224fbbb62c78405627fcec",
		"address": "io1lacjkkue6k20x6dhldx5ksnaz7qxgup6jecd62",
		"signature_type": "ecdsa_recovery"
	}]
}
```

#### `/construction/parse`

Request
```json
{
	"network_identifier": {
		"blockchain": "IoTeX",
		"network": "testnet"
	},
	"signed": false,
	"transaction": "0aa901100118904e220d313030303030303030303030305292010a13313635343839323030373132353139343631301229696f313479743278776b39353377707335643236346776717135336c7676673968667a75346d7a73631a507b2253656e64657241646472657373223a22696f316c61636a6b6b7565366b3230783664686c6478356b736e617a377178677570366a6563643632222c225265616c5061796c6f6164223a6e756c6c7d"
}
```

Response
```json
{
	"operations": [{
		"operation_identifier": {
			"index": 0
		},
		"type": "NATIVE_TRANSFER",
		"status": "",
		"account": {
			"address": "io1lacjkkue6k20x6dhldx5ksnaz7qxgup6jecd62"
		},
		"amount": {
			"value": "-1654892007125194610",
			"currency": {
				"symbol": "IOTX",
				"decimals": 18
			}
		}
	}, {
		"operation_identifier": {
			"index": 1
		},
		"related_operations": [{
			"index": 0
		}],
		"type": "NATIVE_TRANSFER",
		"status": "",
		"account": {
			"address": "io14yt2xwk953wps5d264gvqq53lvvg9hfzu4mzsc"
		},
		"amount": {
			"value": "1654892007125194610",
			"currency": {
				"symbol": "IOTX",
				"decimals": 18
			}
		}
	}],
	"signers": null,
	"metadata": {
		"gasLimit": 10000,
		"gasPrice": 1000000000000,
		"nonce": 1
	}
}
```

#### `/construction/combine`

Request
```json
{
	"network_identifier": {
		"blockchain": "IoTeX",
		"network": "testnet"
	},
	"unsigned_transaction": "0aa901100118904e220d313030303030303030303030305292010a13313635343839323030373132353139343631301229696f313479743278776b39353377707335643236346776717135336c7676673968667a75346d7a73631a507b2253656e64657241646472657373223a22696f316c61636a6b6b7565366b3230783664686c6478356b736e617a377178677570366a6563643632222c225265616c5061796c6f6164223a6e756c6c7d",
	"signatures": [{
		"hex_bytes": "3bf81c991e33188ad0f300cd341637c93fa44e3e6cca7e1467469c078af5f7cf45bd2d782be25dbe1955496334cfa485893dbddf991247df29b76bb72a0b658601",
		"signing_payload": {
			"hex_bytes": "861e40e717c21c7ae80fc9639f26b2726a66cd3a5e224fbbb62c78405627fcec",
			"address": "io1lacjkkue6k20x6dhldx5ksnaz7qxgup6jecd62",
			"signature_type": "ecdsa_recovery"
		},
		"public_key": {
			"hex_bytes": "036a5affe4a290076da7dfc0bd3b3adc722e3310e369b0224feac8ee2a50727443",
			"curve_type": "secp256k1"
		},
		"signature_type": "ecdsa_recovery"
	}]
}
```

Response
```json
{
	"signed_transaction": "0a56100118904e220d3130303030303030303030303052400a13313635343839323030373132353139343631301229696f313479743278776b39353377707335643236346776717135336c7676673968667a75346d7a73631241046a5affe4a290076da7dfc0bd3b3adc722e3310e369b0224feac8ee2a5072744355d424c71120a2c399317ee689d5bafe8ce7776c5e0f877d5bc69ad4369c94a51a413bf81c991e33188ad0f300cd341637c93fa44e3e6cca7e1467469c078af5f7cf45bd2d782be25dbe1955496334cfa485893dbddf991247df29b76bb72a0b658601"
}
```

#### `/construction/parse`

Request
```json
{
	"network_identifier": {
		"blockchain": "IoTeX",
		"network": "testnet"
	},
	"signed": true,
	"transaction": "0a56100118904e220d3130303030303030303030303052400a13313635343839323030373132353139343631301229696f313479743278776b39353377707335643236346776717135336c7676673968667a75346d7a73631241046a5affe4a290076da7dfc0bd3b3adc722e3310e369b0224feac8ee2a5072744355d424c71120a2c399317ee689d5bafe8ce7776c5e0f877d5bc69ad4369c94a51a413bf81c991e33188ad0f300cd341637c93fa44e3e6cca7e1467469c078af5f7cf45bd2d782be25dbe1955496334cfa485893dbddf991247df29b76bb72a0b658601"
}
```

Response
```json
{
	"operations": [{
		"operation_identifier": {
			"index": 0
		},
		"type": "NATIVE_TRANSFER",
		"status": "",
		"account": {
			"address": "io1lacjkkue6k20x6dhldx5ksnaz7qxgup6jecd62"
		},
		"amount": {
			"value": "-1654892007125194610",
			"currency": {
				"symbol": "IOTX",
				"decimals": 18
			}
		}
	}, {
		"operation_identifier": {
			"index": 1
		},
		"related_operations": [{
			"index": 0
		}],
		"type": "NATIVE_TRANSFER",
		"status": "",
		"account": {
			"address": "io14yt2xwk953wps5d264gvqq53lvvg9hfzu4mzsc"
		},
		"amount": {
			"value": "1654892007125194610",
			"currency": {
				"symbol": "IOTX",
				"decimals": 18
			}
		}
	}],
	"signers": ["io1lacjkkue6k20x6dhldx5ksnaz7qxgup6jecd62"],
	"metadata": {
		"gasLimit": 10000,
		"gasPrice": 1000000000000,
		"nonce": 1
	}
}
```

#### `/construction/hash`

Request
```json
{
	"network_identifier": {
		"blockchain": "IoTeX",
		"network": "testnet"
	},
	"signed_transaction": "0a56100118904e220d3130303030303030303030303052400a13313635343839323030373132353139343631301229696f313479743278776b39353377707335643236346776717135336c7676673968667a75346d7a73631241046a5affe4a290076da7dfc0bd3b3adc722e3310e369b0224feac8ee2a5072744355d424c71120a2c399317ee689d5bafe8ce7776c5e0f877d5bc69ad4369c94a51a413bf81c991e33188ad0f300cd341637c93fa44e3e6cca7e1467469c078af5f7cf45bd2d782be25dbe1955496334cfa485893dbddf991247df29b76bb72a0b658601"
}
```

Response
```json
{
	"transaction_identifier": {
		"hash": "6d4cf37497c45f5d06087329ccb7c81b0982a276d0b59dd72a3ebef4ce41f266"
	}
}
```

#### `/construction/submit`

Request
```json
{
	"network_identifier": {
		"blockchain": "IoTeX",
		"network": "testnet"
	},
	"signed_transaction": "0a56100118904e220d3130303030303030303030303052400a13313635343839323030373132353139343631301229696f313479743278776b39353377707335643236346776717135336c7676673968667a75346d7a73631241046a5affe4a290076da7dfc0bd3b3adc722e3310e369b0224feac8ee2a5072744355d424c71120a2c399317ee689d5bafe8ce7776c5e0f877d5bc69ad4369c94a51a413bf81c991e33188ad0f300cd341637c93fa44e3e6cca7e1467469c078af5f7cf45bd2d782be25dbe1955496334cfa485893dbddf991247df29b76bb72a0b658601"
}
```

Response
```json
{
	"transaction_identifier": {
		"hash": "6d4cf37497c45f5d06087329ccb7c81b0982a276d0b59dd72a3ebef4ce41f266"
	}
}
```
