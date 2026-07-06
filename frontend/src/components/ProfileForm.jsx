import { useState } from 'react'
import Field from './ui/Field'
import Button from './ui/Button'

function ProfileForm({ initialDisplayName, initialEmail, onSave, isSaving }) {
  const [displayName, setDisplayName] = useState(initialDisplayName)

  async function handleSubmit(e) {
    e.preventDefault()
    await onSave({ displayName, email: initialEmail })
  }

  return (
    <form className="flex flex-col gap-4" onSubmit={handleSubmit}>
      <Field
        id="displayName"
        label="Display name"
        value={displayName}
        onChange={(e) => setDisplayName(e.target.value)}
      />
      <Field id="profileEmail" label="Email" type="email" value={initialEmail} placeholder="Not set" disabled />
      <Button type="submit" disabled={isSaving} className="mt-1">
        {isSaving ? 'Saving...' : 'Save Profile'}
      </Button>
    </form>
  )
}

export default ProfileForm
