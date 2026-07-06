import { useState } from 'react'
import { BrowserProvider, parseEther } from 'ethers'
import { getCrowdFundingContract } from './lib/crowdFundingContract'
import { fetchCampaigns } from './lib/api'
import ConnectWalletButton from './components/ConnectWalletButton'
import SummaryCard from './components/SummaryCard'
import CreateCampaignForm from './components/CreateCampaignForm'
import CampaignTable from './components/CampaignTable'
import Pagination from './components/Pagination'
import Modal from './components/Modal'
import CampaignDetailsModal from './components/CampaignDetailsModal'
import './App.css'

const SECONDS_PER_DAY = 24 * 60 * 60
const PAGE_SIZE = 10

function App() {
  const [provider, setProvider] = useState(null)
  const [account, setAccount] = useState(null)
  const [campaigns, setCampaigns] = useState([])
  const [totalCampaigns, setTotalCampaigns] = useState(0)
  const [offset, setOffset] = useState(0)
  const [selectedCampaignId, setSelectedCampaignId] = useState(null)
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [error, setError] = useState(null)
  const [isCreating, setIsCreating] = useState(false)
  const [isContributing, setIsContributing] = useState(false)
  const [isWithdrawing, setIsWithdrawing] = useState(false)

  async function refreshCampaigns(targetOffset = offset) {
    const { campaigns: result, total } = await fetchCampaigns({ offset: targetOffset, limit: PAGE_SIZE })
    setCampaigns(result)
    setTotalCampaigns(total)
    setOffset(targetOffset)
  }

  async function connectWallet() {
    if (!window.ethereum) {
      setError('MetaMask is not installed')
      return
    }

    try {
      const browserProvider = new BrowserProvider(window.ethereum)
      const accounts = await browserProvider.send('eth_requestAccounts', [])

      setProvider(browserProvider)
      setAccount(accounts[0])
      setError(null)

      await refreshCampaigns(0)
    } catch (err) {
      setError(err.message)
    }
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

  async function handleWithdraw(campaignId) {
    setError(null)
    setIsWithdrawing(true)

    try {
      const signer = await provider.getSigner()
      const crowdFunding = getCrowdFundingContract(signer)

      const tx = await crowdFunding.withdraw(campaignId)
      await tx.wait()

      await refreshCampaigns()
    } catch (err) {
      setError(err.shortMessage || err.message)
    } finally {
      setIsWithdrawing(false)
    }
  }

  const selectedCampaign = campaigns.find((campaign) => campaign.id === selectedCampaignId) ?? null

  return (
    <div className="app">
      <h1>Crowd Funding</h1>

      {account ? (
        <div className="dashboard">
          <SummaryCard account={account} campaignCount={totalCampaigns} />

          <div className="section-header">
            <h2>Campaigns</h2>
            <button onClick={() => setShowCreateModal(true)}>Create Campaign</button>
          </div>

          <CampaignTable campaigns={campaigns} onSelect={setSelectedCampaignId} />

          <Pagination
            offset={offset}
            pageSize={PAGE_SIZE}
            total={totalCampaigns}
            onPrevious={() => refreshCampaigns(Math.max(0, offset - PAGE_SIZE))}
            onNext={() => refreshCampaigns(offset + PAGE_SIZE)}
          />
        </div>
      ) : (
        <ConnectWalletButton onConnect={connectWallet} />
      )}

      {error && <p className="error">{error}</p>}

      {showCreateModal && (
        <Modal title="Create a Campaign" onClose={() => setShowCreateModal(false)}>
          <CreateCampaignForm onCreate={handleCreateCampaign} isCreating={isCreating} />
        </Modal>
      )}

      {selectedCampaign && (
        <CampaignDetailsModal
          campaign={selectedCampaign}
          account={account}
          onContribute={handleContribute}
          isContributing={isContributing}
          onWithdraw={handleWithdraw}
          isWithdrawing={isWithdrawing}
          onClose={() => setSelectedCampaignId(null)}
        />
      )}
    </div>
  )
}

export default App
