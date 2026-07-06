import Button from './ui/Button'

function WithdrawButton({ onWithdraw, isWithdrawing }) {
  return (
    <div className="border-t border-slate-200 pt-4">
      <Button type="button" onClick={onWithdraw} disabled={isWithdrawing}>
        {isWithdrawing ? 'Withdrawing...' : 'Withdraw Funds'}
      </Button>
    </div>
  )
}

export default WithdrawButton
