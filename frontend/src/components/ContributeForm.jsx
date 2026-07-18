import { useState } from 'react'
import Field from './ui/Field'
import Button from './ui/Button'

function ContributeForm({ currencyMode, tokenSymbol, onContribute, isContributing }) {
  const [amount, setAmount] = useState('')
  const [currency, setCurrency] = useState(currencyMode === 'token' ? 'token' : 'eth')

  async function handleSubmit(e) {
    e.preventDefault()

    try {
      await onContribute(amount, currency)
      setAmount('')
    } catch {
      // parent already surfaces the error; keep the field filled in so the user can retry
    }
  }

  const unitLabel = currency === 'token' ? tokenSymbol || 'Token' : 'ETH'

  return (
    <form className="flex flex-col gap-3 border-t border-slate-200 pt-4" onSubmit={handleSubmit}>
      {currencyMode === 'both' && (
        <div className="grid grid-cols-2 gap-2">
          <button
            type="button"
            onClick={() => setCurrency('eth')}
            className={`cursor-pointer rounded-lg border px-3 py-2 text-sm font-medium transition-colors ${
              currency === 'eth'
                ? 'border-indigo-500 bg-indigo-50 text-indigo-600'
                : 'border-slate-200 text-slate-600 hover:border-indigo-200'
            }`}
          >
            ETH
          </button>
          <button
            type="button"
            onClick={() => setCurrency('token')}
            className={`cursor-pointer rounded-lg border px-3 py-2 text-sm font-medium transition-colors ${
              currency === 'token'
                ? 'border-indigo-500 bg-indigo-50 text-indigo-600'
                : 'border-slate-200 text-slate-600 hover:border-indigo-200'
            }`}
          >
            {tokenSymbol || 'Token'}
          </button>
        </div>
      )}

      <Field
        id="contribution-amount"
        label={`Contribution amount (${unitLabel})`}
        value={amount}
        onChange={(e) => setAmount(e.target.value)}
        required
      />
      <Button type="submit" disabled={isContributing}>
        {isContributing ? 'Contributing...' : 'Contribute'}
      </Button>
    </form>
  )
}

export default ContributeForm
