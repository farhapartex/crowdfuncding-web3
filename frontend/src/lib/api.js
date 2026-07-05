const API_BASE_URL = import.meta.env.VITE_API_BASE_URL

export async function fetchCampaigns() {
  const response = await fetch(`${API_BASE_URL}/campaigns`)
  if (!response.ok) {
    throw new Error(`Failed to load campaigns (status ${response.status})`)
  }
  return response.json()
}
