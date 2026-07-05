import { useState } from 'react'

function ContributeForm({ onContribute, isContributing }) {
  const [amountEth, setAmountEth] = useState('')

  async function handleSubmit(e) {
    e.preventDefault()

    try {
      await onContribute(amountEth)
      setAmountEth('')
    } catch {
      // parent already surfaces the error; keep the field filled in so the user can retry
    }
  }

  return (
    <form className="contribute-form" onSubmit={handleSubmit}>
      <div className="field">
        <label htmlFor="contribution-amount">Contribution amount (ETH)</label>
        <input
          id="contribution-amount"
          value={amountEth}
          onChange={(e) => setAmountEth(e.target.value)}
          required
        />
      </div>
      <button type="submit" disabled={isContributing}>
        {isContributing ? 'Contributing...' : 'Contribute'}
      </button>
    </form>
  )
}

export default ContributeForm
