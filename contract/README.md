## Crowd Funding Contract

This folder holds the Solidity smart contract for the Crowd Funding learning
project (see `../feature.md` for the full feature list). Plain Solidity is
used here, no OpenZeppelin, so the basics can be learned from scratch.

### Project Structure (planned)

```
contract/
├── src/
│   └── CrowdFunding.sol      # Main contract: campaigns, contributions, withdraw, refund
├── script/
│   └── CrowdFunding.s.sol    # Script to deploy CrowdFunding.sol
├── test/
│   └── CrowdFunding.t.sol    # Unit tests for CrowdFunding.sol
├── lib/
│   └── forge-std/            # Foundry standard library (test helpers, cheatcodes)
├── foundry.toml               # Foundry project configuration
└── README.md
```

These files do not exist yet, they will be added as we build the contract
step by step. This section will be updated if the structure changes (for
example, splitting each campaign into its own contract).

### What We Will Write

**`src/CrowdFunding.sol`**
- A `Campaign` struct to hold owner, title, description, goal, deadline,
  amount raised, and status.
- Storage to keep track of all campaigns (for example, an array or mapping).
- `createCampaign(...)` to add a new campaign.
- `contribute(campaignId)` to let a user send ETH to a campaign and record
  how much each address contributed.
- `withdraw(campaignId)` for the campaign owner to take out funds once the
  goal is reached.
- `refund(campaignId)` for contributors to get their money back if the goal
  was not reached by the deadline.
- View functions like `getCampaign(campaignId)` and `getCampaigns()` to read
  campaign data.
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

```shell
$ forge script script/CrowdFunding.s.sol:CrowdFundingScript --rpc-url <your_rpc_url> --private-key <your_private_key>
```

(This script does not exist yet, it will be added later.)

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
