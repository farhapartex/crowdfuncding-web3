import { useEffect, useRef, useState } from 'react'
import { Routes, Route } from 'react-router-dom'
import { BrowserProvider } from 'ethers'
import { fetchSignInMessage, verifySignIn, fetchMe } from './lib/api'
import Navbar from './components/Navbar'
import ToastContainer from './components/ToastContainer'
import CampaignsPage from './pages/CampaignsPage'
import AboutPage from './pages/AboutPage'
import ProfilePage from './pages/ProfilePage'

const SESSION_TOKEN_KEY = 'sessionToken'
const TOAST_DURATION_MS = 4000

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
      <Navbar
        account={account}
        sessionAddress={sessionAddress}
        isSigningIn={isSigningIn}
        onConnect={connectWallet}
        onSignIn={handleSignIn}
        onSignOut={handleSignOut}
      />

      <main className="mx-auto max-w-5xl px-6 py-8">
        {error && (
          <p className="mb-4 rounded-lg border border-rose-200 bg-rose-50 px-4 py-2.5 text-sm text-rose-600">
            {error}
          </p>
        )}

        <Routes>
          <Route
            path="/"
            element={
              <CampaignsPage
                account={account}
                provider={provider}
                onConnectWallet={connectWallet}
                setError={setError}
                showToast={showToast}
              />
            }
          />
          <Route path="/about" element={<AboutPage />} />
          <Route
            path="/profile"
            element={
              <ProfilePage
                account={account}
                sessionToken={sessionToken}
                sessionAddress={sessionAddress}
                isSigningIn={isSigningIn}
                onSignIn={handleSignIn}
              />
            }
          />
        </Routes>
      </main>

      <ToastContainer toasts={toasts} onDismiss={dismissToast} />
    </div>
  )
}

export default App
