import { Contract } from 'ethers'
import crowdFundingAbi from '../contract/CrowdFundingAbi.json'

export const CONTRACT_ADDRESS = import.meta.env.VITE_CONTRACT_ADDRESS

export function getCrowdFundingContract(providerOrSigner) {
  return new Contract(CONTRACT_ADDRESS, crowdFundingAbi, providerOrSigner)
}
