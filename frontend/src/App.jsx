import { useEffect, useState } from 'react'
import { BrowserProvider, parseEther } from 'ethers'
import { getCrowdFundingContract } from './lib/crowdFundingContract'
import { fetchCampaigns, fetchSignInMessage, verifySignIn, fetchMe, fetchMyProfile, updateMyProfile } from './lib/api'
import { shortenAddress } from './utils/format'
import ConnectWalletButton from './components/ConnectWalletButton'
import AuthStatus from './components/AuthStatus'
import ProfileForm from './components/ProfileForm'
import CreateCampaignForm from './components/CreateCampaignForm'
import CampaignTable from './components/CampaignTable'
import Pagination from './components/Pagination'
import Modal from './components/Modal'
import CampaignDetailsModal from './components/CampaignDetailsModal'
import './App.css'

const SECONDS_PER_DAY = 24 * 60 * 60
const PAGE_SIZE = 10
const SESSION_TOKEN_KEY = 'sessionToken'

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
  const [sessionToken, setSessionToken] = useState(null)
  const [sessionAddress, setSessionAddress] = useState(null)
  const [isSigningIn, setIsSigningIn] = useState(false)
  const [showProfileModal, setShowProfileModal] = useState(false)
  const [myProfile, setMyProfile] = useState(null)
  const [isSavingProfile, setIsSavingProfile] = useState(false)

  useEffect(() => {
    refreshCampaigns(0)

    const storedToken = localStorage.getItem(SESSION_TOKEN_KEY)
    if (!storedToken) return

    fetchMe(storedToken)
      .then(({ address }) => {
        setSessionToken(storedToken)
        setSessionAddress(address)
      })
      .catch(() => localStorage.removeItem(SESSION_TOKEN_KEY))
  }, [])

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

  async function handleSignIn() {
    setError(null)
    setIsSigningIn(true)

    try {
      const signer = await provider.getSigner()
      const { message } = await fetchSignInMessage(account)
      const signature = await signer.signMessage(message)
      const { token, address } = await verifySignIn({ address: account, signature })

      setSessionToken(token)
      setSessionAddress(address)
      localStorage.setItem(SESSION_TOKEN_KEY, token)
    } catch (err) {
      setError(err.shortMessage || err.message)
    } finally {
      setIsSigningIn(false)
    }
  }

  function handleSignOut() {
    setSessionToken(null)
    setSessionAddress(null)
    localStorage.removeItem(SESSION_TOKEN_KEY)
  }

  async function handleOpenProfile() {
    setError(null)
    try {
      const profile = await fetchMyProfile(sessionToken)
      setMyProfile(profile)
      setShowProfileModal(true)
    } catch (err) {
      setError(err.message)
    }
  }

  async function handleSaveProfile({ displayName, email }) {
    setError(null)
    setIsSavingProfile(true)

    try {
      const profile = await updateMyProfile(sessionToken, { displayName, email })
      setMyProfile(profile)
      setShowProfileModal(false)
    } catch (err) {
      setError(err.message)
    } finally {
      setIsSavingProfile(false)
    }
  }

  const selectedCampaign = campaigns.find((campaign) => campaign.id === selectedCampaignId) ?? null

  return (
    <div className="app">
      <header className="topbar">
        <h1>Crowd Funding</h1>

        <div className="topbar-actions">
          {account ? (
            <>
              <span className="value mono">{shortenAddress(account)}</span>
              <AuthStatus
                sessionAddress={sessionAddress}
                isSigningIn={isSigningIn}
                onSignIn={handleSignIn}
                onSignOut={handleSignOut}
                onEditProfile={handleOpenProfile}
              />
            </>
          ) : (
            <ConnectWalletButton onConnect={connectWallet} />
          )}
        </div>
      </header>

      <div className="dashboard">
        <div className="section-header">
          <h2>Campaigns ({totalCampaigns})</h2>
          <button
            onClick={() => setShowCreateModal(true)}
            disabled={!account}
            title={!account ? 'Connect your wallet first' : undefined}
          >
            Create Campaign
          </button>
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

      {error && <p className="error">{error}</p>}

      {showCreateModal && (
        <Modal title="Create a Campaign" onClose={() => setShowCreateModal(false)}>
          <CreateCampaignForm onCreate={handleCreateCampaign} isCreating={isCreating} />
        </Modal>
      )}

      {showProfileModal && myProfile && (
        <Modal title="Edit Profile" onClose={() => setShowProfileModal(false)}>
          <ProfileForm
            initialDisplayName={myProfile.displayName}
            initialEmail={myProfile.email ?? ''}
            onSave={handleSaveProfile}
            isSaving={isSavingProfile}
          />
        </Modal>
      )}

      {selectedCampaign && (
        <CampaignDetailsModal
          campaign={selectedCampaign}
          account={account}
          onConnectWallet={connectWallet}
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
