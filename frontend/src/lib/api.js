const API_BASE_URL = import.meta.env.VITE_API_BASE_URL
const API_V1_URL = `${API_BASE_URL}/api/v1`

export async function fetchCampaigns({ offset = 0, limit = 20 } = {}) {
  const [campaignsResponse, countResponse] = await Promise.all([
    fetch(`${API_V1_URL}/campaigns?offset=${offset}&limit=${limit}`),
    fetch(`${API_V1_URL}/campaigns/count`),
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

export async function fetchSignInMessage(address) {
  const response = await fetch(`${API_V1_URL}/auth/nonce?address=${address}`)
  if (!response.ok) {
    throw new Error(`Failed to get sign-in message (status ${response.status})`)
  }
  return response.json()
}

export async function verifySignIn({ address, signature }) {
  const response = await fetch(`${API_V1_URL}/auth/verify`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ address, signature }),
  })
  if (!response.ok) {
    const { error } = await response.json().catch(() => ({}))
    throw new Error(error || `Sign-in failed (status ${response.status})`)
  }
  return response.json()
}

export async function fetchMe(token) {
  const response = await fetch(`${API_V1_URL}/me`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!response.ok) {
    throw new Error(`Failed to load session (status ${response.status})`)
  }
  return response.json()
}

export async function fetchMyProfile(token) {
  const response = await fetch(`${API_V1_URL}/me/profile`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!response.ok) {
    throw new Error(`Failed to load profile (status ${response.status})`)
  }
  return response.json()
}

export async function updateMyProfile(token, { displayName, email }) {
  const response = await fetch(`${API_V1_URL}/me/profile`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
    body: JSON.stringify({ displayName, email }),
  })
  if (!response.ok) {
    throw new Error(`Failed to update profile (status ${response.status})`)
  }
  return response.json()
}

export async function fetchPublicProfile(address) {
  const response = await fetch(`${API_V1_URL}/profiles/${address}`)
  if (!response.ok) {
    throw new Error(`Failed to load public profile (status ${response.status})`)
  }
  return response.json()
}
