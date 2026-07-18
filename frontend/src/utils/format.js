import { formatEther, formatUnits } from 'ethers'

export function shortenAddress(address) {
  return `${address.slice(0, 6)}...${address.slice(-4)}`
}

export function formatEth(amountInWei) {
  return `${formatEther(amountInWei)} ETH`
}

export function formatTokenAmount(amountInSmallestUnits, decimals, symbol) {
  const formatted = Number(formatUnits(amountInSmallestUnits || '0', decimals ?? 18)).toLocaleString(undefined, {
    maximumFractionDigits: 4,
  })
  return symbol ? `${formatted} ${symbol}` : formatted
}

export function formatDate(unixTimestampSeconds) {
  return new Date(Number(unixTimestampSeconds) * 1000).toLocaleString()
}

export function formatEthDisplay(value) {
  const parsed = Number(value)
  if (!value || Number.isNaN(parsed)) return '0'
  return parsed.toLocaleString(undefined, { maximumFractionDigits: 4 })
}
