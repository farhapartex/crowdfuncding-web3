# Solidity Learning Notes — Crowdfunding Project

These notes explain the Solidity concepts actually used in `contract/src/` on this
project, why each one was needed, and the trade-offs involved.

---

## 1. Storing campaigns in a big array (`Campaign[] public campaigns`)

```solidity
Campaign[] public campaigns;
```

Every campaign — from every user, forever — lives in one single dynamic array on
the `CrowdFunding` contract. A new campaign is just `campaigns.push(...)`, and its
id is its index (`campaignId = campaigns.length` before the push).

**Why a list instead of one contract per campaign, or a mapping keyed by some
external id?**

- **One contract to deploy, once.** If each campaign were its own contract, every
  `createCampaign` would mean deploying new bytecode — vastly more gas (deployment
  is one of the most expensive operations in the EVM) and a new address to track
  per campaign.
- **Sequential, predictable ids.** `campaignId` is just an array index, so no id
  generation logic is needed and ids are guaranteed unique and gapless.
- **Cheap existence checks.** `campaignExists` is just `campaignId >= campaigns.length`.
- **Natural pagination.** `getCampaigns(offset, limit)` slices the array directly.

**Why not a `mapping(uint256 => Campaign)` instead?** A mapping would work almost
identically for lookups, but:
- Solidity mappings have no `.length` — you'd have to track a separate counter anyway.
- Arrays support iteration/pagination (`getCampaigns`) directly; mappings don't.
- Since ids are sequential and dense (no deletions), an array is the natural fit —
  a mapping only pays off when keys are sparse or non-sequential (e.g. keyed by an
  address or a hash).

**The trade-off: `MAX_PAGE_SIZE`.**

```solidity
uint256 public constant MAX_PAGE_SIZE = 50;
```

Looping over campaigns in `getCampaigns` costs gas proportional to how many you
read. For a `view` function called off-chain (like our Go backend does via
`crowdFunding.GetCampaigns(...)` — wait, actually we call `GetCampaign(id)` one at a
time from Go), `view` calls are free (no gas, no transaction) as long as they're
called with `eth_call` rather than sent as a transaction. The page-size cap mostly
protects against a *future* on-chain caller (another contract) blowing its gas
limit reading the whole array. This is a defensive habit worth having even when
today's only caller is our indexer.

**Note on `public` array visibility:** `public` on an array auto-generates a
getter `campaigns(uint256 index) view returns (...)`, but it can only return one
element at a time and — importantly — Solidity's auto-getter *skips* nested
`string`/`array` return values in older versions and returns them individually.
That's part of why we wrote an explicit `getCampaign(id)` that returns the whole
`Campaign` struct in one call, rather than relying on the free auto-getter.

---

## 2. Structs and enums as shared types (`CrowdFundingTypes.sol`)

```solidity
enum CampaignStatus { Active, Successful, Failed }

struct Campaign {
    address owner;
    string title;
    string description;
    uint256 goal;
    uint256 deadline;
    uint256 amountRaised;
    bool withdrawn;
}
```

Pulled into their own file so both `CrowdFunding.sol` and any future/alternate
contract (or a V2) can import the same shape without duplicating it. Small detail,
but it's what makes `contract/test/mocks/CrowdFundingV2.sol` able to inherit
`CrowdFunding` and read `campaigns[i].amountRaised` directly — same struct, same
storage layout.

**Why `enum` for status instead of a `string`?** An enum is stored as a `uint8`
under the hood — 1 byte vs. dozens of bytes for a string. `getCampaignStatus`
computes it on the fly from `amountRaised`/`goal`/`deadline` rather than storing it,
so status is always consistent with the underlying numbers and never needs a
separate "sync" step. (Our Go backend does the same thing — see `campaignStatus()`
in `services/chain_status.go` — computing `Active`/`Successful`/`Failed` from raw
numbers instead of trusting a stored flag.)

---

## 3. Custom errors instead of `require(cond, "string")`

```solidity
error GoalMustBeGreaterThanZero();
...
if (goal == 0) revert GoalMustBeGreaterThanZero();
```

Solidity ≥0.8.4 lets you declare `error Foo();` and `revert Foo();` instead of
`require(goal > 0, "goal must be greater than zero")`. Every error string is
otherwise stored in the deployed bytecode *and* included in the revert
data of every failing call. Custom errors:

- Cost far less gas on revert (a few bytes of selector vs. a full ABI-encoded
  string).
- Still show up in traces/ethers.js `error.reason` — our frontend catches these via
  `err.shortMessage || err.message` (see `usePublishCampaign.js`,
  `useWithdrawFunds.js`).
- Can carry typed arguments if needed (ours don't need any — they're all
  zero-argument signals, e.g. `NotCampaignOwner()`, `GoalNotReached()`).

One real gotcha we hit in this project: if `ethers` doesn't have the exact ABI for
the deployed contract (e.g. pointing at a stale contract address after redeploying
`anvil`), it can't decode which custom error fired and reports
`"execution reverted (unknown custom error)"` — that's exactly what showed up when
the withdraw button tried to call a goal that hadn't been reached yet, and matched
what we saw when the contract address in `.env` was stale.

---

## 4. Events, and why the *indexed* args matter

```solidity
event ContributionMade(uint256 indexed campaignId, address indexed contributor, uint256 amount);
event FundsWithdrawn(uint256 indexed campaignId, address indexed owner, uint256 amount);
event ContributionRefunded(uint256 indexed campaignId, address indexed contributor, uint256 amount);
```

Events are the *only* practical way to build transaction history off-chain
without storing it all on-chain again. Every `contribute`/`withdraw`/`refund` call
emits one, and our Go indexer (`services/indexer_service.go`) polls new blocks and
replays these events with `crowdFunding.FilterContributionMade(...)` /
`FilterFundsWithdrawn(...)` / `FilterContributionRefunded(...)` — abigen-generated
iterators over the event logs — into a single `transactions` Postgres table.

**Why mark `campaignId` and the address `indexed`?** Up to 3 event params can be
`indexed`, which puts them into the log's *topics* rather than its *data*. Topics
are what the EVM lets you filter on cheaply (`FilterOpts` + a topic list) without
scanning every log's full data. We filter by nothing here today (`nil, nil` in
`FilterContributionMade(opts, nil, nil)`, meaning "all campaigns, all
contributors") — but the fact that `campaignId` is indexed is what makes it *possible*
to later filter "just this campaign's contributions" efficiently if the log volume
ever grows large.

**Important distinction we learned the hard way:** events give you `event.Raw.BlockNumber`
and `event.Raw.TxHash`, but *not* wall-clock time. `block.timestamp` isn't part of
the event unless you explicitly emit it (we didn't — the struct itself doesn't
store timestamps for gas reasons, see below). Our indexer has to make a *second*
RPC call, `client.HeaderByNumber(ctx, blockNumber).Time`, to resolve a block number
to a real timestamp — cached per block number (`blockTimestampCache` in
`indexer_service.go`) so we don't refetch the same block's header 3 times when a
contribution and a withdrawal land in the same block.

---

## 5. The proxy pattern (`CrowdFundingProxy.sol`) — upgradeability

This is the most involved piece, so it gets the most space.

### The problem it solves

Once a contract is deployed, its code is immutable — you cannot patch a bug or add
a feature to an already-deployed contract address. Normally that's a *feature*
(trustless: nobody can secretly change the rules). But during development, or for
contracts that genuinely need to evolve, that immutability is a liability: every
bugfix means a brand new address, and every existing campaign, every contributor's
recorded contribution, every integration pointing at the old address — all of it
is stranded on the old, frozen contract.

### How the proxy fixes it: separate storage from logic

```solidity
contract CrowdFundingProxy {
    bytes32 private constant IMPLEMENTATION_SLOT = ...;
    bytes32 private constant ADMIN_SLOT = ...;

    fallback() external payable { _delegate(_getImplementation()); }
    receive() external payable { _delegate(_getImplementation()); }

    function _delegate(address impl) private {
        assembly {
            calldatacopy(0, 0, calldatasize())
            let result := delegatecall(gas(), impl, 0, calldatasize(), 0, 0)
            returndatacopy(0, 0, returndatasize())
            switch result
            case 0 { revert(0, returndatasize()) }
            default { return(0, returndatasize()) }
        }
    }
}
```

Two contracts, two jobs:

- **The proxy** is the address everyone talks to (frontend, backend, users). It
  holds *all the storage* (the `campaigns` array actually lives at the proxy's
  address) but has almost no logic of its own.
- **The implementation** (plain `CrowdFunding.sol`) holds *all the code* but its
  own storage is never used in production — it's just a library of logic the proxy
  borrows.

The magic is `delegatecall`. A normal `call` runs the target's code *in the
target's own storage context*. `delegatecall` runs the target's code but *in the
caller's storage context* — so when `CrowdFunding.createCampaign()` does
`campaigns.push(...)`, and that code is running via `delegatecall` from the proxy,
`campaigns` resolves to *the proxy's* storage slot for `campaigns`, not the
implementation contract's. The implementation is never actually holding any of the
real data — it's borrowed code, running against someone else's storage.

`fallback()`/`receive()` catch *every* call the proxy doesn't explicitly define
(which is all of them — the proxy defines no business functions at all) and
forward it via `_delegate`. That's why calling `proxy.createCampaign(...)`
from ethers.js works exactly as if you'd called it on `CrowdFunding` directly —
the ABI is the same, only the address is the proxy's.

### Upgrading

```solidity
function upgradeTo(address newImplementation) external onlyAdmin {
    _setImplementation(newImplementation);
}
```

To ship a fix, you deploy a *new* implementation contract (new address, new
bytecode) and call `upgradeTo` on the proxy. The proxy address never changes.
Every campaign already stored, every contribution already recorded — all of it is
still there, because it was always in the proxy's storage, never the old
implementation's. `contract/test/mocks/CrowdFundingV2.sol` is a working example
in this repo:

```solidity
contract CrowdFundingV2 is CrowdFunding {
    function totalRaised() external view returns (uint256 total) {
        uint256 count = campaigns.length;
        for (uint256 i = 0; i < count; i++) {
            total += campaigns[i].amountRaised;
        }
    }
}
```

It inherits everything from `CrowdFunding` (same storage layout — critical, see
below) and *adds* a new read-only function. After `upgradeTo(v2Address)`, the
exact same proxy address now also answers `totalRaised()`, with all the old data
intact. `contract/test/CrowdFundingProxy.t.sol` exercises exactly this: deploy V1
behind a proxy, create campaigns, upgrade to V2, and prove the old data survived
and the new function works.

### Storage slot collisions — why `IMPLEMENTATION_SLOT`/`ADMIN_SLOT` look so strange

```solidity
bytes32 private constant IMPLEMENTATION_SLOT =
    bytes32(uint256(keccak256("eip1967.proxy.implementation")) - 1);
```

The proxy needs to remember *its own* two pieces of state — "who is the current
implementation" and "who is allowed to upgrade" — but it must store them somewhere
that can *never* collide with any storage variable the implementation contract
might declare (like `campaigns`, which lives at slot 0 in `CrowdFunding`). If the
proxy just used `address implementation;` as a normal Solidity variable, it would
claim slot 0 — the exact slot `CrowdFunding.campaigns` wants — and the two would
silently overwrite each other.

The fix (this is the [EIP-1967](https://eips.ethereum.org/EIPS/eip-1967) standard):
pick a storage slot by hashing an arbitrary, unique-sounding string
(`"eip1967.proxy.implementation"`) and using that huge, effectively-random 256-bit
number as the slot index instead of 0, 1, 2, .... The `- 1` is a deliberate extra
step recommended by the EIP so that even if someone could find a preimage for the
hash (they can't, feasibly), they still couldn't produce the *exact* slot without
also breaking `keccak256` itself. Reading/writing that slot requires raw
`assembly { sload(slot) }` / `sstore(slot, ...)` because Solidity's normal variable
system has no way to say "put this variable at this exact arbitrary slot" — that's
only expressible in inline assembly.

### `onlyAdmin` and the upgrade key

```solidity
modifier onlyAdmin() {
    if (msg.sender != _getAdmin()) revert NotAdmin();
    _;
}
```

Upgrading is extremely powerful — a new implementation can do *anything* with the
existing storage, including draining it. So `upgradeTo` is gated to one admin
address, set once at deploy time in `DeployProxy.s.sol`:

```solidity
CrowdFundingProxy proxy = new CrowdFundingProxy(address(implementation), deployer);
```

In this project the admin is just the deployer's EOA (fine for local/dev). In a
real production system this would typically be a multisig or a timelock contract,
so that no single private key can unilaterally rewrite the rules of every
campaign's funds.

---

## 6. Reentrancy and checks-effects-interactions

```solidity
function withdraw(uint256 campaignId) external campaignExists(campaignId) {
    Campaign storage campaign = campaigns[campaignId];

    if (msg.sender != campaign.owner) revert NotCampaignOwner();
    if (campaign.amountRaised < campaign.goal) revert GoalNotReached();
    if (campaign.withdrawn) revert FundsAlreadyWithdrawn();

    campaign.withdrawn = true;                       // <-- effect, before the call
    uint256 amountToWithdraw = campaign.amountRaised;

    (bool success,) = payable(campaign.owner).call{value: amountToWithdraw}("");
    if (!success) revert TransferFailed();

    emit FundsWithdrawn(campaignId, campaign.owner, amountToWithdraw);
}
```

`campaign.withdrawn = true` is set *before* the external `.call{value: ...}("")`
that actually sends the ETH. This ordering — checks, then state changes
("effects"), then external calls ("interactions") — is the standard defense
against **reentrancy**: if `campaign.owner` is a contract whose `receive()`
function calls back into `withdraw(campaignId)` again mid-transfer, the second
call sees `campaign.withdrawn == true` already and reverts with
`FundsAlreadyWithdrawn`, rather than being able to drain the campaign twice.

This project has tests that prove it:
`contract/test/attackers/ReentrantWithdrawer.sol` is a contract whose `receive()`
calls `crowdFunding.withdraw(campaignId)` a second time the moment it gets paid.
`CrowdFundingReentrancy.t.sol` asserts the attack reverts the *entire* transaction
(including the first, legitimate withdrawal) with `TransferFailed` — because the
low-level `.call` bubbles the revert up, `success` is `false`, and the whole
outer `withdraw` reverts, leaving the campaign's `withdrawn` flag untouched and the
funds still safely in the contract. `refund` follows the identical pattern
(`contributions[...][msg.sender] = 0` before the `.call`).

We don't use OpenZeppelin's `ReentrancyGuard` modifier here — the manual
checks-effects-interactions ordering is enough for these two functions and avoids
the extra `SLOAD`/`SSTORE` a guard modifier would cost on every call. That's a
deliberate simplicity/gas trade-off that only holds because both `withdraw` and
`refund` are careful to flip their state flag before the transfer — if a third
function like this were added carelessly, it would need the same discipline (or a
guard) to stay safe.

---

## 7. Reducing gas: things this contract already does, and things it could do

Gas is paid by whoever sends the transaction (the campaign creator, a contributor,
etc.), so every unnecessary storage write or byte of calldata is real cost to a
real user. What we've applied and what else is available:

**Already applied in this contract:**

- **Custom errors over require-strings** (section 3) — smaller bytecode, cheaper
  reverts.
- **`indexed` event args, not extra storage** (section 4) — history lives in logs
  (cheap-ish, and outside contract storage entirely) instead of in a second
  on-chain array of past transactions. Storage (`SSTORE`) is one of the most
  expensive EVM operations; logs are far cheaper and are exactly what off-chain
  indexers are for.
- **`external` over `public` for functions nobody calls internally** — `external`
  parameters can be read directly from calldata; `public` functions must copy
  calldata into memory first in case they're called internally. Every function
  in `CrowdFunding.sol` is `external`.
- **`calldata` over `memory` for `string` params** (`createCampaign(string calldata
  title, string calldata description, ...)`) — avoids copying the string into
  memory when the function never needs to modify it.
- **Minimal fields per campaign.** We discussed this directly during planning:
  more struct fields = more storage slots = more gas on every `createCampaign`. The
  `Campaign` struct only carries what's needed to enforce the funding rules
  on-chain (`owner`, `goal`, `deadline`, `amountRaised`, `withdrawn`) plus `title`/
  `description` for display. Everything else the app needs (country, category,
  cover image, draft/publish workflow) lives in Postgres instead — see section 8.
- **No on-chain status field.** `CampaignStatus` is computed on read
  (`getCampaignStatus`) from numbers already being stored, instead of adding a
  `status` field that would need a write (`SSTORE`, ~5,000-20,000 gas) every time
  a campaign transitions.

**Gas-reduction techniques this contract *doesn't* use, but are worth knowing:**

- **Struct packing.** Solidity packs multiple small fields into a single 32-byte
  storage slot if they're declared adjacently and together fit in 32 bytes (e.g.
  a `uint128` + a `uint96` + a `bool`). Our `Campaign` struct doesn't really
  benefit from this — most fields (`address`, two `uint256`s, another `uint256`,
  a `bool`) are either already 32 bytes or would need awkward type-shrinking
  (e.g. `uint64` for a deadline) to pack, which risks overflow bugs for
  marginal savings here.
- **`immutable`/`constant` for values fixed at deploy time.** `MAX_PAGE_SIZE` is
  already `constant` (baked into bytecode, zero storage cost). If this contract
  ever took a configurable fee recipient or similar fixed-at-deploy value, marking
  it `immutable` avoids a storage read (`SLOAD`, ~2,100 gas cold) every time it's used.
- **Unchecked arithmetic blocks** for loops/counters that provably can't overflow
  (`unchecked { i++; }` inside `getCampaigns`'s loop) — Solidity ≥0.8 adds
  overflow checks to every arithmetic op by default; skipping that check where
  you've proven it's safe saves a little gas per iteration. Not done here since
  the loop bound is small (`MAX_PAGE_SIZE = 50`) and clarity mattered more.
- **Batching writes.** If a future feature needed to update many campaigns at
  once, doing it in one transaction is cheaper than N separate transactions,
  because each transaction pays a flat 21,000 gas base cost on top of its actual
  work.

---

## 8. Why campaign metadata is split between the chain and Postgres

Directly related to the "minimal struct" gas point above: this project
deliberately keeps two sources of truth, and neither one is "extra" — they store
different things for different reasons.

- **On-chain (`CrowdFunding.sol`) — financial truth.** `goal`, `deadline`,
  `amountRaised`, `withdrawn`. These *must* be on-chain because they're the numbers
  that determine who can withdraw, who can be refunded, and how much ETH actually
  moves. Trusting a centralized database for this would defeat the entire purpose
  of building this on a blockchain — anyone can independently verify these numbers
  by reading the contract directly.
- **Off-chain (Postgres, via `models.Campaign`) — everything else.** `country`,
  `category`, cover image, draft-vs-published workflow state, owner's Auth0
  identity. None of this needs consensus or trustlessness — it's presentation and
  workflow data, and storing it on-chain would mean paying gas for every campaign
  creator to write a cover image URL into permanent, replicated blockchain
  storage, for data nobody needs a trustless guarantee about.

The backend's `services/campaign_service.go` (`GetMyCampaign`) and
`services/public_campaign_service.go` (`toCampaignResponse`) merge both sources
per request: on-chain numbers are always read live from the contract (never
cached/trusted from the DB), while descriptive fields come from Postgres. This
mirrors a pattern we hit directly during development — a campaign's `/my-campaigns/:id`
page once showed "0 ETH raised" because the *API* wasn't merging on-chain data at
all, not because the contract was wrong. The fix was in the Go layer, not Solidity
— a good reminder that "gas efficiency" and "correctness" are solved in two very
different places in a dApp like this.
