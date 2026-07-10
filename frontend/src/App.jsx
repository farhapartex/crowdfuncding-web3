import { useEffect, useRef, useState } from 'react'
import { Routes, Route } from 'react-router-dom'
import { BrowserProvider } from 'ethers'
import { fetchSignInMessage, verifySignIn, fetchMe } from './lib/api'
import Navbar from './components/Navbar'
import ToastContainer from './components/ToastContainer'
import CampaignsPage from './pages/CampaignsPage'
import CampaignDetailsPage from './pages/CampaignDetailsPage'
import CampaignManagePage from './pages/CampaignManagePage'
import ProfilePage from './pages/ProfilePage'
import MyCampaignsPage from './pages/MyCampaignsPage'
import MyCampaignDetailsPage from './pages/MyCampaignDetailsPage'
import CreateCampaignPage from './pages/CreateCampaignPage'
import HomePage from './landing/HomePage'
import AboutUsPage from './landing/AboutUsPage'
import withAuthGuard from './auth/withAuthGuard'

const SESSION_TOKEN_KEY = 'sessionToken'
const TOAST_DURATION_MS = 4000

const GuardedMyCampaignsPage = withAuthGuard(MyCampaignsPage)
const GuardedMyCampaignDetailsPage = withAuthGuard(MyCampaignDetailsPage)
const GuardedCreateCampaignPage = withAuthGuard(CreateCampaignPage)
const GuardedProfilePage = withAuthGuard(ProfilePage)

function App() {
  const [provider, setProvider] = useState(null)
  const [account, setAccount] = useState(null)
  const [error, setError] = useState(null)
  const [sessionToken, setSessionToken] = useState(null)
  const [sessionAddress, setSessionAddress] = useState(null)
  const [isSigningIn, setIsSigningIn] = useState(false)
  const [toasts, setToasts] = useState([])
  const nextToastId = useRef(0)

  useEffect(() => {
    const storedToken = localStorage.getItem(SESSION_TOKEN_KEY)
    if (!storedToken) return

    fetchMe(storedToken)
      .then(({ address }) => {
        setSessionToken(storedToken)
        setSessionAddress(address)
      })
      .catch(() => localStorage.removeItem(SESSION_TOKEN_KEY))
  }, [])

  useEffect(() => {
    if (!window.ethereum) return

    const browserProvider = new BrowserProvider(window.ethereum)
    browserProvider.send('eth_accounts', []).then((accounts) => {
      if (accounts.length > 0) {
        setProvider(browserProvider)
        setAccount(accounts[0])
      }
    })
  }, [])

  async function connectWallet() {
    if (!window.ethereum) {
      setError('MetaMask is not installed')
      return null
    }

    try {
      const browserProvider = new BrowserProvider(window.ethereum)
      const accounts = await browserProvider.send('eth_requestAccounts', [])

      setProvider(browserProvider)
      setAccount(accounts[0])
      setError(null)

      return { provider: browserProvider, account: accounts[0] }
    } catch (err) {
      setError(err.message)
      return null
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

  function dismissToast(id) {
    setToasts((prev) => prev.filter((toast) => toast.id !== id))
  }

  function showToast(message) {
    const id = ++nextToastId.current
    setToasts((prev) => [...prev, { id, message }])
    setTimeout(() => dismissToast(id), TOAST_DURATION_MS)
  }

  return (
    <div className="min-h-screen bg-slate-50">
      <Navbar account={account} onConnect={connectWallet} />

      <main className="mx-auto max-w-5xl px-6 py-8">
        {error && (
          <p className="mb-4 rounded-lg border border-rose-200 bg-rose-50 px-4 py-2.5 text-sm text-rose-600">
            {error}
          </p>
        )}

        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/campaigns" element={<CampaignsPage />} />
          <Route
            path="/campaigns/:id"
            element={
              <CampaignDetailsPage
                provider={provider}
                account={account}
                onConnectWallet={connectWallet}
                setError={setError}
                showToast={showToast}
              />
            }
          />
          <Route
            path="/campaigns/:id/manage"
            element={
              <CampaignManagePage
                provider={provider}
                account={account}
                sessionAddress={sessionAddress}
                setError={setError}
                showToast={showToast}
              />
            }
          />
          <Route path="/about" element={<AboutUsPage />} />
          <Route path="/my-campaigns" element={<GuardedMyCampaignsPage />} />
          <Route
            path="/my-campaigns/:id"
            element={
              <GuardedMyCampaignDetailsPage
                provider={provider}
                account={account}
                onConnectWallet={connectWallet}
                showToast={showToast}
              />
            }
          />
          <Route
            path="/create-campaign"
            element={
              <GuardedCreateCampaignPage
                provider={provider}
                account={account}
                onConnectWallet={connectWallet}
                showToast={showToast}
              />
            }
          />
          <Route
            path="/profile"
            element={
              <GuardedProfilePage account={account} onConnectWallet={connectWallet} />
            }
          />
        </Routes>
      </main>

      <ToastContainer toasts={toasts} onDismiss={dismissToast} />
    </div>
  )
}

export default App
