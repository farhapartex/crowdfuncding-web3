import { useEffect, useState } from 'react'
import Modal from './Modal'
import ContributeForm from './ContributeForm'
import WithdrawButton from './WithdrawButton'
import StatusBadge from './ui/StatusBadge'
import Button from './ui/Button'
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
      <p className="text-sm text-slate-600">{campaign.description}</p>

      <dl className="mt-4 grid grid-cols-2 gap-4">
        <div>
          <dt className="text-xs text-slate-500">Campaign ID</dt>
          <dd className="text-sm text-slate-900">{campaign.id}</dd>
        </div>
        <div>
          <dt className="text-xs text-slate-500">Owner</dt>
          <dd className="font-mono text-sm text-slate-900">{ownerDisplayName || shortenAddress(campaign.owner)}</dd>
        </div>
        <div>
          <dt className="text-xs text-slate-500">Status</dt>
          <dd className="mt-0.5">
            <StatusBadge status={campaign.status} />
          </dd>
        </div>
        <div>
          <dt className="text-xs text-slate-500">Goal</dt>
          <dd className="text-sm text-slate-900">{formatEth(campaign.goal)}</dd>
        </div>
        <div>
          <dt className="text-xs text-slate-500">Amount raised</dt>
          <dd className="text-sm text-slate-900">{formatEth(campaign.amountRaised)}</dd>
        </div>
        <div>
          <dt className="text-xs text-slate-500">Deadline</dt>
          <dd className="text-sm text-slate-900">{formatDate(campaign.deadline)}</dd>
        </div>
        <div>
          <dt className="text-xs text-slate-500">Withdrawn</dt>
          <dd className="text-sm text-slate-900">{campaign.withdrawn ? 'Yes' : 'No'}</dd>
        </div>
      </dl>

      {canContribute &&
        (account ? (
          <ContributeForm
            onContribute={(amountEth) => onContribute(campaign.id, amountEth)}
            isContributing={isContributing}
          />
        ) : (
          <div className="mt-4 flex flex-col gap-3 border-t border-slate-200 pt-4">
            <p className="text-sm text-slate-600">Connect your wallet to contribute to this campaign.</p>
            <Button type="button" onClick={onConnectWallet}>
              Connect Wallet
            </Button>
          </div>
        ))}

      {canWithdraw && (
        <WithdrawButton onWithdraw={() => onWithdraw(campaign.id)} isWithdrawing={isWithdrawing} />
      )}
    </Modal>
  )
}

export default CampaignDetailsModal
