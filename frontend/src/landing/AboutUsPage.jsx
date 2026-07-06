function AboutUsPage() {
  return (
    <div className="mx-auto max-w-2xl">
      <h1 className="text-2xl font-semibold text-slate-900">About Us</h1>
      <p className="mt-3 leading-relaxed text-slate-600">
        We built this platform to make crowdfunding simple, transparent, and open to anyone — no
        middleman holding your funds, no approval process, and no hidden fees. Every campaign,
        contribution, and withdrawal happens directly on the blockchain, where anyone can verify
        it.
      </p>

      <h2 className="mt-10 text-xl font-semibold text-slate-900">What is Crowd Funding</h2>
      <p className="mt-3 leading-relaxed text-slate-600">
        Crowd funding is a way to collect small amounts of money from many people to reach one big
        goal. Someone creates a campaign and says how much money they need and by when. Other
        people send money to that campaign if they like the idea. If the campaign reaches its goal
        before the deadline, the owner can take the money out. If it doesn't, contributors can get
        their money back.
      </p>

      <h2 className="mt-10 text-xl font-semibold text-slate-900">How It Works Here</h2>
      <ul className="mt-4 flex flex-col gap-3">
        <li className="rounded-lg border border-slate-200 bg-white p-4">
          <span className="font-medium text-slate-900">Create a campaign</span>
          <p className="mt-1 text-sm text-slate-600">
            Sign in, then set a title, description, funding goal, and deadline.
          </p>
        </li>
        <li className="rounded-lg border border-slate-200 bg-white p-4">
          <span className="font-medium text-slate-900">Contribute</span>
          <p className="mt-1 text-sm text-slate-600">
            Anyone can send ETH to a campaign directly from their wallet, before its deadline.
          </p>
        </li>
        <li className="rounded-lg border border-slate-200 bg-white p-4">
          <span className="font-medium text-slate-900">Withdraw</span>
          <p className="mt-1 text-sm text-slate-600">
            Once a campaign reaches its goal, only its owner can withdraw the funds.
          </p>
        </li>
        <li className="rounded-lg border border-slate-200 bg-white p-4">
          <span className="font-medium text-slate-900">Refund</span>
          <p className="mt-1 text-sm text-slate-600">
            If a campaign's deadline passes without reaching its goal, contributors can claim
            their contribution back.
          </p>
        </li>
      </ul>

      <h2 className="mt-10 text-xl font-semibold text-slate-900">About This Project</h2>
      <p className="mt-3 leading-relaxed text-slate-600">
        This is a learning project built to explore Solidity, Foundry, Go, and React together end
        to end.
      </p>
    </div>
  )
}

export default AboutUsPage
