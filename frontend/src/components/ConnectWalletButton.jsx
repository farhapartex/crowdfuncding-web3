function ConnectWalletButton({ onConnect }) {
  return (
    <div className="connect">
      <button onClick={onConnect}>Connect Wallet</button>
    </div>
  )
}

export default ConnectWalletButton
