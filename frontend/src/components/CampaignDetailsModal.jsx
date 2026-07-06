import { useEffect, useState } from 'react'
import Modal from './Modal'
import ContributeForm from './ContributeForm'
import WithdrawButton from './WithdrawButton'
import { shortenAddress, formatEth, formatDate } from '../utils/format'
import { fetchPublicProfile } from '../lib/api'

function CampaignDetailsModal({
  campaign,
  account,
  onConnectWallet,
  onContribute,
  isContributing,
  onWithdraw,
  isWithdrawing,
  onClose,
}) {
  const [ownerDisplayName, setOwnerDisplayName] = useState('')

  useEffect(() => {
    fetchPublicProfile(campaign.owner)
      .then(({ displayName }) => setOwnerDisplayName(displayName))
      .catch(() => setOwnerDisplayName(''))
  }, [campaign.owner])

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
          <dd className="mono">{ownerDisplayName || shortenAddress(campaign.owner)}</dd>
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

      {canContribute &&
        (account ? (
          <ContributeForm
            onContribute={(amountEth) => onContribute(campaign.id, amountEth)}
            isContributing={isContributing}
          />
        ) : (
          <div className="contribute-form">
            <p className="campaign-description">Connect your wallet to contribute to this campaign.</p>
            <button type="button" onClick={onConnectWallet}>
              Connect Wallet
            </button>
          </div>
        ))}

      {canWithdraw && (
        <WithdrawButton onWithdraw={() => onWithdraw(campaign.id)} isWithdrawing={isWithdrawing} />
      )}
    </Modal>
  )
}

export default CampaignDetailsModal
