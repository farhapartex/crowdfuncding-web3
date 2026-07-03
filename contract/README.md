## Crowd Funding Contract

This folder holds the Solidity smart contract for the Crowd Funding learning
project (see `../feature.md` for the full feature list). Plain Solidity is
used here, no OpenZeppelin, so the basics can be learned from scratch.

### Project Structure

```
contract/
├── src/
│   ├── CrowdFunding.sol       # Main contract: campaigns, contributions, withdraw, refund
│   └── CrowdFundingTypes.sol  # Shared Campaign struct and CampaignStatus enum
├── script/
│   └── CrowdFunding.s.sol     # Script to deploy CrowdFunding.sol
├── test/
│   └── CrowdFunding.t.sol     # Unit tests for CrowdFunding.sol
├── lib/
│   └── forge-std/             # Foundry standard library (test helpers, cheatcodes)
├── foundry.toml                # Foundry project configuration
└── README.md
```

This section will be updated if the structure changes (for example, splitting
each campaign into its own contract).

### What Each File Has

**`src/CrowdFundingTypes.sol`**
- The `Campaign` struct (owner, title, description, goal, deadline, amount
  raised, withdrawn) and the `CampaignStatus` enum (Active, Successful,
  Failed), shared by the contract and its tests.

**`src/CrowdFunding.sol`**
- Storage to keep track of all campaigns and how much each address
  contributed.
- `createCampaign(...)` to add a new campaign.
- `contribute(campaignId)` to let a user send ETH to a campaign and record
  how much each address contributed.
- `withdraw(campaignId)` for the campaign owner to take out funds once the
  goal is reached.
- `refund(campaignId)` for contributors to get their money back if the goal
  was not reached by the deadline.
- View functions like `getCampaign(campaignId)`, `getCampaigns()`, and
  `getCampaignStatus(campaignId)` to read campaign data.
- Events such as `CampaignCreated` and `ContributionMade` so the backend can
  listen for changes.

**`script/CrowdFunding.s.sol`**
- A Foundry script that deploys `CrowdFunding.sol` to a local Anvil node or a
  testnet.

**`test/CrowdFunding.t.sol`**
- Unit tests covering creating a campaign, contributing, withdrawing after
  success, refunding after failure, and checking the view functions return
  correct data.

## Foundry

**Foundry is a blazing fast, portable and modular toolkit for Ethereum application development written in Rust.**

Foundry consists of:

- **Forge**: Ethereum testing framework (like Truffle, Hardhat and DappTools).
- **Cast**: Swiss army knife for interacting with EVM smart contracts, sending transactions and getting chain data.
- **Anvil**: Local Ethereum node, akin to Ganache, Hardhat Network.
- **Chisel**: Fast, utilitarian, and verbose solidity REPL.

## Documentation

https://book.getfoundry.sh/

## Usage

### Build

```shell
$ forge build
```

### Test

```shell
$ forge test
```

### Format

```shell
$ forge fmt
```

### Gas Snapshots

```shell
$ forge snapshot
```

### Anvil

```shell
$ anvil
```

### Deploy

Copy `.env.example` to `.env` and fill in `RPC_URL` and `PRIVATE_KEY`, then run:

```shell
$ forge script script/CrowdFunding.s.sol:CrowdFundingScript --rpc-url anvil --broadcast
```

`anvil` here is an alias defined in `foundry.toml` under `[rpc_endpoints]`, pointing
at the `RPC_URL` from `.env`. The private key is read inside the script itself
via `vm.envUint("PRIVATE_KEY")`, so it never needs to be typed on the command
line.

### Cast

```shell
$ cast <subcommand>
```

### Help

```shell
$ forge --help
$ anvil --help
$ cast --help
```
