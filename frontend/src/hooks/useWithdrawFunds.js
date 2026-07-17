import { useCallback, useState } from 'react'
import { getCrowdFundingContract } from '../lib/crowdFundingContract'

export function useWithdrawFunds({ campaign, provider, account, onConnectWallet, onWithdrawn }) {
  const [phase, setPhase] = useState('idle')
  const [error, setError] = useState(null)

  const withdraw = useCallback(async () => {
    setError(null)

    let activeProvider = provider
    let activeAccount = account

    if (!activeAccount) {
      setPhase('connecting')
      const connected = await onConnectWallet()
      if (!connected) {
        setPhase('idle')
        return
      }
      activeProvider = connected.provider
      activeAccount = connected.account
    }

    if (activeAccount.toLowerCase() !== campaign.walletAddress.toLowerCase()) {
      setError('Connected wallet does not match the campaign owner wallet.')
      setPhase('error')
      return
    }

    setPhase('signing')
    try {
      const signer = await activeProvider.getSigner()
      const crowdFunding = getCrowdFundingContract(signer)

      const tx = await crowdFunding.withdraw(campaign.onChainCampaignId)

      setPhase('confirming')
      await tx.wait()

      setPhase('idle')
      onWithdrawn?.()
    } catch (err) {
      setError(err.shortMessage || err.message)
      setPhase('error')
    }
  }, [account, provider, onConnectWallet, campaign, onWithdrawn])

  return { phase, error, withdraw }
}
