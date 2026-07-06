import { formatEth, formatDate } from '../utils/format'
import StatusBadge from './ui/StatusBadge'
import Button from './ui/Button'

function CampaignTable({ campaigns, onSelect }) {
  if (campaigns.length === 0) {
    return (
      <div className="rounded-xl border border-dashed border-slate-300 bg-white px-6 py-16 text-center">
        <p className="text-sm text-slate-500">No campaigns yet. Be the first to create one.</p>
      </div>
    )
  }

  return (
    <div className="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b border-slate-200 bg-slate-50">
            <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wide text-slate-500">Title</th>
            <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wide text-slate-500">Goal</th>
            <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wide text-slate-500">Raised</th>
            <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wide text-slate-500">Status</th>
            <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wide text-slate-500">
              Deadline
            </th>
            <th className="px-4 py-3"></th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-100">
          {campaigns.map((campaign) => (
            <tr key={campaign.id} className="transition-colors hover:bg-slate-50">
              <td className="px-4 py-3 font-medium text-slate-900">{campaign.title}</td>
              <td className="px-4 py-3 text-slate-600">{formatEth(campaign.goal)}</td>
              <td className="px-4 py-3 text-slate-600">{formatEth(campaign.amountRaised)}</td>
              <td className="px-4 py-3">
                <StatusBadge status={campaign.status} />
              </td>
              <td className="px-4 py-3 text-slate-600">{formatDate(campaign.deadline)}</td>
              <td className="px-4 py-3 text-right">
                <Button variant="link" onClick={() => onSelect(campaign.id)}>
                  Details
                </Button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default CampaignTable
