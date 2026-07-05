function WithdrawButton({ onWithdraw, isWithdrawing }) {
  return (
    <div className="withdraw-section">
      <button type="button" onClick={onWithdraw} disabled={isWithdrawing}>
        {isWithdrawing ? 'Withdrawing...' : 'Withdraw Funds'}
      </button>
    </div>
  )
}

export default WithdrawButton
