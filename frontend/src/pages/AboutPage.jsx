function AboutPage() {
  return (
    <div className="mx-auto max-w-2xl">
      <h1 className="text-2xl font-semibold text-slate-900">What is Crowd Funding</h1>
      <p className="mt-3 leading-relaxed text-slate-600">
        Crowd funding is a way to collect small amounts of money from many people to reach one big goal. Someone
        creates a campaign and says how much money they need and by when. Other people send money to that campaign
        if they like the idea. If the campaign reaches its goal before the deadline, the owner of the campaign can
        take the money out. If the campaign does not reach the goal before the deadline, the people who sent money
        can get their money back.
      </p>

      <h2 className="mt-10 text-xl font-semibold text-slate-900">How This App Works</h2>
      <p className="mt-3 leading-relaxed text-slate-600">
        This app runs entirely on a smart contract on the Ethereum blockchain, so nobody, including the people who
        built it, can secretly change the rules or move funds around.
      </p>
      <ul className="mt-4 flex flex-col gap-3">
        <li className="rounded-lg border border-slate-200 bg-white p-4">
          <span className="font-medium text-slate-900">Create a campaign</span>
          <p className="mt-1 text-sm text-slate-600">
            Anyone can create a campaign with a title, description, funding goal, and deadline.
          </p>
        </li>
        <li className="rounded-lg border border-slate-200 bg-white p-4">
          <span className="font-medium text-slate-900">Contribute</span>
          <p className="mt-1 text-sm text-slate-600">
            Anyone can send ETH to a campaign they want to support, as long as its deadline hasn't passed.
          </p>
        </li>
        <li className="rounded-lg border border-slate-200 bg-white p-4">
          <span className="font-medium text-slate-900">Withdraw</span>
          <p className="mt-1 text-sm text-slate-600">
            If a campaign reaches its goal, only its owner can withdraw the funds.
          </p>
        </li>
        <li className="rounded-lg border border-slate-200 bg-white p-4">
          <span className="font-medium text-slate-900">Refund</span>
          <p className="mt-1 text-sm text-slate-600">
            If a campaign's deadline passes without reaching its goal, contributors can get their own contribution
            back.
          </p>
        </li>
      </ul>

      <h2 className="mt-10 text-xl font-semibold text-slate-900">About This Project</h2>
      <p className="mt-3 leading-relaxed text-slate-600">
        This is a learning project built to explore Solidity, Foundry, Go, and React together end to end.
      </p>
    </div>
  )
}

export default AboutPage
