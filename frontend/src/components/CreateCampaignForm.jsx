import { useState } from 'react'
import Field from './ui/Field'
import Button from './ui/Button'

function CreateCampaignForm({ onCreate, isCreating }) {
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [goalEth, setGoalEth] = useState('')
  const [durationDays, setDurationDays] = useState('')

  async function handleSubmit(e) {
    e.preventDefault()

    try {
      await onCreate({ title, description, goalEth, durationDays })
      setTitle('')
      setDescription('')
      setGoalEth('')
      setDurationDays('')
    } catch {
      // parent already surfaces the error; keep the form filled in so the user can retry
    }
  }

  return (
    <form className="flex flex-col gap-4" onSubmit={handleSubmit}>
      <Field id="title" label="Title" value={title} onChange={(e) => setTitle(e.target.value)} required />
      <Field
        id="description"
        label="Description"
        value={description}
        onChange={(e) => setDescription(e.target.value)}
        required
      />
      <Field id="goal" label="Goal (ETH)" value={goalEth} onChange={(e) => setGoalEth(e.target.value)} required />
      <Field
        id="duration"
        label="Duration (days)"
        value={durationDays}
        onChange={(e) => setDurationDays(e.target.value)}
        required
      />
      <Button type="submit" disabled={isCreating} className="mt-1">
        {isCreating ? 'Creating...' : 'Create Campaign'}
      </Button>
    </form>
  )
}

export default CreateCampaignForm
