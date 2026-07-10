import { useState } from 'react'
import { useAuth0 } from '@auth0/auth0-react'
import { useCurrentUser } from '../auth/CurrentUserContext'
import { shortenAddress } from '../utils/format'
import Button from '../components/ui/Button'
import TabButton from '../components/ui/TabButton'
import UserTransactionsTab from '../components/UserTransactionsTab'

function getWalletProviderName() {
  if (typeof window === 'undefined' || !window.ethereum) return 'Wallet'
  if (window.ethereum.isMetaMask) return 'MetaMask'
  return 'Wallet'
}

function ProfilePage({ account, onConnectWallet }) {
  const { user } = useAuth0()
  const { currentUser } = useCurrentUser()
  const [copied, setCopied] = useState(false)
  const [activeTab, setActiveTab] = useState('profile')

  function handleCopyAddress() {
    navigator.clipboard.writeText(account)
    setCopied(true)
    setTimeout(() => setCopied(false), 1500)
  }

  const accountLabel = currentUser?.displayName || currentUser?.email || user?.name || user?.email || 'Your account'
  const accountEmail = currentUser?.email || user?.email
  const accountInitial = accountLabel.charAt(0).toUpperCase()

  return (
    <div className="mx-auto flex max-w-2xl flex-col gap-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight text-slate-900">Profile</h1>
        <p className="mt-1 text-sm text-slate-500">Manage your account details and connected wallet.</p>
      </div>

      <div className="flex gap-6 border-b border-slate-200">
        <TabButton active={activeTab === 'profile'} onClick={() => setActiveTab('profile')}>
          Profile
        </TabButton>
        <TabButton active={activeTab === 'transactions'} onClick={() => setActiveTab('transactions')}>
          Transactions
        </TabButton>
      </div>

      {activeTab === 'profile' ? (
        <>
          <div className="flex items-center gap-4 rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
            {user?.picture ? (
              <img src={user.picture} alt={accountLabel} className="h-16 w-16 rounded-full object-cover" />
            ) : (
              <div className="flex h-16 w-16 shrink-0 items-center justify-center rounded-full bg-indigo-600 text-xl font-semibold text-white">
                {accountInitial}
              </div>
            )}
            <div className="min-w-0">
              <p className="truncate text-lg font-semibold text-slate-900">{accountLabel}</p>
              {accountEmail && <p className="truncate text-sm text-slate-500">{accountEmail}</p>}
            </div>
          </div>

          <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
            <h2 className="text-sm font-semibold text-slate-900">Wallet</h2>

            {account ? (
              <div className="mt-3 flex flex-col gap-3 rounded-xl bg-slate-50 px-4 py-3">
                <div className="flex items-center justify-between gap-3">
                  <div className="min-w-0">
                    <p className="text-sm font-medium text-slate-900">{getWalletProviderName()}</p>
                    <p className="truncate font-mono text-sm text-slate-500">{shortenAddress(account)}</p>
                  </div>
                  <button
                    type="button"
                    onClick={handleCopyAddress}
                    className="shrink-0 cursor-pointer text-xs font-medium text-indigo-600 hover:text-indigo-500"
                  >
                    {copied ? 'Copied!' : 'Copy'}
                  </button>
                </div>

                <div className="flex justify-end">
                  <Button variant="secondary" onClick={() => {}}>
                    Disconnect Wallet
                  </Button>
                </div>
              </div>
            ) : (
              <div className="mt-3 flex items-center justify-between gap-3 rounded-xl bg-slate-50 px-4 py-3">
                <p className="text-sm text-slate-500">No wallet connected yet.</p>
                <Button onClick={onConnectWallet}>Connect Wallet</Button>
              </div>
            )}
          </div>
        </>
      ) : (
        <UserTransactionsTab />
      )}
    </div>
  )
}

export default ProfilePage
