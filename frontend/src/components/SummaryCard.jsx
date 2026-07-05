import { shortenAddress } from '../utils/format'

function SummaryCard({ account, campaignCount }) {
  return (
    <div className="summary-card">
      <div className="summary-item">
        <span className="label">Connected account</span>
        <span className="value mono">{shortenAddress(account)}</span>
      </div>
      <div className="summary-item">
        <span className="label">Campaign count</span>
        <span className="value">{campaignCount}</span>
      </div>
    </div>
  )
}

export default SummaryCard
