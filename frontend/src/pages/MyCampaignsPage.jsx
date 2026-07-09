import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth0 } from '@auth0/auth0-react'
import { fetchMyCampaigns } from '../lib/api'
import CampaignGrid from '../components/CampaignGrid'
import Pagination from '../components/Pagination'
import Button from '../components/ui/Button'

const PAGE_SIZE = 9

function MyCampaignsPage() {
  const navigate = useNavigate()
  const { getAccessTokenSilently } = useAuth0()
  const [campaigns, setCampaigns] = useState([])
  const [total, setTotal] = useState(0)
  const [offset, setOffset] = useState(0)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    refreshCampaigns(0)
  }, [])

  async function refreshCampaigns(targetOffset = offset) {
    setIsLoading(true)
    setError(null)
    try {
      const accessToken = await getAccessTokenSilently()
      const { items, total: totalCount } = await fetchMyCampaigns(accessToken, {
        offset: targetOffset,
        limit: PAGE_SIZE,
      })
      setCampaigns(items)
      setTotal(totalCount)
      setOffset(targetOffset)
    } catch (err) {
      setError(err.message)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="flex flex-col gap-5">
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-semibold text-slate-900">My Campaigns</h1>
        <Button onClick={() => navigate('/create-campaign')}>Create Campaign</Button>
      </div>

      {error && <p className="text-sm text-rose-600">{error}</p>}

      {isLoading ? (
        <p className="text-sm text-slate-500">Loading your campaigns...</p>
      ) : (
        <>
          <CampaignGrid campaigns={campaigns} onSelect={() => {}} showOwner={false} />

          <Pagination
            offset={offset}
            pageSize={PAGE_SIZE}
            total={total}
            onPrevious={() => refreshCampaigns(Math.max(0, offset - PAGE_SIZE))}
            onNext={() => refreshCampaigns(offset + PAGE_SIZE)}
          />
        </>
      )}
    </div>
  )
}

export default MyCampaignsPage
