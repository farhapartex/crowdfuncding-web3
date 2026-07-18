import { Contract } from 'ethers'

const ERC20_ABI = [
  'function decimals() view returns (uint8)',
  'function balanceOf(address owner) view returns (uint256)',
  'function allowance(address owner, address spender) view returns (uint256)',
  'function approve(address spender, uint256 amount) returns (bool)',
]

export function getErc20Contract(tokenAddress, providerOrSigner) {
  return new Contract(tokenAddress, ERC20_ABI, providerOrSigner)
}
