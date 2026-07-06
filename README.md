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

### Backend Setup

The backend is a Go REST API (using Gin) that reads data from the deployed
contract through [go-ethereum](https://github.com/ethereum/go-ethereum).

```shell
cd backend
cp .env.example .env
```

Set `RPC_URL` (your Anvil endpoint), `CONTRACT_ADDRESS` (from the deploy step
above), `JWT_SECRET` (any long random string, used to sign session tokens),
and the `POSTGRES_*` values in `.env`.

Both Postgres and the backend itself run via Docker Compose (the backend
auto-migrates the database schema on startup):

```shell
docker compose up -d --build
```

Inside the container, the backend reaches Postgres via the Docker network
(`POSTGRES_HOST=postgres`, overridden automatically in `docker-compose.yml`)
and reaches Anvil on your host machine via `host.docker.internal` instead of
`127.0.0.1`.

Alternatively, for faster iteration while developing, run just Postgres in
Docker and the backend directly on your host:

```shell
docker compose up -d postgres
go run .
```

Try it out (with Anvil running and the contract deployed):

```shell
curl http://localhost:8080/health
curl http://localhost:8080/campaigns/count
```

If a contract change is made, regenerate the Go bindings from the Foundry
ABI:

```shell
cd ../contract
forge inspect CrowdFunding abi --json > ../backend/contract/CrowdFunding.abi.json
cd ../backend/contract
abigen --abi=CrowdFunding.abi.json --pkg=contract --type=CrowdFunding --out=crowdfunding.go
```

## Make a Real Test in Local Blockchain

`cast` has two modes: `cast call` reads data (free, no gas, no private key),
and `cast send` writes data (sends a real transaction, needs a private key
to sign it).

- Find the deployed contract address in
  `contract/broadcast/CrowdFunding.s.sol/31337/run-latest.json` under
  `"contractAddress"`.
- Get test accounts from the `anvil` terminal output. It prints 10 funded
  accounts with their private keys. The account whose key you put in
  `PRIVATE_KEY` in `.env` is account `#0` (the campaign owner below). Grab
  account `#1` too, to act as a contributor.
- Read the current campaign count:

  ```shell
  cast call <CONTRACT_ADDRESS> "campaignCount()(uint256)" --rpc-url anvil
  ```

- Create a campaign (goal `1000000000000000000` = 1 ETH, duration `604800`
  = 7 days in seconds):

  ```shell
  cast send <CONTRACT_ADDRESS> \
    "createCampaign(string,string,uint256,uint256)" \
    "Save the Turtles" "Help us protect sea turtles" 1000000000000000000 604800 \
    --rpc-url anvil --private-key <ACCOUNT_0_PRIVATE_KEY>
  ```

- Confirm it was created and read it back. The return type must match the
  `Campaign` struct field order (`owner, title, description, goal, deadline,
  amountRaised, withdrawn`):

  ```shell
  cast call <CONTRACT_ADDRESS> "campaignCount()(uint256)" --rpc-url anvil

  cast call <CONTRACT_ADDRESS> \
    "getCampaign(uint256)((address,string,string,uint256,uint256,uint256,bool))" 0 \
    --rpc-url anvil
  ```

- Contribute from a different account:

  ```shell
  cast send <CONTRACT_ADDRESS> "contribute(uint256)" 0 \
    --value 0.5ether \
    --rpc-url anvil --private-key <ACCOUNT_1_PRIVATE_KEY>
  ```

- Check the contribution landed:

  ```shell
  cast call <CONTRACT_ADDRESS> \
    "getContribution(uint256,address)(uint256)" 0 <ACCOUNT_1_ADDRESS> \
    --rpc-url anvil
  ```

- Check campaign status (`0` = Active, `1` = Successful, `2` = Failed,
  matching the `CampaignStatus` enum order):

  ```shell
  cast call <CONTRACT_ADDRESS> "getCampaignStatus(uint256)(uint8)" 0 --rpc-url anvil
  ```

- Once `amountRaised >= goal`, the owner (account `#0`) can withdraw:

  ```shell
  cast send <CONTRACT_ADDRESS> "withdraw(uint256)" 0 --rpc-url anvil --private-key <ACCOUNT_0_PRIVATE_KEY>
  ```

- To test a refund, a campaign's deadline needs to pass without reaching
  its goal. On Anvil you can fast-forward time yourself:

  ```shell
  cast rpc evm_increaseTime 604801
  cast rpc evm_mine
  ```

  Then a contributor can call `refund(uint256)` the same way `contribute`
  was called above.
