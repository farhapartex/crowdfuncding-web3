import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { fetchCampaigns } from '../lib/api'
import CampaignGrid from '../components/CampaignGrid'
import Pagination from '../components/Pagination'

const PAGE_SIZE = 10

function CampaignsPage() {
  const navigate = useNavigate()
  const [campaigns, setCampaigns] = useState([])
  const [totalCampaigns, setTotalCampaigns] = useState(0)
  const [offset, setOffset] = useState(0)

  useEffect(() => {
    refreshCampaigns(0)
  }, [])

  async function refreshCampaigns(targetOffset = offset) {
    const { campaigns: result, total } = await fetchCampaigns({ offset: targetOffset, limit: PAGE_SIZE })
    setCampaigns(result)
    setTotalCampaigns(total)
    setOffset(targetOffset)
  }

  function handleSelectCampaign(campaignId) {
    navigate(`/campaigns/${campaignId}`)
  }

  return (
    <div className="flex flex-col gap-5">
      <h1 className="text-xl font-semibold text-slate-900">Campaigns ({totalCampaigns})</h1>

      <CampaignGrid campaigns={campaigns} onSelect={handleSelectCampaign} />

      <Pagination
        offset={offset}
        pageSize={PAGE_SIZE}
        total={totalCampaigns}
        onPrevious={() => refreshCampaigns(Math.max(0, offset - PAGE_SIZE))}
        onNext={() => refreshCampaigns(offset + PAGE_SIZE)}
      />
    </div>
  )
}

export default CampaignsPage
