import { formatEth, formatDate } from '../utils/format'

function CampaignTable({ campaigns, onSelect }) {
  if (campaigns.length === 0) {
    return <p className="empty-state">No campaigns yet. Create the first one.</p>
  }

  return (
    <table className="campaign-table">
      <thead>
        <tr>
          <th>Title</th>
          <th>Goal</th>
          <th>Raised</th>
          <th>Status</th>
          <th>Deadline</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        {campaigns.map((campaign) => (
          <tr key={campaign.id}>
            <td>{campaign.title}</td>
            <td>{formatEth(campaign.goal)}</td>
            <td>{formatEth(campaign.amountRaised)}</td>
            <td>
              <span className={`status-badge status-${campaign.status.toLowerCase()}`}>{campaign.status}</span>
            </td>
            <td>{formatDate(campaign.deadline)}</td>
            <td>
              <button type="button" className="link-button" onClick={() => onSelect(campaign.id)}>
                Details
              </button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}

export default CampaignTable
