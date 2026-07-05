import { useState } from 'react'

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
    <form className="campaign-form" onSubmit={handleSubmit}>
      <div className="field">
        <label htmlFor="title">Title</label>
        <input id="title" value={title} onChange={(e) => setTitle(e.target.value)} required />
      </div>
      <div className="field">
        <label htmlFor="description">Description</label>
        <input
          id="description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          required
        />
      </div>
      <div className="field">
        <label htmlFor="goal">Goal (ETH)</label>
        <input id="goal" value={goalEth} onChange={(e) => setGoalEth(e.target.value)} required />
      </div>
      <div className="field">
        <label htmlFor="duration">Duration (days)</label>
        <input
          id="duration"
          value={durationDays}
          onChange={(e) => setDurationDays(e.target.value)}
          required
        />
      </div>
      <button type="submit" disabled={isCreating}>
        {isCreating ? 'Creating...' : 'Create Campaign'}
      </button>
    </form>
  )
}

export default CreateCampaignForm
