import Modal from './Modal'
import ContributeForm from './ContributeForm'
import WithdrawButton from './WithdrawButton'
import { shortenAddress, formatEth, formatDate } from '../utils/format'

function CampaignDetailsModal({ campaign, account, onContribute, isContributing, onWithdraw, isWithdrawing, onClose }) {
  const canContribute = Date.now() / 1000 < Number(campaign.deadline)
  const isOwner = account?.toLowerCase() === campaign.owner.toLowerCase()
  const canWithdraw = isOwner && campaign.status === 'Successful' && !campaign.withdrawn

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

      {canContribute && (
        <ContributeForm
          onContribute={(amountEth) => onContribute(campaign.id, amountEth)}
          isContributing={isContributing}
        />
      )}

      {canWithdraw && (
        <WithdrawButton onWithdraw={() => onWithdraw(campaign.id)} isWithdrawing={isWithdrawing} />
      )}
    </Modal>
  )
}

export default CampaignDetailsModal
