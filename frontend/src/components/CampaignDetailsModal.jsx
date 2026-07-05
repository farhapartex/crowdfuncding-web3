import Modal from './Modal'
import { shortenAddress, formatEth, formatDate } from '../utils/format'

function CampaignDetailsModal({ campaign, onClose }) {
  return (
    <Modal title={campaign.title} onClose={onClose}>
      <p className="campaign-description">{campaign.description}</p>

      <dl className="campaign-details">
        <div>
          <dt>Campaign ID</dt>
          <dd>{campaign.id}</dd>
        </div>
        <div>
          <dt>Owner</dt>
          <dd className="mono">{shortenAddress(campaign.owner)}</dd>
        </div>
        <div>
          <dt>Status</dt>
          <dd>{campaign.status}</dd>
        </div>
        <div>
          <dt>Goal</dt>
          <dd>{formatEth(campaign.goal)}</dd>
        </div>
        <div>
          <dt>Amount raised</dt>
          <dd>{formatEth(campaign.amountRaised)}</dd>
        </div>
        <div>
          <dt>Deadline</dt>
          <dd>{formatDate(campaign.deadline)}</dd>
        </div>
        <div>
          <dt>Withdrawn</dt>
          <dd>{campaign.withdrawn ? 'Yes' : 'No'}</dd>
        </div>
      </dl>
    </Modal>
  )
}

export default CampaignDetailsModal
