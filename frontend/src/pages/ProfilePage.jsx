import { useEffect, useState } from 'react'
import { fetchMyProfile, updateMyProfile } from '../lib/api'
import ProfileForm from '../components/ProfileForm'
import Button from '../components/ui/Button'

function ProfilePage({ account, sessionToken, sessionAddress, isSigningIn, onSignIn }) {
  const [profile, setProfile] = useState(null)
  const [isSaving, setIsSaving] = useState(false)
  const [error, setError] = useState(null)
  const [savedMessage, setSavedMessage] = useState(false)
  const [copied, setCopied] = useState(false)

  function handleCopyAddress() {
    navigator.clipboard.writeText(sessionAddress)
    setCopied(true)
    setTimeout(() => setCopied(false), 1500)
  }

  useEffect(() => {
    if (!sessionToken) return

    fetchMyProfile(sessionToken)
      .then(setProfile)
      .catch((err) => setError(err.message))
  }, [sessionToken])

  async function handleSave({ displayName, email }) {
    setError(null)
    setSavedMessage(false)
    setIsSaving(true)

    try {
      const updated = await updateMyProfile(sessionToken, { displayName, email })
      setProfile(updated)
      setSavedMessage(true)
    } catch (err) {
      setError(err.message)
    } finally {
      setIsSaving(false)
    }
  }

  if (!account) {
    return (
      <div className="mx-auto max-w-md">
        <h1 className="text-xl font-semibold text-slate-900">Profile</h1>
        <p className="mt-3 text-sm text-slate-600">Connect your wallet to view your profile.</p>
      </div>
    )
  }

  if (!sessionAddress) {
    return (
      <div className="mx-auto max-w-md">
        <h1 className="text-xl font-semibold text-slate-900">Profile</h1>
        <p className="mt-3 text-sm text-slate-600">Sign in to view and edit your profile.</p>
        <Button className="mt-4" onClick={onSignIn} disabled={isSigningIn}>
          {isSigningIn ? 'Signing in...' : 'Sign In'}
        </Button>
      </div>
    )
  }

  return (
    <div className="mx-auto max-w-md">
      <h1 className="text-xl font-semibold text-slate-900">Profile</h1>
      <div className="mt-2 mb-6 flex items-center gap-2">
        <p className="break-all font-mono text-sm text-slate-500">{sessionAddress}</p>
        <button
          type="button"
          onClick={handleCopyAddress}
          className="shrink-0 text-xs font-medium text-indigo-600 hover:text-indigo-500"
        >
          {copied ? 'Copied!' : 'Copy'}
        </button>
      </div>

      <div className="rounded-xl border border-slate-200 bg-white p-6 shadow-sm">
        {profile && (
          <ProfileForm
            initialDisplayName={profile.displayName}
            initialEmail={profile.email ?? ''}
            onSave={handleSave}
            isSaving={isSaving}
          />
        )}

        {savedMessage && <p className="mt-3 text-sm text-emerald-600">Profile saved.</p>}
        {error && <p className="mt-3 text-sm text-rose-600">{error}</p>}
      </div>
    </div>
  )
}

export default ProfilePage
