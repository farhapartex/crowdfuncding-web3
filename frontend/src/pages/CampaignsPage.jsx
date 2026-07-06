import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { parseEther } from 'ethers'
import { getCrowdFundingContract } from '../lib/crowdFundingContract'
import { fetchCampaigns } from '../lib/api'
import CreateCampaignForm from '../components/CreateCampaignForm'
import CampaignTable from '../components/CampaignTable'
import Pagination from '../components/Pagination'
import Modal from '../components/Modal'
import CampaignDetailsModal from '../components/CampaignDetailsModal'
import Button from '../components/ui/Button'

const SECONDS_PER_DAY = 24 * 60 * 60
const PAGE_SIZE = 10

function CampaignsPage({ account, provider, onConnectWallet, setError, showToast }) {
  const navigate = useNavigate()
  const [campaigns, setCampaigns] = useState([])
  const [totalCampaigns, setTotalCampaigns] = useState(0)
  const [offset, setOffset] = useState(0)
  const [selectedCampaignId, setSelectedCampaignId] = useState(null)
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [isCreating, setIsCreating] = useState(false)
  const [isContributing, setIsContributing] = useState(false)

  useEffect(() => {
    refreshCampaigns(0)
  }, [])

  async function refreshCampaigns(targetOffset = offset) {
    const { campaigns: result, total } = await fetchCampaigns({ offset: targetOffset, limit: PAGE_SIZE })
    setCampaigns(result)
    setTotalCampaigns(total)
    setOffset(targetOffset)
  }

  async function handleCreateCampaign({ title, description, goalEth, durationDays }) {
    setError(null)
    setIsCreating(true)

    try {
      const signer = await provider.getSigner()
      const crowdFunding = getCrowdFundingContract(signer)

      const goalInWei = parseEther(goalEth)
      const durationInSeconds = Number(durationDays) * SECONDS_PER_DAY

      const tx = await crowdFunding.createCampaign(title, description, goalInWei, durationInSeconds)
      await tx.wait()

      await refreshCampaigns()
      setShowCreateModal(false)
      showToast('Your campaign is live! People can now contribute to it.')
    } catch (err) {
      setError(err.shortMessage || err.message)
      throw err
    } finally {
      setIsCreating(false)
    }
  }

  async function handleContribute(campaignId, amountEth) {
    setError(null)
    setIsContributing(true)

    try {
      const signer = await provider.getSigner()
      const crowdFunding = getCrowdFundingContract(signer)

      const amountInWei = parseEther(amountEth)

      const tx = await crowdFunding.contribute(campaignId, { value: amountInWei })
      await tx.wait()

      await refreshCampaigns()
    } catch (err) {
      setError(err.shortMessage || err.message)
      throw err
    } finally {
      setIsContributing(false)
    }
  }

  function handleSelectCampaign(campaignId) {
    const campaign = campaigns.find((c) => c.id === campaignId)
    const isOwner = account && campaign && account.toLowerCase() === campaign.owner.toLowerCase()

    if (isOwner) {
      navigate(`/campaigns/${campaignId}/manage`)
    } else {
      setSelectedCampaignId(campaignId)
    }
  }

  const selectedCampaign = campaigns.find((campaign) => campaign.id === selectedCampaignId) ?? null

  return (
    <div className="flex flex-col gap-5">
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-semibold text-slate-900">Campaigns ({totalCampaigns})</h1>
        <Button
          onClick={() => setShowCreateModal(true)}
          disabled={!account}
          title={!account ? 'Connect your wallet first' : undefined}
        >
          Create Campaign
        </Button>
      </div>

      <CampaignTable campaigns={campaigns} onSelect={handleSelectCampaign} />

      <Pagination
        offset={offset}
        pageSize={PAGE_SIZE}
        total={totalCampaigns}
        onPrevious={() => refreshCampaigns(Math.max(0, offset - PAGE_SIZE))}
        onNext={() => refreshCampaigns(offset + PAGE_SIZE)}
      />

      {showCreateModal && (
        <Modal title="Create a Campaign" onClose={() => setShowCreateModal(false)}>
          <CreateCampaignForm onCreate={handleCreateCampaign} isCreating={isCreating} />
        </Modal>
      )}

      {selectedCampaign && (
        <CampaignDetailsModal
          campaign={selectedCampaign}
          account={account}
          onConnectWallet={onConnectWallet}
          onContribute={handleContribute}
          isContributing={isContributing}
          onClose={() => setSelectedCampaignId(null)}
        />
      )}
    </div>
  )
}

export default CampaignsPage
