import { useCallback, useState } from 'react'
import { parseEther, parseUnits, ZeroAddress } from 'ethers'
import { getCrowdFundingContract } from '../lib/crowdFundingContract'
import { publishMyCampaign } from '../lib/api'
import { useAccessToken } from './useAccessToken'

const SECONDS_PER_DAY = 24 * 60 * 60

const CURRENCY_MODE_TO_ONCHAIN = { eth: 0, token: 1, both: 2 }

export function usePublishCampaign({ campaign, provider, account, onConnectWallet, onPublished }) {
  const getAccessToken = useAccessToken()
  const [phase, setPhase] = useState('idle')
  const [error, setError] = useState(null)
  const [pendingLink, setPendingLink] = useState(null)

  const finishLinking = useCallback(
    async (linkData) => {
      setPhase('linking')
      setError(null)

      try {
        const accessToken = await getAccessToken()
        await publishMyCampaign(accessToken, campaign.id, linkData)
        setPendingLink(null)
        setPhase('idle')
        onPublished?.()
      } catch (err) {
        setError(err.message)
        setPendingLink(linkData)
        setPhase('error')
      }
    },
    [campaign, getAccessToken, onPublished],
  )

  const publish = useCallback(async () => {
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

    setPhase('signing')
    try {
      const signer = await activeProvider.getSigner()
      const crowdFunding = getCrowdFundingContract(signer)

      const currencyMode = CURRENCY_MODE_TO_ONCHAIN[campaign.currencyMode]
      const isEthMode = campaign.currencyMode === 'eth' || campaign.currencyMode === 'both'
      const isTokenMode = campaign.currencyMode === 'token' || campaign.currencyMode === 'both'

      const goalEth = isEthMode ? parseEther(campaign.targetEth) : 0n
      const goalToken = isTokenMode ? parseUnits(campaign.goalToken, campaign.tokenDecimals ?? 18) : 0n
      const tokenAddress = isTokenMode ? campaign.tokenAddress : ZeroAddress
      const durationInSeconds = Number(campaign.durationDays) * SECONDS_PER_DAY

      const tx = await crowdFunding.createCampaign(
        campaign.title,
        '',
        currencyMode,
        tokenAddress,
        goalEth,
        goalToken,
        durationInSeconds,
      )

      setPhase('confirming')
      const receipt = await tx.wait()

      let onChainCampaignId = null
      for (const log of receipt.logs) {
        try {
          const parsed = crowdFunding.interface.parseLog(log)
          if (parsed?.name === 'CampaignCreated') {
            onChainCampaignId = Number(parsed.args.campaignId)
            break
          }
        } catch {
          continue
        }
      }

      if (onChainCampaignId === null) {
        throw new Error('Published on-chain, but could not read the new campaign id from the transaction.')
      }

      await finishLinking({ walletAddress: activeAccount, onChainCampaignId, txHash: tx.hash })
    } catch (err) {
      setError(err.shortMessage || err.message)
      setPhase('error')
    }
  }, [account, provider, onConnectWallet, campaign, finishLinking])

  const retryLinking = useCallback(() => {
    if (pendingLink) finishLinking(pendingLink)
  }, [pendingLink, finishLinking])

  return { phase, error, pendingLink, publish, retryLinking }
}
