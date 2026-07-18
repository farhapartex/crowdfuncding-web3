# Crowd Funding

## What is this project?

This is a crowdfunding website, like FundMe, but built on the blockchain.

Anyone can create a campaign asking for money. Other people can then give
money (called "contributing") to campaigns they like. All the important
rules (how much has been raised, who can take the money out, when to give
refunds) live inside a smart contract on the blockchain. This means no one,
not even the person running the website, can secretly touch the money.

People can contribute using **ETH** (the normal currency of Ethereum) or an
**ERC20 token** (a different kind of coin that lives on the same
blockchain, like a stablecoin). The campaign creator picks which one(s)
they accept.

### Key Features

- Create a campaign and choose how people can pay: **ETH only**, **Token
  only**, or **Both**.
- Contribute to a campaign using ETH or an ERC20 token.
- Money is locked in the smart contract. It can only be taken out
  ("withdrawn") by the campaign owner, and only after the goal is reached.
- If a campaign fails to reach its goal in time, contributors can get a
  refund, in the exact same currency they paid with.
- Log in with email/Google (using Auth0), then connect your crypto wallet.
- Add a title, description, cover photo, category, and country to each
  campaign.
- See a live list of contributors and every transaction (contribute,
  withdraw, refund) on a campaign.
- Comment on a campaign and reply to other people's comments.
- Archive a campaign once it's finished.
- Watch live logs and errors from the backend, the comment service, and
  even the website itself (in the browser), using Grafana.

### The Four Parts of This Project

- `contract/` – the smart contract (written in Solidity, using the Foundry
  tool)
- `backend/` – a Go API that reads the blockchain and stores extra info
  (like descriptions and photos) in a database
- `comment/` – a small Go service that only handles comments, talks to the
  backend
- `frontend/` – the actual website (built with React and Vite)

## How to Set It Up

Set things up in this order: **contract first, then backend, then
frontend**. Each part depends on the one before it.

### 1. Smart Contract (Solidity)

The smart contract is the "source of truth" — it holds the real money and
the real rules. Everything else (backend, frontend) just reads from it or
sends transactions to it.

We use [Foundry](https://book.getfoundry.sh/getting-started/installation)
to build, test, and deploy it. Install Foundry first if you don't have it.

```shell
cd contract
forge build
forge test
```

Now let's run a fake, private blockchain on your own computer to test on
(this is called **Anvil**, and it comes with Foundry):

1. Open a new terminal and run:

   ```shell
   anvil
   ```

   This starts a fake blockchain and prints out 10 test accounts, each
   with 10,000 fake ETH and a private key. Keep this terminal open. Copy
   the private key of "Account #0" — you'll use it as the deployer.

2. Copy the example env file and fill it in:

   ```shell
   cp .env.example .env
   ```

   - `RPC_URL` — leave as `http://127.0.0.1:8545` (this is Anvil's address)
   - `PRIVATE_KEY` — paste the private key you copied above

3. Deploy the contract:

   ```shell
   forge script script/DeployProxy.s.sol:DeployProxyScript --rpc-url anvil --broadcast
   ```

   This prints two addresses. Copy the one called `proxyAddress` — that's
   the address you'll use everywhere else (backend, frontend, MetaMask).
   You can ignore `implementationAddress`.

4. (Optional, only needed if you want to test paying with a token) Deploy
   a fake test token:

   ```shell
   forge script script/DeployMockToken.s.sol:DeployMockTokenScript --rpc-url anvil --broadcast
   ```

   This creates a fake stablecoin called `tUSDC` and gives 1,000,000
   tUSDC to each of the 10 Anvil test accounts. Copy the token address
   that gets printed — you'll need it later.

If you change anything in the contract, rebuild it and regenerate the
files the backend and frontend use to talk to it:

```shell
forge build
jq '.abi' out/CrowdFunding.sol/CrowdFunding.json > ../frontend/src/contract/CrowdFundingAbi.json
jq '.abi' out/CrowdFunding.sol/CrowdFunding.json > ../backend/contract/CrowdFunding.abi.json
cd ../backend/contract
abigen --abi=CrowdFunding.abi.json --pkg=contract --type=CrowdFunding --out=crowdfunding.go
```

#### Add this test network to MetaMask

So you can actually click around the website with a real wallet, add your
local Anvil blockchain to MetaMask (or any wallet you use):

1. Open MetaMask → click the network dropdown → **Add network** → **Add a
   network manually**.
2. Fill in:
   - **Network name**: anything you like, e.g. `Anvil Local`
   - **New RPC URL**: `http://127.0.0.1:8545`
   - **Chain ID**: `31337`
   - **Currency symbol**: `ETH`
3. Save, then switch to this new network.
4. Import a test account: click the account icon → **Import Account** →
   paste one of the private keys Anvil printed (pick a different one than
   the deployer, e.g. "Account #1", so you have someone who can act as a
   contributor). Now MetaMask has fake ETH (and fake tUSDC, if you
   deployed it) ready to use.

### 2. Backend (Go API)

The backend reads campaign data straight from the blockchain, and keeps
extra info in its own database (titles, descriptions, photos, comments,
login sessions, etc. — things that don't need to live on the blockchain).

1. Copy the example env file:

   ```shell
   cd backend
   cp .env.example .env
   ```

2. Fill in `.env`:
   - `RPC_URL` — same as the contract, `http://127.0.0.1:8545`
   - `CONTRACT_ADDRESS` — the `proxyAddress` you copied during contract
     setup
   - `JWT_SECRET` — any long random text (used to sign login sessions)
   - `SCOPE_MASK_SECRET` — any long random text (used to hide real
     database IDs from the outside world)
   - `POSTGRES_USER` / `POSTGRES_PASSWORD` / `POSTGRES_DB` /
     `POSTGRES_PORT` — fine to leave as the defaults
   - `COMMENT_SERVICE_ADDR` / `COMMENT_SERVICE_TOKEN` — fine to leave as
     the defaults for local use
   - `SUPPORTED_TOKENS` — optional, only needed for token campaigns (see
     the [testing section](#testing-with-different-currencies) below)

3. Set up login (Auth0): the website uses [Auth0](https://auth0.com) (a
   free login service) so people can sign in with email/Google/etc.
   - Sign up for a free Auth0 account.
   - Create an **Application** (type: *Single Page Application*).
   - Create an **API** (Applications → APIs → Create API) — give it any
     name and an identifier like `https://crowdfunding-api`. This
     identifier is the "audience".
   - From the Application's settings, copy the **Domain** and
     **Client ID**.
   - Put these into `.env`: `AUTH0_APP_DOMAIN` (the Domain),
     `AUTH0_AUDIENCE` (the API identifier you made up).
   - You'll reuse the Domain, Client ID, and audience again in the
     frontend setup below.

4. Set up photo storage (Cloudflare R2): campaign cover photos are stored
   on [Cloudflare R2](https://developers.cloudflare.com/r2/) (cheap,
   S3-like storage — the free tier is enough for testing).
   - Sign up for a free Cloudflare account and turn on R2.
   - Create a bucket (any name).
   - Create an API token for R2 with read/write access — this gives you
     an Access Key ID and Secret Access Key.
   - Turn on public access for the bucket (or connect your own domain) so
     photos can be viewed from a browser — this gives you a Public URL.
   - Put these into `.env`: `R2_ACCOUNT_ID`, `R2_ACCESS_KEY_ID`,
     `R2_SECRET_ACCESS_KEY`, `R2_BUCKET`, `R2_PUBLIC_URL`, and
     `R2_S3_API` (this looks like
     `https://<your-account-id>.r2.cloudflarestorage.com`).

5. Start everything with Docker (this starts Postgres, the backend, and
   the comment service together):

   ```shell
   cd ..
   docker compose up -d --build
   ```

   The backend automatically creates its database tables on startup, and
   reaches your Anvil node (running on your computer, not in Docker)
   through `host.docker.internal` instead of `127.0.0.1`.

6. Check it's alive:

   ```shell
   curl http://localhost:8080/health
   ```

   If you'd rather run the backend directly on your computer instead of
   in Docker (useful while making code changes), run just Postgres in
   Docker and the backend on its own:

   ```shell
   docker compose up -d postgres
   go run .
   ```

#### Optional: Watch logs with Grafana

This project also comes with Grafana + Loki, so you can see live logs
from the backend, the comment service, and even errors from people's
browsers, all in one place.

```shell
docker compose up -d loki promtail grafana alloy
```

Open [http://localhost:3000](http://localhost:3000) and go to **Explore**
→ pick the **Loki** datasource. Try a query like `{service="backend"}` or
`{app_name="crowdfunding-frontend"}`.

### 3. Frontend (the website)

1. Copy the example env file:

   ```shell
   cd frontend
   cp .env.example .env
   ```

2. Fill in `.env`:
   - `VITE_CONTRACT_ADDRESS` — the same `proxyAddress` as the backend
   - `VITE_API_BASE_URL` — `http://localhost:8080`
   - `VITE_AUTH0_DOMAIN` / `VITE_AUTH0_CLIENT_ID` / `VITE_AUTH0_AUDIENCE`
     — the same Domain, Client ID, and audience from the Auth0
     application you made in the backend step
   - `VITE_FARO_URL` — only needed if you set up Grafana above, leave as
     `http://localhost:12347/collect`

3. In your Auth0 Application's settings, add these so login actually
   works:
   - **Allowed Callback URLs**: `http://localhost:5173/auth/callback`
   - **Allowed Logout URLs**: `http://localhost:5173`
   - **Allowed Web Origins**: `http://localhost:5173`

4. Install and run:

   ```shell
   npm install
   npm run dev
   ```

5. Open [http://localhost:5173](http://localhost:5173) in the same
   browser where you set up your MetaMask test account.

## Testing With Different Currencies

Once everything above is running, here's how to try out all three ways of
paying: ETH only, Token only, and Both.

1. Make sure you deployed the fake token (contract setup, step 4) and put
   its address into the backend's `SUPPORTED_TOKENS`, like this:

   ```
   SUPPORTED_TOKENS=[{"symbol":"tUSDC","address":"<token address here>","decimals":6}]
   ```

   Restart the backend after changing this.

2. Log in on the website and go to **Create Campaign**.
3. Fill in the basic details, then in the **Funding** section, pick one:
   - **ETH** — you only set an ETH goal. People pay with ETH.
   - **Token** — you pick a token from the list and set a token goal.
     People pay with that token.
   - **Both** — you set both an ETH goal and a token goal. People can pay
     with either one.
4. Click **Publish**. MetaMask will pop up **once**, asking you to
   confirm creating the campaign. This always costs a small ETH gas fee,
   even for token campaigns — gas is a separate thing from the campaign's
   currency, and it's always paid in ETH.
5. Now try contributing (open the campaign page and use the "Contribute"
   box):
   - **ETH**: MetaMask asks you to confirm **once**, and you're done.
   - **Token**: MetaMask asks you to confirm **twice**:
     1. **Approve** — lets the app move that amount of your tokens.
     2. **Contribute** — actually sends the tokens.

     (This is normal for ERC20 tokens. If you contribute again with a
     different amount later, you'll likely see the "Approve" step again.)
   - **Both**: you'll see a small switch to pick ETH or the token before
     typing in an amount.
6. Don't have any test tokens? The `DeployMockToken` script already gave
   1,000,000 tUSDC to each of the 10 default Anvil accounts. If you're
   using a different account, the test token's `mint` function is open to
   anyone, so you (or the app) could give any address more of it — it's
   fake money for testing only.
7. **Withdraw**: once the ETH goal *or* the token goal is reached
   (whichever happens first), the campaign owner can click
   **Withdraw Funds** and receive everything raised so far, in whichever
   currencies were used.
8. **Refund**: let the campaign's deadline pass without reaching either
   goal. On Anvil, you can skip time forward yourself instead of waiting:

   ```shell
   cast rpc evm_increaseTime 604801
   cast rpc evm_mine
   ```

   After that, contributors can request a refund, and they always get
   back the exact same currency they originally paid with.
