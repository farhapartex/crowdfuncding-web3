import { useState } from 'react'

function ProfileForm({ initialDisplayName, initialEmail, onSave, isSaving }) {
  const [displayName, setDisplayName] = useState(initialDisplayName)
  const [email, setEmail] = useState(initialEmail)

  async function handleSubmit(e) {
    e.preventDefault()
    await onSave({ displayName, email })
  }

  return (
    <form className="campaign-form" onSubmit={handleSubmit}>
      <div className="field">
        <label htmlFor="displayName">Display name</label>
        <input id="displayName" value={displayName} onChange={(e) => setDisplayName(e.target.value)} />
      </div>
      <div className="field">
        <label htmlFor="profileEmail">Email</label>
        <input id="profileEmail" type="email" value={email} onChange={(e) => setEmail(e.target.value)} />
      </div>
      <button type="submit" disabled={isSaving}>
        {isSaving ? 'Saving...' : 'Save Profile'}
      </button>
    </form>
  )
}

export default ProfileForm
