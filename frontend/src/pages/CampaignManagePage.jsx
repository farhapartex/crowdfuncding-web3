import { useEffect, useState } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { parseEther } from 'ethers'
import { fetchCampaign, fetchContributors } from '../lib/api'
import { getCrowdFundingContract } from '../lib/crowdFundingContract'
import { shortenAddress, formatEth, formatDate } from '../utils/format'
import StatusBadge from '../components/ui/StatusBadge'
import WithdrawButton from '../components/WithdrawButton'
import ContributeForm from '../components/ContributeForm'

function DetailsTab({ campaign, canContribute, isContributing, onContribute, canWithdraw, isWithdrawing, onWithdraw }) {
  return (
    <div className="rounded-xl border border-slate-200 bg-white p-6 shadow-sm">
      <p className="text-sm text-slate-600">{campaign.description}</p>

      <dl className="mt-4 grid grid-cols-2 gap-4">
        <div>
          <dt className="text-xs text-slate-500">Campaign ID</dt>
          <dd className="text-sm text-slate-900">{campaign.id}</dd>
        </div>
        <div>
          <dt className="text-xs text-slate-500">Owner</dt>
          <dd className="font-mono text-sm text-slate-900">{shortenAddress(campaign.owner)}</dd>
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

      {canContribute && <ContributeForm onContribute={onContribute} isContributing={isContributing} />}
      {canWithdraw && <WithdrawButton onWithdraw={onWithdraw} isWithdrawing={isWithdrawing} />}
    </div>
  )
}

function ContributorsTab({ campaignId, setError }) {
  const [contributors, setContributors] = useState(null)

  useEffect(() => {
    fetchContributors(campaignId)
      .then(setContributors)
      .catch((err) => setError(err.message))
  }, [campaignId])

  if (!contributors) {
    return <p className="text-sm text-slate-500">Loading contributors...</p>
  }

  if (contributors.length === 0) {
    return (
      <div className="rounded-xl border border-dashed border-slate-300 bg-white px-6 py-16 text-center">
        <p className="text-sm text-slate-500">No contributions yet.</p>
      </div>
    )
  }

  return (
    <div className="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b border-slate-200 bg-slate-50">
            <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wide text-slate-500">
              Contributor
            </th>
            <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wide text-slate-500">
              Amount
            </th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-100">
          {contributors.map((contributor) => (
            <tr key={contributor.address}>
              <td className="px-4 py-3 font-mono text-slate-900">
                {contributor.displayName || shortenAddress(contributor.address)}
              </td>
              <td className="px-4 py-3 text-slate-600">{formatEth(contributor.amount)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function CampaignManagePage({ provider, account, setError, showToast }) {
  const { id } = useParams()
  const navigate = useNavigate()
  const [campaign, setCampaign] = useState(null)
  const [activeTab, setActiveTab] = useState('details')
  const [isWithdrawing, setIsWithdrawing] = useState(false)
  const [isContributing, setIsContributing] = useState(false)

  useEffect(() => {
    fetchCampaign(id)
      .then(setCampaign)
      .catch((err) => setError(err.message))
  }, [id])

  useEffect(() => {
    if (!campaign) return

    const isOwner = account && account.toLowerCase() === campaign.owner.toLowerCase()
    if (!isOwner) {
      navigate('/', { replace: true })
    }
  }, [campaign, account, navigate])

  async function handleContribute(amountEth) {
    setError(null)
    setIsContributing(true)

    try {
      const signer = await provider.getSigner()
      const crowdFunding = getCrowdFundingContract(signer)

      const amountInWei = parseEther(amountEth)

      const tx = await crowdFunding.contribute(id, { value: amountInWei })
      await tx.wait()

      const updated = await fetchCampaign(id)
      setCampaign(updated)
      showToast('Contribution successful!')
    } catch (err) {
      setError(err.shortMessage || err.message)
      throw err
    } finally {
      setIsContributing(false)
    }
  }

  async function handleWithdraw() {
    setError(null)
    setIsWithdrawing(true)

    try {
      const signer = await provider.getSigner()
      const crowdFunding = getCrowdFundingContract(signer)

      const tx = await crowdFunding.withdraw(id)
      await tx.wait()

      const updated = await fetchCampaign(id)
      setCampaign(updated)
      showToast('Funds withdrawn successfully.')
    } catch (err) {
      setError(err.shortMessage || err.message)
    } finally {
      setIsWithdrawing(false)
    }
  }

  if (!campaign) {
    return <p className="text-sm text-slate-500">Loading campaign...</p>
  }

  const isOwner = account && account.toLowerCase() === campaign.owner.toLowerCase()

  if (!isOwner) {
    return null
  }

  const canContribute = Date.now() / 1000 < Number(campaign.deadline)
  const canWithdraw = campaign.status === 'Successful' && !campaign.withdrawn

  return (
    <div className="flex flex-col gap-5">
      <div>
        <Link to="/" className="text-sm text-indigo-600 hover:text-indigo-500">
          &larr; Back to campaigns
        </Link>
        <h1 className="mt-2 text-xl font-semibold text-slate-900">{campaign.title}</h1>
      </div>

      <div className="flex gap-6 border-b border-slate-200">
        <button
          type="button"
          onClick={() => setActiveTab('details')}
          className={`border-b-2 px-1 pb-3 text-sm font-medium transition-colors ${
            activeTab === 'details'
              ? 'border-indigo-600 text-indigo-600'
              : 'border-transparent text-slate-500 hover:text-slate-700'
          }`}
        >
          Details
        </button>
        <button
          type="button"
          onClick={() => setActiveTab('contributors')}
          className={`border-b-2 px-1 pb-3 text-sm font-medium transition-colors ${
            activeTab === 'contributors'
              ? 'border-indigo-600 text-indigo-600'
              : 'border-transparent text-slate-500 hover:text-slate-700'
          }`}
        >
          Contributors
        </button>
      </div>

      {activeTab === 'details' ? (
        <DetailsTab
          campaign={campaign}
          canContribute={canContribute}
          isContributing={isContributing}
          onContribute={handleContribute}
          canWithdraw={canWithdraw}
          isWithdrawing={isWithdrawing}
          onWithdraw={handleWithdraw}
        />
      ) : (
        <ContributorsTab campaignId={id} setError={setError} />
      )}
    </div>
  )
}

export default CampaignManagePage
