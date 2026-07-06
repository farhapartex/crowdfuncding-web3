const API_BASE_URL = import.meta.env.VITE_API_BASE_URL

export async function fetchCampaigns({ offset = 0, limit = 20 } = {}) {
  const [campaignsResponse, countResponse] = await Promise.all([
    fetch(`${API_BASE_URL}/campaigns?offset=${offset}&limit=${limit}`),
    fetch(`${API_BASE_URL}/campaigns/count`),
  ])

  if (!campaignsResponse.ok) {
    throw new Error(`Failed to load campaigns (status ${campaignsResponse.status})`)
  }
  if (!countResponse.ok) {
    throw new Error(`Failed to load campaign count (status ${countResponse.status})`)
  }

  const campaigns = await campaignsResponse.json()
  const { count } = await countResponse.json()

  return { campaigns, total: Number(count) }
}
