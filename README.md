# Crowd Funding

## Local Setup

### Contract Setup

```shell
cd contract
forge build
```

Run tests:

```shell
forge test
```

Deploy to a local Anvil node:

1. Run `anvil` in a separate terminal, then copy one of its printed private keys.
2. Copy `.env.example` to `.env` and set `RPC_URL` and `PRIVATE_KEY`:

   ```shell
   cp .env.example .env
   ```

3. Run the deploy script:

   ```shell
   forge script script/CrowdFunding.s.sol:CrowdFundingScript --rpc-url anvil --broadcast
   ```
