### Testing

The test spins up a standalone IoTeX chain node and a iotex-core-rosetta-gateway first, then start checking the chain with rosetta-cli.
Meanwhile, it starts injecting different actions to the chain including:


1. 20 token transfer actions
2. 1 contract deployment
3. 20 multi-send contract executions.
4. Staking related actions: candidate register, stake bucket creation, stake bucket add deposit, bucket unstake, stake bucket withdraw

The whole test will take around 6 to 10 mins to finish.

Notices that this test exempts two protocol accounts, that is because staking and rewarding protocl addresses are not accessiable in standalone mode.
