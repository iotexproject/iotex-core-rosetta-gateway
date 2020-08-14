### Testing

The test first spins up a standalone IoTeX chain node and an iotex-core-rosetta-gateway locally, and then it starts testing with rosetta-cli commands: `check:construction/data` and `view:block/account`.

For `check:construction`, it guarantees at least one confirmed transaction is made.

For `check:data`, it injects following different actions onto the chain to ensure balance changes and transaction calculations are correct:

1. 20 token transfer actions
2. 1 contract deployment
3. 20 multi-send contract executions.
4. Staking related actions: candidate register, stake bucket creation, stake bucket add deposit, bucket unstake, stake bucket withdraw

The whole test will take around 6 to 10 mins to finish.

Notices that this test exempts two protocol accounts, that is because staking and rewarding protocl addresses are not accessiable in IoTeX standalone mode node.
