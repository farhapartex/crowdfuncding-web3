import { useState } from 'react'
import { BrowserProvider, Contract } from 'ethers'
import crowdFundingAbi from './contract/CrowdFundingAbi.json'
import './App.css'

const CONTRACT_ADDRESS = import.meta.env.VITE_CONTRACT_ADDRESS

function App() {
  const [account, setAccount] = useState(null)
  const [campaignCount, setCampaignCount] = useState(null)
  const [error, setError] = useState(null)

  async function connectWallet() {
    if (!window.ethereum) {
      setError('MetaMask is not installed')
      return
    }

    try {
      const provider = new BrowserProvider(window.ethereum)
      const accounts = await provider.send('eth_requestAccounts', [])
      setAccount(accounts[0])

      const crowdFunding = new Contract(CONTRACT_ADDRESS, crowdFundingAbi, provider)
      const count = await crowdFunding.campaignCount()
      setCampaignCount(count.toString())

      setError(null)
    } catch (err) {
      setError(err.message)
    }
  }

  return (
    <div>
      <h1>Crowd Funding</h1>

      {account ? (
        <div>
          <p>Connected account: {account}</p>
          <p>Campaign count: {campaignCount}</p>
        </div>
      ) : (
        <button onClick={connectWallet}>Connect Wallet</button>
      )}

      {error && <p style={{ color: 'red' }}>{error}</p>}
    </div>
  )
}

export default App
