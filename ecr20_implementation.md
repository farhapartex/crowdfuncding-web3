# ERC20 Contribution Support — Implementation Plan

This is a design doc, not yet implemented. It covers everything needed to let a
campaign accept **ETH only**, **a single ERC20 token only**, or **Both**,
across the contract, backend, and frontend — plus a mock ERC20 token so we can
actually test this locally (Anvil has no real USDC/DAI deployed).

## 1. Core design decisions

**One accepted ERC20 token per campaign, chosen at creation.** No mixing
multiple different ERC20s in one campaign, and no price oracle. This keeps the
contract simple and gas-cheap, consistent with how the rest of this contract
has been kept minimal.

**"Both" mode tracks ETH and the token as two independent, parallel pools —
never converted or summed together.** A campaign in "Both" mode has *two*
goals (`goalEth`, `goalToken`) and *two* raised totals
(`amountRaisedEth`, `amountRaisedToken`). There is no unified "total raised"
number, because 1 ETH and 1 USDC aren't comparable without a price feed, and
adding a price oracle (e.g. Chainlink) is a large additional subsystem with its
own trust/failure assumptions — out of scope for what's being asked here.

**"Goal reached" = at least one of the enabled currencies hit its own goal.**
For `EthOnly`/`TokenOnly` this is unambiguous (only one goal exists). For
`Both`, reaching *either* goal unlocks withdrawal of *everything* raised so
far, in both currencies. This mirrors how a real campaign would work — hitting
the target is the trigger, regardless of which currency got you there.

**Refunds are always paid back in the same currency the contributor sent.**
No conversion, ever. Simplest possible logic, and it's the only thing that's
actually fair without a price oracle.

**⚠️ Storage-layout warning (ties back to the proxy pattern):** this changes
the shape of the `Campaign` struct — splitting `goal` into `goalEth`/
`goalToken` and `amountRaised` into `amountRaisedEth`/`amountRaisedToken`, plus
adding `token`/`currencyMode`. That is **not** a safe `upgradeTo()` change
under the existing proxy (see `learning/solidity/note.md` §5) — reusing/
reinterpreting existing storage slots corrupts old data. Two real options:

- **(a) Fresh deploy** (what we've done every time we've reset Anvil this
  project anyway) — simplest, loses existing on-chain campaign data, fine for
  a dev/learning environment.
- **(b) Genuinely additive V2** — keep the old `goal`/`amountRaised` fields
  exactly as-is (they become "the ETH-only legacy fields"), and *append* new
  fields (`token`, `goalToken`, `amountRaisedToken`, `currencyMode`) at the end
  of the struct. This is upgrade-compatible but means every function has to
  keep supporting the old two-field shape forever alongside the new one.

Given this project resets Anvil constantly during development anyway, **(a)
fresh deploy** is the pragmatic choice unless real production data already
exists that must be preserved.

---

## 2. Solidity changes

### 2.1 New dependency: OpenZeppelin Contracts

Currently only `forge-std` is installed. Add OpenZeppelin for `IERC20`,
`SafeERC20` (handles non-standard tokens like USDT that don't return `bool`
correctly), and `ERC20` (base for the mock token):

```bash
forge install OpenZeppelin/openzeppelin-contracts
```

### 2.2 `CrowdFundingTypes.sol`

```solidity
enum CurrencyMode { EthOnly, TokenOnly, Both }

struct Campaign {
    address owner;
    string title;
    string description;
    CurrencyMode currencyMode;
    address token;              // address(0) for EthOnly
    uint256 goalEth;             // 0 if TokenOnly
    uint256 goalToken;           // 0 if EthOnly
    uint256 deadline;
    uint256 amountRaisedEth;
    uint256 amountRaisedToken;
    bool withdrawn;
}
```

### 2.3 `CrowdFunding.sol`

**New imports:**

```solidity
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

using SafeERC20 for IERC20;
```

**New errors:**

```solidity
error EthNotAccepted();
error TokenNotAccepted();
error TokenAddressRequired();
```

**New per-token contribution tracking** (alongside the existing ETH mapping):

```solidity
mapping(uint256 => mapping(address => uint256)) public contributionsEth;   // renamed from `contributions`
mapping(uint256 => mapping(address => uint256)) public contributionsToken;
```

**`createCampaign`** — gains `currencyMode`, `token`, and splits `goal` into
`goalEth`/`goalToken`:

```solidity
function createCampaign(
    string calldata title,
    string calldata description,
    CurrencyMode currencyMode,
    address token,
    uint256 goalEth,
    uint256 goalToken,
    uint256 durationInSeconds
) external returns (uint256 campaignId) {
    if (durationInSeconds == 0) revert DurationMustBeGreaterThanZero();

    if (currencyMode == CurrencyMode.EthOnly) {
        if (goalEth == 0) revert GoalMustBeGreaterThanZero();
    } else if (currencyMode == CurrencyMode.TokenOnly) {
        if (goalToken == 0) revert GoalMustBeGreaterThanZero();
        if (token == address(0)) revert TokenAddressRequired();
    } else {
        if (goalEth == 0 || goalToken == 0) revert GoalMustBeGreaterThanZero();
        if (token == address(0)) revert TokenAddressRequired();
    }

    campaignId = campaigns.length;
    campaigns.push(Campaign({
        owner: msg.sender,
        title: title,
        description: description,
        currencyMode: currencyMode,
        token: token,
        goalEth: goalEth,
        goalToken: goalToken,
        deadline: block.timestamp + durationInSeconds,
        amountRaisedEth: 0,
        amountRaisedToken: 0,
        withdrawn: false
    }));

    emit CampaignCreated(campaignId, msg.sender, currencyMode, token, goalEth, goalToken, campaigns[campaignId].deadline);
}
```

**`contribute` splits into two functions** (keeping one `payable` function that
silently ignores/rejects a token contribution is a footgun — better to make
the two paths structurally distinct so a contributor can't accidentally send
ETH into a token-only campaign and have it get stuck):

```solidity
function contributeEth(uint256 campaignId) external payable campaignExists(campaignId) {
    Campaign storage campaign = campaigns[campaignId];
    if (campaign.currencyMode == CurrencyMode.TokenOnly) revert EthNotAccepted();
    // forge-lint: disable-next-line(block-timestamp)
    if (block.timestamp >= campaign.deadline) revert CampaignHasEnded();
    if (msg.value == 0) revert ContributionMustBeGreaterThanZero();

    campaign.amountRaisedEth += msg.value;
    contributionsEth[campaignId][msg.sender] += msg.value;

    emit ContributionMade(campaignId, msg.sender, address(0), msg.value);
}

function contributeToken(uint256 campaignId, uint256 amount) external campaignExists(campaignId) {
    Campaign storage campaign = campaigns[campaignId];
    if (campaign.currencyMode == CurrencyMode.EthOnly) revert TokenNotAccepted();
    // forge-lint: disable-next-line(block-timestamp)
    if (block.timestamp >= campaign.deadline) revert CampaignHasEnded();
    if (amount == 0) revert ContributionMustBeGreaterThanZero();

    IERC20(campaign.token).safeTransferFrom(msg.sender, address(this), amount);

    campaign.amountRaisedToken += amount;
    contributionsToken[campaignId][msg.sender] += amount;

    emit ContributionMade(campaignId, msg.sender, campaign.token, amount);
}
```

Note `ContributionMade` now carries a `token` address (`address(0)` = ETH) so
one event shape covers both currencies.

**`withdraw`** — goal check becomes an OR across whichever currencies are
active; pays out whatever balance exists in each:

```solidity
function withdraw(uint256 campaignId) external campaignExists(campaignId) {
    Campaign storage campaign = campaigns[campaignId];

    if (msg.sender != campaign.owner) revert NotCampaignOwner();
    if (campaign.withdrawn) revert FundsAlreadyWithdrawn();

    bool goalReached =
        (campaign.goalEth > 0 && campaign.amountRaisedEth >= campaign.goalEth) ||
        (campaign.goalToken > 0 && campaign.amountRaisedToken >= campaign.goalToken);
    if (!goalReached) revert GoalNotReached();

    campaign.withdrawn = true;

    uint256 ethAmount = campaign.amountRaisedEth;
    uint256 tokenAmount = campaign.amountRaisedToken;

    if (ethAmount > 0) {
        (bool success,) = payable(campaign.owner).call{value: ethAmount}("");
        if (!success) revert TransferFailed();
    }
    if (tokenAmount > 0) {
        IERC20(campaign.token).safeTransfer(campaign.owner, tokenAmount);
    }

    emit FundsWithdrawn(campaignId, campaign.owner, ethAmount, tokenAmount);
}
```

**`refund`** — same OR-based "did it succeed" check (inverted), pays back
whichever currencies this specific contributor put in:

```solidity
function refund(uint256 campaignId) external campaignExists(campaignId) {
    Campaign storage campaign = campaigns[campaignId];

    // forge-lint: disable-next-line(block-timestamp)
    if (block.timestamp < campaign.deadline) revert CampaignStillActive();

    bool goalReached =
        (campaign.goalEth > 0 && campaign.amountRaisedEth >= campaign.goalEth) ||
        (campaign.goalToken > 0 && campaign.amountRaisedToken >= campaign.goalToken);
    if (goalReached) revert GoalAlreadyReached();

    uint256 ethContributed = contributionsEth[campaignId][msg.sender];
    uint256 tokenContributed = contributionsToken[campaignId][msg.sender];
    if (ethContributed == 0 && tokenContributed == 0) revert NoContributionToRefund();

    if (ethContributed > 0) {
        contributionsEth[campaignId][msg.sender] = 0;
        (bool success,) = payable(msg.sender).call{value: ethContributed}("");
        if (!success) revert TransferFailed();
    }
    if (tokenContributed > 0) {
        contributionsToken[campaignId][msg.sender] = 0;
        IERC20(campaign.token).safeTransfer(msg.sender, tokenContributed);
    }

    emit ContributionRefunded(campaignId, msg.sender, ethContributed, tokenContributed);
}
```

**Updated event signatures:**

```solidity
event CampaignCreated(uint256 indexed campaignId, address indexed owner, CurrencyMode currencyMode, address token, uint256 goalEth, uint256 goalToken, uint256 deadline);
event ContributionMade(uint256 indexed campaignId, address indexed contributor, address token, uint256 amount);
event FundsWithdrawn(uint256 indexed campaignId, address indexed owner, uint256 ethAmount, uint256 tokenAmount);
event ContributionRefunded(uint256 indexed campaignId, address indexed contributor, uint256 ethAmount, uint256 tokenAmount);
```

`getCampaignStatus`/`getCampaign`/`getCampaigns`/`campaignCount` need no
structural changes beyond reading the new struct shape; `getCampaignStatus`'s
"successful" check becomes the same OR condition used in `withdraw`/`refund`.

### 2.4 Mock ERC20 token, for local testing

Anvil has no real USDC/DAI — we need our own. New file
`contract/src/MockERC20.sol`:

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract MockERC20 is ERC20 {
    uint8 private immutable _decimals;

    constructor(string memory name_, string memory symbol_, uint8 decimals_, uint256 initialSupply)
        ERC20(name_, symbol_)
    {
        _decimals = decimals_;
        _mint(msg.sender, initialSupply);
    }

    function decimals() public view override returns (uint8) {
        return _decimals;
    }

    // Unrestricted on purpose — this is a local test token, not real money.
    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}
```

**Deploy script** (`contract/script/DeployMockToken.s.sol`) — deploys the mock
token and mints a large balance to all 10 default Anvil accounts, so any test
account can act as a token contributor immediately:

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Script} from "forge-std/Script.sol";
import {MockERC20} from "../src/MockERC20.sol";

contract DeployMockTokenScript is Script {
    address[10] internal ANVIL_ACCOUNTS = [
        0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266,
        0x70997970C51812dc3A010C7d01b50e0d17dc79C8,
        0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC,
        0x90F79bf6EB2c4f870365E785982E1f101E93b906,
        0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65,
        0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc,
        0x976EA74026E726554dB657fA54763abd0C3a0aa9,
        0x14dC79964da2C08b23698B3D3cc7Ca32193d9955,
        0x23618e81E3f5cdF7f54C3d65f7FBc0aBf5B21E8f,
        0xa0Ee7A142d267C1f36714E4a8F75612F20a79720
    ];

    function run() external returns (address tokenAddress) {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");

        vm.startBroadcast(deployerPrivateKey);
        MockERC20 token = new MockERC20("Test USD Coin", "tUSDC", 6, 0);
        for (uint256 i = 0; i < ANVIL_ACCOUNTS.length; i++) {
            token.mint(ANVIL_ACCOUNTS[i], 1_000_000 * 10 ** 6); // 1,000,000 tUSDC each
        }
        vm.stopBroadcast();

        return address(token);
    }
}
```

6 decimals to deliberately mirror real USDC (and exercise the "decimals aren't
always 18" handling on the backend/frontend rather than accidentally only ever
testing with an 18-decimal token).

---

## 3. Backend changes

### 3.1 Regenerate bindings

After the contract changes, regenerate the abigen Go bindings — `contract.Campaign`
gains `CurrencyMode`, `Token`, `GoalEth`, `GoalToken`, `AmountRaisedEth`,
`AmountRaisedToken`; new `contract.MockERC20` bindings if we want the backend
to read token metadata (`Symbol`, `Decimals`) directly rather than hardcoding
it.

### 3.2 `models/campaign.go`

Add fields captured **at draft-creation time** (same pattern as `category`/
`durationDays` today):

```go
const (
    CurrencyModeEth   = "eth"
    CurrencyModeToken = "token"
    CurrencyModeBoth  = "both"
)

// on Campaign struct:
CurrencyMode  string  `gorm:"not null;default:eth" json:"currencyMode"`
TokenAddress  *string `json:"tokenAddress"`
TokenSymbol   *string `json:"tokenSymbol"`
TokenDecimals *uint8   `json:"tokenDecimals"`
GoalTokenRaw  *string `json:"goalToken"` // string, same reasoning as existing TargetEth (avoid float precision issues)
```

`CreateCampaign` gains `currencyMode`, `tokenAddress`, `goalToken` params;
validates:
- `currencyMode` is one of the three allowed values.
- If `token`/`both`, `tokenAddress` is a well-formed address and `goalToken` is
  a valid positive number.
- If `eth`/`both`, existing `targetEth` validation still applies.

### 3.3 Known-tokens list (mirrors `CampaignCategories`)

Following the same pattern already used for categories (shared, hardcoded list
on both backend and frontend to avoid an extra API call), add a small
allow-list of supported tokens rather than accepting arbitrary addresses —
safer, and avoids a whole class of "user pastes a broken/malicious ERC20"
problems:

```go
type SupportedToken struct {
    Symbol   string
    Address  string
    Decimals uint8
}

var SupportedTokens = []SupportedToken{
    {Symbol: "tUSDC", Address: "0x...", Decimals: 6}, // the mock token's deployed address
}
```

*(Open question to confirm with you: curated list only for now, or also allow
an "advanced: paste any ERC20 address" escape hatch? Recommend curated-only
for v1, revisit later.)*

### 3.4 `services/campaign_service.go`

- `CreateCampaignInput` gains `CurrencyMode`, `TokenAddress`, `GoalToken`.
- `PublishCampaignInput`/`PublishCampaign` — the existing on-chain verification
  (owner match) extends to also verify the on-chain `currencyMode`/`token`
  match what's stored in the DB draft, so a client can't publish a draft
  against a mismatched on-chain campaign.
- `GetMyCampaign`/response building — surface `goalEth`, `goalToken`,
  `amountRaisedEth`, `amountRaisedToken`, `tokenSymbol`, `tokenDecimals`
  instead of the current single `goal`/`amountRaised`.

### 3.5 `services/public_campaign_service.go`

`CampaignResponse` gains the same per-currency fields; `toCampaignResponse`
reads both `GoalEth`/`GoalToken`/`AmountRaisedEth`/`AmountRaisedToken` off the
chain struct. `campaignStatus()` helper's "successful" check becomes the same
OR condition used in the contract.

### 3.6 `models/transaction.go` + indexer

- Add `TokenAddress *string` (nullable — `nil` means ETH) to the `Transaction`
  model.
- `indexContributions` — reads the new `token` field straight off
  `ContributionMade`, one Transaction row per contribution as today.
- `indexWithdrawals`/`indexRefunds` — since `FundsWithdrawn`/
  `ContributionRefunded` now carry *two* amounts (`ethAmount`, `tokenAmount`)
  in one event, insert **up to two Transaction rows** per event (one if
  `ethAmount > 0`, another if `tokenAmount > 0`, using `LogIndex` +
  a suffix or a synthetic sub-index to keep the existing
  `(tx_hash, log_index)` uniqueness constraint intact — needs a small tweak,
  e.g. `log_index * 2` / `log_index * 2 + 1` for the two synthetic rows).

### 3.7 Amount formatting

Nothing new needed backend-side beyond exposing `tokenDecimals` — formatting
"raw units → human-readable" is a display concern, stays in the frontend (see
below), matching how ETH formatting already works today (backend always
returns raw wei strings).

---

## 4. Frontend changes

### 4.1 `lib/erc20Contract.js` (new)

Minimal ERC20 ABI + helper, parallel to `lib/crowdFundingContract.js`:

```js
const ERC20_ABI = [
  'function balanceOf(address) view returns (uint256)',
  'function allowance(address,address) view returns (uint256)',
  'function approve(address,uint256) returns (bool)',
  'function decimals() view returns (uint8)',
  'function symbol() view returns (string)',
]

export function getErc20Contract(address, providerOrSigner) {
  return new Contract(address, ERC20_ABI, providerOrSigner)
}
```

### 4.2 `utils/format.js`

Add a decimals-aware formatter alongside the existing (18-decimals-assumed)
`formatEth`/`formatEthDisplay`:

```js
export function formatTokenAmount(rawAmount, decimals, symbol) {
  return `${Number(formatUnits(rawAmount, decimals)).toLocaleString(undefined, { maximumFractionDigits: 4 })} ${symbol}`
}
```

### 4.3 Campaign creation (`CreateCampaignPage.jsx`)

- New "Accepted currency" selector: ETH / Token / Both (mirrors the existing
  `CATEGORIES` dropdown pattern).
- If Token/Both selected: a token dropdown (from the shared `SUPPORTED_TOKENS`
  list, mirrored from the backend's list) and a separate "Token goal" field.
- If ETH/Both selected: existing ETH goal field stays.
- `createMyCampaign` payload gains `currencyMode`, `tokenAddress`, `goalToken`.

### 4.4 Publish flow (`usePublishCampaign.js`)

- `createCampaign` on-chain call gains `currencyMode`, `token`, `goalEth`,
  `goalToken` params, sourced from the draft.

### 4.5 Contribution flow — the two-phase approve/contribute state machine

New hook `useContributeToken` (mirrors the existing `useWithdrawFunds`/
`usePublishCampaign` phase-machine pattern: `idle → checkingAllowance →
approving → confirmingApproval → contributing → confirming → idle/error`):

```js
async function contributeWithToken(campaignId, tokenAddress, amount, signer) {
  const token = getErc20Contract(tokenAddress, signer)
  const crowdFunding = getCrowdFundingContract(signer)
  const owner = await signer.getAddress()

  const currentAllowance = await token.allowance(owner, CONTRACT_ADDRESS)
  if (currentAllowance < amount) {
    const approveTx = await token.approve(CONTRACT_ADDRESS, amount)
    await approveTx.wait()
  }

  const tx = await crowdFunding.contributeToken(campaignId, amount)
  await tx.wait()
}
```

`ContributeForm.jsx` — when `campaign.currencyMode === 'both'`, add a
currency toggle (ETH / token symbol) before the amount input; route to
`contributeEth`/the new token flow accordingly. Show "Approve" as a distinct
UI step/label (matching `PUBLISH_LABELS`/`WITHDRAW_LABELS` conventions) so the
two-signature flow doesn't look broken.

Balance display: for token contributions, check `token.balanceOf(account)`
instead of (or alongside) the wallet's native ETH balance.

### 4.6 Display (`CampaignDetailsPage.jsx`, `CampaignPreview.jsx`)

- Replace the single progress bar with per-currency progress when in `both`
  mode (two small bars/labels: "X ETH raised of Y ETH goal" and "X tUSDC
  raised of Y tUSDC goal"), a single bar for `eth`/`token`-only campaigns.
- `computeProgressPercent` needs a currency-scoped variant (compute
  separately per currency, never combined).

### 4.7 Transaction tables (`CampaignTransactionsTab.jsx`, `UserTransactionsTab.jsx`)

- Show the currency/token symbol per row (ETH vs the token's symbol) instead
  of assuming everything is ETH.
- CSV export includes the currency column.

---

## 5. Open decisions before implementation

1. **Curated token list vs. arbitrary address input** — recommend curated-only
   for v1 (safer, mirrors the `CampaignCategories` pattern already in place).
2. **Fresh deploy vs. additive V2 struct** — recommend fresh deploy, given
   Anvil gets reset frequently in this project anyway.
3. Anything else you want to lock in before we start writing code.
