import { formatEther } from 'ethers'

export function shortenAddress(address) {
  return `${address.slice(0, 6)}...${address.slice(-4)}`
}

export function formatEth(amountInWei) {
  return `${formatEther(amountInWei)} ETH`
}

export function formatDate(unixTimestampSeconds) {
  return new Date(Number(unixTimestampSeconds) * 1000).toLocaleString()
}
