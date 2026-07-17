const API_BASE_URL = import.meta.env.VITE_API_BASE_URL
const API_V1_URL = `${API_BASE_URL}/api/v1`

export async function fetchCampaigns({ offset = 0, limit = 20, category = '' } = {}) {
  const params = new URLSearchParams({ offset, limit })
  if (category) params.set('category', category)

  const response = await fetch(`${API_V1_URL}/campaigns?${params.toString()}`)
  if (!response.ok) {
    throw new Error(`Failed to load campaigns (status ${response.status})`)
  }

  const { items, total } = await response.json()
  return { campaigns: items, total }
}

export async function fetchCampaign(id) {
  const response = await fetch(`${API_V1_URL}/campaigns/${id}`)
  if (!response.ok) {
    throw new Error(`Failed to load campaign (status ${response.status})`)
  }
  return response.json()
}

export async function fetchContributors(campaignId) {
  const response = await fetch(`${API_V1_URL}/campaigns/${campaignId}/contributors`)
  if (!response.ok) {
    throw new Error(`Failed to load contributors (status ${response.status})`)
  }
  return response.json()
}

export async function fetchCampaignTransactions(campaignId, { offset = 0, limit = 20 } = {}) {
  const response = await fetch(`${API_V1_URL}/campaigns/${campaignId}/transactions?offset=${offset}&limit=${limit}`)
  if (!response.ok) {
    throw new Error(`Failed to load transactions (status ${response.status})`)
  }
  return response.json()
}

export async function fetchWalletTransactions(address, { offset = 0, limit = 20 } = {}) {
  const response = await fetch(`${API_V1_URL}/wallets/${address}/transactions?offset=${offset}&limit=${limit}`)
  if (!response.ok) {
    throw new Error(`Failed to load transactions (status ${response.status})`)
  }
  return response.json()
}

export async function fetchCampaignComments(campaignId, { offset = 0, limit = 20 } = {}) {
  const response = await fetch(`${API_V1_URL}/campaigns/${campaignId}/comments?offset=${offset}&limit=${limit}`)
  if (!response.ok) {
    throw new Error(`Failed to load comments (status ${response.status})`)
  }
  return response.json()
}

export async function postCampaignComment(accessToken, campaignId, text) {
  const response = await fetch(`${API_V1_URL}/campaigns/${campaignId}/comments`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${accessToken}` },
    body: JSON.stringify({ text }),
  })
  if (!response.ok) {
    const { error } = await response.json().catch(() => ({}))
    throw new Error(error || `Failed to post comment (status ${response.status})`)
  }
  return response.json()
}

export async function postCommentReply(accessToken, commentId, text) {
  const response = await fetch(`${API_V1_URL}/comments/${commentId}/replies`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${accessToken}` },
    body: JSON.stringify({ text }),
  })
  if (!response.ok) {
    const { error } = await response.json().catch(() => ({}))
    throw new Error(error || `Failed to post reply (status ${response.status})`)
  }
  return response.json()
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

export async function fetchAuth0Me(accessToken) {
  const response = await fetch(`${API_V1_URL}/auth0/me`, {
    headers: { Authorization: `Bearer ${accessToken}` },
  })
  if (!response.ok) {
    const error = new Error(`Failed to load current user (status ${response.status})`)
    error.status = response.status
    throw error
  }
  return response.json()
}

export async function syncAuth0User(accessToken) {
  const response = await fetch(`${API_V1_URL}/auth0/sync`, {
    method: 'POST',
    headers: { Authorization: `Bearer ${accessToken}` },
  })
  if (!response.ok) {
    throw new Error(`Failed to sync user (status ${response.status})`)
  }
  return response.json()
}

export async function uploadAsset(accessToken, file) {
  const formData = new FormData()
  formData.append('file', file)

  const response = await fetch(`${API_V1_URL}/assets`, {
    method: 'POST',
    headers: { Authorization: `Bearer ${accessToken}` },
    body: formData,
  })
  if (!response.ok) {
    const { error } = await response.json().catch(() => ({}))
    throw new Error(error || `Failed to upload image (status ${response.status})`)
  }
  return response.json()
}

export async function fetchMyCampaigns(accessToken, { offset = 0, limit = 20 } = {}) {
  const response = await fetch(`${API_V1_URL}/my-campaigns?offset=${offset}&limit=${limit}`, {
    headers: { Authorization: `Bearer ${accessToken}` },
  })
  if (!response.ok) {
    throw new Error(`Failed to load your campaigns (status ${response.status})`)
  }
  return response.json()
}

export async function createMyCampaign(accessToken, payload) {
  const response = await fetch(`${API_V1_URL}/my-campaigns`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${accessToken}` },
    body: JSON.stringify(payload),
  })
  if (!response.ok) {
    const { error } = await response.json().catch(() => ({}))
    throw new Error(error || `Failed to create campaign (status ${response.status})`)
  }
  return response.json()
}

export async function fetchMyCampaign(accessToken, id) {
  const response = await fetch(`${API_V1_URL}/my-campaigns/${id}`, {
    headers: { Authorization: `Bearer ${accessToken}` },
  })
  if (!response.ok) {
    throw new Error(`Failed to load campaign (status ${response.status})`)
  }
  return response.json()
}

export async function publishMyCampaign(accessToken, id, payload) {
  const response = await fetch(`${API_V1_URL}/my-campaigns/${id}/publish`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${accessToken}` },
    body: JSON.stringify(payload),
  })
  if (!response.ok) {
    const { error } = await response.json().catch(() => ({}))
    throw new Error(error || `Failed to publish campaign (status ${response.status})`)
  }
  return response.json()
}

export async function deleteMyCampaign(accessToken, id) {
  const response = await fetch(`${API_V1_URL}/my-campaigns/${id}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${accessToken}` },
  })
  if (!response.ok) {
    const { error } = await response.json().catch(() => ({}))
    throw new Error(error || `Failed to delete campaign (status ${response.status})`)
  }
}

export async function fetchPublicProfile(address) {
  const response = await fetch(`${API_V1_URL}/profiles/${address}`)
  if (!response.ok) {
    throw new Error(`Failed to load public profile (status ${response.status})`)
  }
  return response.json()
}
