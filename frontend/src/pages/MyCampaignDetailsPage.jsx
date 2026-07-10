import { useCallback, useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useAuth0 } from '@auth0/auth0-react'
import { fetchMyCampaign, deleteMyCampaign } from '../lib/api'
import { usePublishCampaign } from '../hooks/usePublishCampaign'
import CampaignPreview from '../components/CampaignPreview'
import CampaignTransactionsTab from '../components/CampaignTransactionsTab'
import ConfirmDialog from '../components/ConfirmDialog'
import TabButton from '../components/ui/TabButton'

const PUBLISH_LABELS = {
  connecting: 'Connecting wallet...',
  signing: 'Confirm in wallet...',
  confirming: 'Waiting for confirmation...',
  linking: 'Finalizing...',
}

function MyCampaignDetailsPage({ provider, account, onConnectWallet, showToast }) {
  const { id } = useParams()
  const navigate = useNavigate()
  const { getAccessTokenSilently } = useAuth0()
  const [campaign, setCampaign] = useState(null)
  const [error, setError] = useState(null)
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false)
  const [isDeleting, setIsDeleting] = useState(false)
  const [activeTab, setActiveTab] = useState('details')

  const loadCampaign = useCallback(() => {
    return getAccessTokenSilently()
      .then((accessToken) => fetchMyCampaign(accessToken, id))
      .then(setCampaign)
      .catch((err) => setError(err.message))
  }, [id, getAccessTokenSilently])

  useEffect(() => {
    loadCampaign()
  }, [loadCampaign])

  const publishHook = usePublishCampaign({
    campaign,
    provider,
    account,
    onConnectWallet,
    onPublished: () => {
      showToast?.('Your campaign is live!')
      loadCampaign()
    },
  })

  async function handleConfirmDelete() {
    setError(null)
    setIsDeleting(true)

    try {
      const accessToken = await getAccessTokenSilently()
      await deleteMyCampaign(accessToken, id)
      navigate('/my-campaigns')
    } catch (err) {
      setError(err.message)
      setIsDeleting(false)
      setShowDeleteConfirm(false)
    }
  }

  if (error) {
    return <p className="text-sm text-rose-600">{error}</p>
  }

  if (!campaign) {
    return <p className="text-sm text-slate-500">Loading campaign...</p>
  }

  const isPublishing = ['connecting', 'signing', 'confirming', 'linking'].includes(publishHook.phase)
  const publishLabel = PUBLISH_LABELS[publishHook.phase] || (publishHook.pendingLink ? 'Finish Publishing' : 'Publish')

  return (
    <>
      <div className="mx-auto mb-6 flex max-w-5xl gap-6 border-b border-slate-200">
        <TabButton active={activeTab === 'details'} onClick={() => setActiveTab('details')}>
          Details
        </TabButton>
        <TabButton active={activeTab === 'transactions'} onClick={() => setActiveTab('transactions')}>
          Transactions
        </TabButton>
      </div>

      {activeTab === 'details' ? (
        <CampaignPreview
          campaign={campaign}
          onBack={() => navigate('/my-campaigns')}
          onPublish={publishHook.pendingLink ? publishHook.retryLinking : publishHook.publish}
          publishLabel={publishLabel}
          isPublishing={isPublishing}
          publishError={publishHook.error}
          onDelete={() => setShowDeleteConfirm(true)}
        />
      ) : (
        <div className="mx-auto max-w-5xl">
          <CampaignTransactionsTab campaign={campaign} />
        </div>
      )}

      {showDeleteConfirm && (
        <ConfirmDialog
          title="Delete campaign?"
          message="This draft campaign will be permanently deleted. This cannot be undone."
          confirmLabel="Delete"
          isConfirming={isDeleting}
          onCancel={() => setShowDeleteConfirm(false)}
          onConfirm={handleConfirmDelete}
        />
      )}
    </>
  )
}

export default MyCampaignDetailsPage
