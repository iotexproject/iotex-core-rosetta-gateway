## Construction
- [Workflow](#workflow)
- [Test on testnet](#testnet)

## <a name="workflow"/>Construction Workflow

#### `/construction/derive`

Request
```json
{
	"network_identifier": {
		"blockchain": "IoTeX",
		"network": "testnet"
	},
	"public_key": {
		"hex_bytes": "03a1ace8521ab5b2dcadba9af4c0ca65963c3e3e3a3306031486acdc83ec54e7ad",
		"curve_type": "secp256k1"
	}
}
```

Response
```json
{
	"address": "io1pc5nr6s6047ldtg0whkjuujs9yeerzcexwz6nh"
}
```

#### `/construction/preprocess`

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
			"address": "io1pc5nr6s6047ldtg0whkjuujs9yeerzcexwz6nh"
		},
		"amount": {
			"value": "-5619726348293826415",
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
			"address": "io1p395hdexlwsqphhg8p42n6ljmnzzzd0d3y8kta"
		},
		"amount": {
			"value": "5619726348293826415",
			"currency": {
				"symbol": "IOTX",
				"decimals": 18
			}
		}
	}]
}
```

Response
```json
{
	"options": {
		"amount": "5619726348293826415",
		"decimals": 18,
		"recipient": "io1p395hdexlwsqphhg8p42n6ljmnzzzd0d3y8kta",
		"sender": "io1pc5nr6s6047ldtg0whkjuujs9yeerzcexwz6nh",
		"symbol": "IOTX",
		"type": "NATIVE_TRANSFER"
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
		"amount": "5619726348293826415",
		"decimals": 18,
		"recipient": "io1p395hdexlwsqphhg8p42n6ljmnzzzd0d3y8kta",
		"sender": "io1pc5nr6s6047ldtg0whkjuujs9yeerzcexwz6nh",
		"symbol": "IOTX",
		"type": "NATIVE_TRANSFER"
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
			"address": "io1pc5nr6s6047ldtg0whkjuujs9yeerzcexwz6nh"
		},
		"amount": {
			"value": "-5619726348293826415",
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
			"address": "io1p395hdexlwsqphhg8p42n6ljmnzzzd0d3y8kta"
		},
		"amount": {
			"value": "5619726348293826415",
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
	"unsigned_transaction": "0a56100118904e220d3130303030303030303030303052400a13353631393732363334383239333832363431351229696f3170333935686465786c77737170686867387034326e366c6a6d6e7a7a7a6430643379386b74611229696f317063356e723673363034376c6474673077686b6a75756a7339796565727a636578777a366e68",
	"payloads": [{
		"hex_bytes": "ca671c6d94c90608d5ee6ac8372cb262308285b46f38052663e9ca7773a84480",
		"address": "io1pc5nr6s6047ldtg0whkjuujs9yeerzcexwz6nh",
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
	"transaction": "0a56100118904e220d3130303030303030303030303052400a13353631393732363334383239333832363431351229696f3170333935686465786c77737170686867387034326e366c6a6d6e7a7a7a6430643379386b74611229696f317063356e723673363034376c6474673077686b6a75756a7339796565727a636578777a366e68"
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
			"address": "io1pc5nr6s6047ldtg0whkjuujs9yeerzcexwz6nh"
		},
		"amount": {
			"value": "-5619726348293826415",
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
			"address": "io1p395hdexlwsqphhg8p42n6ljmnzzzd0d3y8kta"
		},
		"amount": {
			"value": "5619726348293826415",
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
	"unsigned_transaction": "0a56100118904e220d3130303030303030303030303052400a13353631393732363334383239333832363431351229696f3170333935686465786c77737170686867387034326e366c6a6d6e7a7a7a6430643379386b74611229696f317063356e723673363034376c6474673077686b6a75756a7339796565727a636578777a366e68",
	"signatures": [{
		"hex_bytes": "2d2d5cb0b6096710a2e186ad4760bfb4f83ecec8d7c21961c9cab966d4517e8a727cf6d6cfb31bcbc9affc6f561221ca9950619a776496be0dacbcd46913979500",
		"signing_payload": {
			"hex_bytes": "ca671c6d94c90608d5ee6ac8372cb262308285b46f38052663e9ca7773a84480",
			"address": "io1pc5nr6s6047ldtg0whkjuujs9yeerzcexwz6nh",
			"signature_type": "ecdsa_recovery"
		},
		"public_key": {
			"hex_bytes": "03a1ace8521ab5b2dcadba9af4c0ca65963c3e3e3a3306031486acdc83ec54e7ad",
			"curve_type": "secp256k1"
		},
		"signature_type": "ecdsa_recovery"
	}]
}
```

Response
```json
{
	"signed_transaction": "0a56100118904e220d3130303030303030303030303052400a13353631393732363334383239333832363431351229696f3170333935686465786c77737170686867387034326e366c6a6d6e7a7a7a6430643379386b7461124104a1ace8521ab5b2dcadba9af4c0ca65963c3e3e3a3306031486acdc83ec54e7ad73cd26c6097157561098da5768bb2f9f052227b882d366fc52b1681b2e7806f11a412d2d5cb0b6096710a2e186ad4760bfb4f83ecec8d7c21961c9cab966d4517e8a727cf6d6cfb31bcbc9affc6f561221ca9950619a776496be0dacbcd46913979500"
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
	"transaction": "0a56100118904e220d3130303030303030303030303052400a13353631393732363334383239333832363431351229696f3170333935686465786c77737170686867387034326e366c6a6d6e7a7a7a6430643379386b7461124104a1ace8521ab5b2dcadba9af4c0ca65963c3e3e3a3306031486acdc83ec54e7ad73cd26c6097157561098da5768bb2f9f052227b882d366fc52b1681b2e7806f11a412d2d5cb0b6096710a2e186ad4760bfb4f83ecec8d7c21961c9cab966d4517e8a727cf6d6cfb31bcbc9affc6f561221ca9950619a776496be0dacbcd46913979500"
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
			"address": "io1pc5nr6s6047ldtg0whkjuujs9yeerzcexwz6nh"
		},
		"amount": {
			"value": "-5619726348293826415",
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
			"address": "io1p395hdexlwsqphhg8p42n6ljmnzzzd0d3y8kta"
		},
		"amount": {
			"value": "5619726348293826415",
			"currency": {
				"symbol": "IOTX",
				"decimals": 18
			}
		}
	}],
	"signers": ["io1pc5nr6s6047ldtg0whkjuujs9yeerzcexwz6nh"],
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
	"signed_transaction": "0a56100118904e220d3130303030303030303030303052400a13353631393732363334383239333832363431351229696f3170333935686465786c77737170686867387034326e366c6a6d6e7a7a7a6430643379386b7461124104a1ace8521ab5b2dcadba9af4c0ca65963c3e3e3a3306031486acdc83ec54e7ad73cd26c6097157561098da5768bb2f9f052227b882d366fc52b1681b2e7806f11a412d2d5cb0b6096710a2e186ad4760bfb4f83ecec8d7c21961c9cab966d4517e8a727cf6d6cfb31bcbc9affc6f561221ca9950619a776496be0dacbcd46913979500"
}
```

Response
```json
{
	"transaction_identifier": {
		"hash": "0c931bf5f2754e1464e6d335be9f498f5618b8b9fd1f038178b4697813628c45"
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
	"signed_transaction": "0a56100118904e220d3130303030303030303030303052400a13353631393732363334383239333832363431351229696f3170333935686465786c77737170686867387034326e366c6a6d6e7a7a7a6430643379386b7461124104a1ace8521ab5b2dcadba9af4c0ca65963c3e3e3a3306031486acdc83ec54e7ad73cd26c6097157561098da5768bb2f9f052227b882d366fc52b1681b2e7806f11a412d2d5cb0b6096710a2e186ad4760bfb4f83ecec8d7c21961c9cab966d4517e8a727cf6d6cfb31bcbc9affc6f561221ca9950619a776496be0dacbcd46913979500"
}
```

Response
```json
{
	"transaction_identifier": {
		"hash": "0c931bf5f2754e1464e6d335be9f498f5618b8b9fd1f038178b4697813628c45"
	}
}
```

## <a name="testnet"/>Test Constructions on IoTeX Testnet 
1. (Optional) Run iotex-core-rosetta-gateway locally 
2. Set `online_url` (and `offline_url` if skipped step 1) to be `https://rosetta.testnet.iotex.one` in [`rosetta-cli-config/testnet/iotex.json`](https://github.com/iotexproject/iotex-core-rosetta-gateway/blob/master/rosetta-cli-config/testnet/iotex.json)
3. Run
``` bash
cd rosetta-cli-config
rosetta-cli check:construction --configuration-file testnet/iotex.json
```
4. Request funds from [IoTeX Faucet](https://faucet.iotex.io/) (Note: every Google account can request only once)
