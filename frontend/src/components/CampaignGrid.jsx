import { formatEth } from '../utils/format'

function CampaignGrid({ campaigns, onSelect }) {
  if (campaigns.length === 0) {
    return (
      <div className="rounded-xl border border-dashed border-slate-300 bg-white px-6 py-16 text-center">
        <p className="text-sm text-slate-500">No campaigns yet. Be the first to create one.</p>
      </div>
    )
  }

  return (
    <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
      {campaigns.map((campaign) => (
        <button
          key={campaign.id}
          type="button"
          onClick={() => onSelect(campaign.id)}
          className="overflow-hidden rounded-xl border border-slate-200 bg-white text-left shadow-sm transition-shadow hover:shadow-md"
        >
          <div className="flex aspect-video w-full items-center justify-center bg-gradient-to-br from-indigo-50 to-indigo-100">
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="1.5"
              className="h-10 w-10 text-indigo-300"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M3 16.5V6.75A2.25 2.25 0 0 1 5.25 4.5h13.5A2.25 2.25 0 0 1 21 6.75v9.75M3 16.5l4.72-4.72a1.5 1.5 0 0 1 2.12 0l3.66 3.66a1.5 1.5 0 0 0 2.12 0l1.66-1.66a1.5 1.5 0 0 1 2.12 0L21 16.5M3 16.5V18a2.25 2.25 0 0 0 2.25 2.25h13.5A2.25 2.25 0 0 0 21 18v-1.5"
              />
            </svg>
          </div>
          <div className="p-4">
            <h3 className="truncate font-medium text-slate-900">{campaign.title}</h3>
            <p className="mt-1 text-sm text-slate-500">{formatEth(campaign.amountRaised)} raised</p>
          </div>
        </button>
      ))}
    </div>
  )
}

export default CampaignGrid
