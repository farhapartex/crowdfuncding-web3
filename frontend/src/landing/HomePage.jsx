import { Link } from 'react-router-dom'

function HomePage() {
  return (
    <div className="flex flex-col gap-16">
      <section className="flex flex-col items-center gap-6 py-12 text-center">
        <h1 className="max-w-2xl text-4xl font-bold text-slate-900 sm:text-5xl">
          Bring your ideas to life, together
        </h1>
        <p className="max-w-xl text-lg text-slate-600">
          Create a campaign, share your goal, and let people who believe in it help you reach it —
          secured by a smart contract, not a middleman.
        </p>
        <div className="flex gap-3">
          <Link
            to="/campaigns"
            className="rounded-lg bg-indigo-600 px-5 py-2.5 text-sm font-medium text-white shadow-sm hover:bg-indigo-500"
          >
            Explore Campaigns
          </Link>
          <Link
            to="/campaigns"
            className="rounded-lg border border-slate-300 bg-white px-5 py-2.5 text-sm font-medium text-slate-700 hover:bg-slate-50"
          >
            Start a Campaign
          </Link>
        </div>
      </section>

      <section className="grid gap-6 sm:grid-cols-3">
        <div className="rounded-xl border border-slate-200 bg-white p-6">
          <span className="text-sm font-semibold text-indigo-600">Step 1</span>
          <h3 className="mt-2 text-lg font-semibold text-slate-900">Create a campaign</h3>
          <p className="mt-2 text-sm text-slate-600">
            Set a title, a goal, and a deadline. Anyone can start one — no approval needed.
          </p>
        </div>
        <div className="rounded-xl border border-slate-200 bg-white p-6">
          <span className="text-sm font-semibold text-indigo-600">Step 2</span>
          <h3 className="mt-2 text-lg font-semibold text-slate-900">Get supported</h3>
          <p className="mt-2 text-sm text-slate-600">
            Anyone can contribute directly from their wallet — no account required to give.
          </p>
        </div>
        <div className="rounded-xl border border-slate-200 bg-white p-6">
          <span className="text-sm font-semibold text-indigo-600">Step 3</span>
          <h3 className="mt-2 text-lg font-semibold text-slate-900">Reach your goal</h3>
          <p className="mt-2 text-sm text-slate-600">
            If the goal is met, the owner withdraws the funds. If not, contributors are refunded.
          </p>
        </div>
      </section>

      <section className="flex flex-col items-center gap-4 rounded-xl bg-slate-900 px-6 py-12 text-center">
        <h2 className="text-2xl font-semibold text-white">Ready to get started?</h2>
        <p className="max-w-md text-sm text-slate-300">
          Connect your wallet, sign in, and create your first campaign in a couple of minutes.
        </p>
        <Link
          to="/campaigns"
          className="rounded-lg bg-white px-5 py-2.5 text-sm font-medium text-slate-900 hover:bg-slate-100"
        >
          Browse Campaigns
        </Link>
      </section>
    </div>
  )
}

export default HomePage
