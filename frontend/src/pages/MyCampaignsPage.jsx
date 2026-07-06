import { useNavigate } from 'react-router-dom'
import CampaignGrid from '../components/CampaignGrid'
import Button from '../components/ui/Button'

const PLACEHOLDER_CAMPAIGNS = [
  {
    id: 0,
    title: 'Help Build a Community Garden',
    amountRaised: '2500000000000000000',
  },
  {
    id: 1,
    title: 'New Laptop for Freelance Work',
    amountRaised: '1000000000000000000',
  },
  {
    id: 2,
    title: 'Local Animal Shelter Renovation',
    amountRaised: '4200000000000000000',
  },
]

function MyCampaignsPage() {
  const navigate = useNavigate()

  return (
    <div className="flex flex-col gap-5">
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-semibold text-slate-900">My Campaigns</h1>
        <Button onClick={() => navigate('/create-campaign')}>Create Campaign</Button>
      </div>

      <CampaignGrid campaigns={PLACEHOLDER_CAMPAIGNS} onSelect={() => {}} />
    </div>
  )
}

export default MyCampaignsPage
