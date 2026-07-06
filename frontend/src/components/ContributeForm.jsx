import { useState } from 'react'
import Field from './ui/Field'
import Button from './ui/Button'

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
    <form className="flex flex-col gap-3 border-t border-slate-200 pt-4" onSubmit={handleSubmit}>
      <Field
        id="contribution-amount"
        label="Contribution amount (ETH)"
        value={amountEth}
        onChange={(e) => setAmountEth(e.target.value)}
        required
      />
      <Button type="submit" disabled={isContributing}>
        {isContributing ? 'Contributing...' : 'Contribute'}
      </Button>
    </form>
  )
}

export default ContributeForm
