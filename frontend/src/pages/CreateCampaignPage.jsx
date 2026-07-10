import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth0 } from '@auth0/auth0-react'
import { uploadAsset, createMyCampaign, fetchMyCampaign } from '../lib/api'
import { formatEthDisplay } from '../utils/format'
import { usePublishCampaign } from '../hooks/usePublishCampaign'
import Button from '../components/ui/Button'
import CampaignPreview from '../components/CampaignPreview'

const PUBLISH_LABELS = {
  connecting: 'Connecting wallet...',
  signing: 'Confirm in wallet...',
  confirming: 'Waiting for confirmation...',
  linking: 'Finalizing...',
}

const COUNTRIES = [
  'United States',
  'United Kingdom',
  'Canada',
  'Australia',
  'Germany',
  'France',
  'India',
  'Bangladesh',
  'Japan',
  'Brazil',
  'Other',
]

const FUNDRAISING_FOR_OPTIONS = ['Yourself', 'Someone else', 'Charity']

const CATEGORIES = [
  'Medical & Health',
  'Education',
  'Community & Environment',
  'Animals & Pets',
  'Emergency Relief',
  'Other',
]

const inputClasses =
  'w-full rounded-xl border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 placeholder:text-slate-400 shadow-sm transition focus:outline-none focus:ring-2 focus:ring-indigo-500/40 focus:border-indigo-500'

const labelClasses = 'text-sm font-semibold text-slate-800'
const hintClasses = 'text-xs text-slate-400'

function SectionCard({ title, description, children }) {
  return (
    <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm sm:p-8">
      <div className="mb-6">
        <h2 className="text-base font-semibold text-slate-900">{title}</h2>
        {description && <p className="mt-1 text-sm text-slate-500">{description}</p>}
      </div>
      <div className="flex flex-col gap-5">{children}</div>
    </div>
  )
}

function StepPill({ index, label, active, done }) {
  return (
    <div className="flex items-center gap-2">
      <div
        className={`flex h-7 w-7 shrink-0 items-center justify-center rounded-full text-xs font-semibold ${
          done
            ? 'bg-indigo-600 text-white'
            : active
              ? 'bg-indigo-600 text-white'
              : 'bg-slate-100 text-slate-400'
        }`}
      >
        {done ? (
          <svg viewBox="0 0 20 20" fill="currentColor" className="h-4 w-4">
            <path
              fillRule="evenodd"
              d="M16.7 5.3a1 1 0 0 1 0 1.4l-7.5 7.5a1 1 0 0 1-1.4 0l-3.5-3.5a1 1 0 1 1 1.4-1.4l2.8 2.8 6.8-6.8a1 1 0 0 1 1.4 0Z"
              clipRule="evenodd"
            />
          </svg>
        ) : (
          index
        )}
      </div>
      <span className={`text-sm font-medium ${active ? 'text-slate-900' : 'text-slate-400'}`}>{label}</span>
    </div>
  )
}

function Spinner() {
  return (
    <svg viewBox="0 0 24 24" fill="none" className="h-6 w-6 animate-spin text-white">
      <circle cx="12" cy="12" r="9" stroke="currentColor" strokeWidth="3" strokeOpacity="0.3" />
      <path d="M21 12a9 9 0 0 0-9-9" stroke="currentColor" strokeWidth="3" strokeLinecap="round" />
    </svg>
  )
}

function LivePreviewCard({ title, description, target, country, category, durationDays, fundraisingFor, coverPhotoPreview }) {
  return (
    <div className="sticky top-24 overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
      {coverPhotoPreview ? (
        <img src={coverPhotoPreview} alt="Cover" className="aspect-[16/10] w-full object-cover" />
      ) : (
        <div className="flex aspect-[16/10] w-full flex-col items-center justify-center gap-2 bg-gradient-to-br from-indigo-50 via-white to-emerald-50 text-indigo-300">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="h-10 w-10">
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M3 16.5V6a1.5 1.5 0 0 1 1.5-1.5h15A1.5 1.5 0 0 1 21 6v10.5m-18 0A1.5 1.5 0 0 0 4.5 18h15a1.5 1.5 0 0 0 1.5-1.5m-18 0 5.47-5.47a1.5 1.5 0 0 1 2.12 0l2.66 2.66m7.75 2.81-3.66-3.66a1.5 1.5 0 0 0-2.12 0l-.97.97"
            />
          </svg>
          <span className="text-xs font-medium">Cover photo preview</span>
        </div>
      )}

      <div className="flex flex-col gap-3 p-5">
        <div className="flex flex-wrap items-center gap-2">
          <span className="inline-flex items-center rounded-full bg-emerald-50 px-2.5 py-1 text-xs font-medium text-emerald-600">
            {country}
          </span>
          <span className="inline-flex items-center rounded-full bg-indigo-50 px-2.5 py-1 text-xs font-medium text-indigo-600">
            {category}
          </span>
        </div>

        <h3 className="text-lg font-semibold leading-snug text-slate-900">
          {title || 'Your campaign title will appear here'}
        </h3>

        <p className="line-clamp-3 text-sm text-slate-500">
          {description || 'A short, compelling description of your campaign will show up here as you type.'}
        </p>

        <div className="mt-1">
          <div className="h-2 w-full overflow-hidden rounded-full bg-slate-100">
            <div className="h-full w-0 rounded-full bg-indigo-600" />
          </div>
          <div className="mt-2 flex items-baseline justify-between">
            <span className="text-sm font-semibold text-slate-900">0 ETH raised</span>
            <span className="text-xs text-slate-400">of {formatEthDisplay(target)} ETH goal</span>
          </div>
        </div>

        <div className="flex flex-col gap-1.5 border-t border-slate-100 pt-3 text-xs text-slate-500">
          <div className="flex items-center gap-2">
            <svg viewBox="0 0 24 24" fill="currentColor" className="h-4 w-4 text-slate-400">
              <path d="M12 12a5 5 0 1 0 0-10 5 5 0 0 0 0 10Zm0 2c-4.42 0-8 2.24-8 5v1a1 1 0 0 0 1 1h14a1 1 0 0 0 1-1v-1c0-2.76-3.58-5-8-5Z" />
            </svg>
            Fundraising for {fundraisingFor.toLowerCase()}
          </div>
          <div className="flex items-center gap-2">
            <svg viewBox="0 0 24 24" fill="currentColor" className="h-4 w-4 text-slate-400">
              <path
                fillRule="evenodd"
                d="M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20Zm1 5a1 1 0 1 0-2 0v5a1 1 0 0 0 .5.87l3.5 2a1 1 0 0 0 1-1.74L13 11.4V7Z"
                clipRule="evenodd"
              />
            </svg>
            Runs for {durationDays || 0} days once published
          </div>
        </div>
      </div>
    </div>
  )
}

function CreateCampaignPreviewStep({ campaign, provider, account, onConnectWallet, showToast, onBack, onEdit, onRefresh }) {
  const { getAccessTokenSilently } = useAuth0()

  const publishHook = usePublishCampaign({
    campaign,
    provider,
    account,
    onConnectWallet,
    onPublished: async () => {
      showToast?.('Your campaign is live!')
      const accessToken = await getAccessTokenSilently()
      const updated = await fetchMyCampaign(accessToken, campaign.id)
      onRefresh(updated)
    },
  })

  const isPublishing = ['connecting', 'signing', 'confirming', 'linking'].includes(publishHook.phase)
  const publishLabel = PUBLISH_LABELS[publishHook.phase] || (publishHook.pendingLink ? 'Finish Publishing' : 'Publish')

  return (
    <CampaignPreview
      campaign={campaign}
      onBack={onBack}
      onEdit={onEdit}
      onPublish={publishHook.pendingLink ? publishHook.retryLinking : publishHook.publish}
      publishLabel={publishLabel}
      isPublishing={isPublishing}
      publishError={publishHook.error}
    />
  )
}

function CreateCampaignPage({ provider, account, onConnectWallet, showToast }) {
  const navigate = useNavigate()
  const { getAccessTokenSilently } = useAuth0()
  const [step, setStep] = useState('form')
  const [country, setCountry] = useState(COUNTRIES[0])
  const [category, setCategory] = useState(CATEGORIES[0])
  const [title, setTitle] = useState('')
  const [target, setTarget] = useState('')
  const [durationDays, setDurationDays] = useState('30')
  const [description, setDescription] = useState('')
  const [fundraisingFor, setFundraisingFor] = useState(FUNDRAISING_FOR_OPTIONS[0])
  const [coverPhotoPreview, setCoverPhotoPreview] = useState('')
  const [coverAsset, setCoverAsset] = useState(null)
  const [isUploadingImage, setIsUploadingImage] = useState(false)
  const [uploadError, setUploadError] = useState(null)
  const [isDragging, setIsDragging] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [submitError, setSubmitError] = useState(null)
  const [campaignPreview, setCampaignPreview] = useState(null)

  useEffect(() => {
    if (!coverPhotoPreview) return
    return () => URL.revokeObjectURL(coverPhotoPreview)
  }, [coverPhotoPreview])

  async function uploadCoverFile(file) {
    if (!file || !file.type.startsWith('image/') || isUploadingImage) return

    setCoverPhotoPreview(URL.createObjectURL(file))
    setCoverAsset(null)
    setUploadError(null)
    setIsUploadingImage(true)

    try {
      const accessToken = await getAccessTokenSilently()
      const asset = await uploadAsset(accessToken, file)
      setCoverAsset(asset)
    } catch (err) {
      setUploadError(err.message)
    } finally {
      setIsUploadingImage(false)
    }
  }

  function handleCoverPhotoChange(e) {
    uploadCoverFile(e.target.files?.[0])
  }

  function handleDrop(e) {
    e.preventDefault()
    setIsDragging(false)
    uploadCoverFile(e.dataTransfer.files?.[0])
  }

  function handleRemovePhoto() {
    setCoverPhotoPreview('')
    setCoverAsset(null)
    setUploadError(null)
  }

  async function handleSubmit(e) {
    e.preventDefault()

    if (!coverAsset) {
      setSubmitError('Please add at least one image before continuing.')
      return
    }

    setSubmitError(null)
    setIsSubmitting(true)

    try {
      const accessToken = await getAccessTokenSilently()
      const created = await createMyCampaign(accessToken, {
        country,
        category,
        title,
        description,
        targetEth: target,
        durationDays: Number(durationDays),
        fundraisingFor,
        assetIds: [coverAsset.id],
      })
      const details = await fetchMyCampaign(accessToken, created.id)
      setCampaignPreview(details)
      setStep('preview')
    } catch (err) {
      setSubmitError(err.message)
    } finally {
      setIsSubmitting(false)
    }
  }

  if (step === 'preview' && campaignPreview) {
    return (
      <CreateCampaignPreviewStep
        campaign={campaignPreview}
        provider={provider}
        account={account}
        onConnectWallet={onConnectWallet}
        showToast={showToast}
        onBack={() => navigate('/my-campaigns')}
        onEdit={() => setStep('form')}
        onRefresh={setCampaignPreview}
      />
    )
  }

  return (
    <div className="mx-auto max-w-5xl">
      <div className="mb-8">
        <h1 className="text-2xl font-bold tracking-tight text-slate-900">Create a campaign</h1>
        <p className="mt-1 text-sm text-slate-500">Tell your story and set a goal — you can preview before publishing.</p>
      </div>

      <div className="mb-8 flex flex-wrap items-center gap-6 rounded-2xl border border-slate-200 bg-white px-5 py-4 shadow-sm">
        <StepPill index={1} label="Campaign details" active />
        <div className="h-px w-8 bg-slate-200 sm:w-16" />
        <StepPill index={2} label="Preview" />
      </div>

      <form onSubmit={handleSubmit} className="grid grid-cols-1 gap-8 lg:grid-cols-[minmax(0,1fr)_360px]">
        <div className="flex flex-col gap-6">
          <SectionCard title="Basics" description="Where is this campaign based, and who is it for?">
            <div className="grid grid-cols-1 gap-5 sm:grid-cols-2">
              <div className="flex flex-col gap-1.5">
                <label htmlFor="country" className={labelClasses}>
                  Country
                </label>
                <select
                  id="country"
                  value={country}
                  onChange={(e) => setCountry(e.target.value)}
                  className={inputClasses}
                >
                  {COUNTRIES.map((option) => (
                    <option key={option} value={option}>
                      {option}
                    </option>
                  ))}
                </select>
              </div>

              <div className="flex flex-col gap-1.5">
                <label htmlFor="category" className={labelClasses}>
                  Category
                </label>
                <select
                  id="category"
                  value={category}
                  onChange={(e) => setCategory(e.target.value)}
                  className={inputClasses}
                >
                  {CATEGORIES.map((option) => (
                    <option key={option} value={option}>
                      {option}
                    </option>
                  ))}
                </select>
              </div>

              <div className="flex flex-col gap-1.5">
                <label htmlFor="fundraisingFor" className={labelClasses}>
                  Who are you fundraising for
                </label>
                <select
                  id="fundraisingFor"
                  value={fundraisingFor}
                  onChange={(e) => setFundraisingFor(e.target.value)}
                  className={inputClasses}
                >
                  {FUNDRAISING_FOR_OPTIONS.map((option) => (
                    <option key={option} value={option}>
                      {option}
                    </option>
                  ))}
                </select>
              </div>
            </div>
          </SectionCard>

          <SectionCard title="Your story" description="Give supporters a clear title and a compelling description.">
            <div className="flex flex-col gap-1.5">
              <label htmlFor="title" className={labelClasses}>
                Campaign title
              </label>
              <input
                id="title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="e.g. Help rebuild our community library"
                required
                className={inputClasses}
              />
            </div>

            <div className="flex flex-col gap-1.5">
              <label htmlFor="target" className={labelClasses}>
                Funding target
              </label>
              <div className="relative">
                <input
                  id="target"
                  type="number"
                  min="0"
                  step="0.01"
                  value={target}
                  onChange={(e) => setTarget(e.target.value)}
                  placeholder="0.00"
                  required
                  className={`${inputClasses} pr-14`}
                />
                <span className="absolute inset-y-0 right-4 flex items-center text-sm font-semibold text-slate-400">
                  ETH
                </span>
              </div>
              <span className={hintClasses}>This is the amount you're aiming to raise, in ETH.</span>
            </div>

            <div className="flex flex-col gap-1.5">
              <label htmlFor="durationDays" className={labelClasses}>
                Campaign duration
              </label>
              <div className="relative">
                <input
                  id="durationDays"
                  type="number"
                  min="1"
                  max="365"
                  step="1"
                  value={durationDays}
                  onChange={(e) => setDurationDays(e.target.value)}
                  placeholder="30"
                  required
                  className={`${inputClasses} pr-16`}
                />
                <span className="absolute inset-y-0 right-4 flex items-center text-sm font-semibold text-slate-400">
                  days
                </span>
              </div>
              <span className={hintClasses}>How many days should this campaign run once published (1-365).</span>
            </div>

            <div className="flex flex-col gap-1.5">
              <label htmlFor="description" className={labelClasses}>
                Description
              </label>
              <textarea
                id="description"
                rows={6}
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Explain what you're raising money for, why it matters, and how the funds will be used."
                required
                className={`${inputClasses} resize-none`}
              />
            </div>
          </SectionCard>

          <SectionCard
            title="Cover photo"
            description="A clear, high-quality image helps your campaign stand out. At least one image is required."
          >
            <label
              htmlFor="coverPhoto"
              onDragOver={(e) => {
                e.preventDefault()
                setIsDragging(true)
              }}
              onDragLeave={() => setIsDragging(false)}
              onDrop={handleDrop}
              className={`relative flex cursor-pointer flex-col items-center justify-center gap-3 rounded-xl border-2 border-dashed px-6 py-10 text-center transition ${
                isDragging ? 'border-indigo-400 bg-indigo-50' : 'border-slate-200 bg-slate-50 hover:border-indigo-300 hover:bg-indigo-50/50'
              }`}
            >
              {coverPhotoPreview ? (
                <div className="relative w-full max-w-sm">
                  <img
                    src={coverPhotoPreview}
                    alt="Cover preview"
                    className="aspect-video w-full rounded-lg object-cover shadow-sm"
                  />
                  {isUploadingImage && (
                    <div className="absolute inset-0 flex flex-col items-center justify-center gap-2 rounded-lg bg-slate-900/60">
                      <Spinner />
                      <span className="text-xs font-medium text-white">Uploading...</span>
                    </div>
                  )}
                  {!isUploadingImage && coverAsset && (
                    <span className="absolute right-2 top-2 inline-flex items-center gap-1 rounded-full bg-emerald-500 px-2 py-1 text-xs font-medium text-white">
                      <svg viewBox="0 0 20 20" fill="currentColor" className="h-3 w-3">
                        <path
                          fillRule="evenodd"
                          d="M16.7 5.3a1 1 0 0 1 0 1.4l-7.5 7.5a1 1 0 0 1-1.4 0l-3.5-3.5a1 1 0 1 1 1.4-1.4l2.8 2.8 6.8-6.8a1 1 0 0 1 1.4 0Z"
                          clipRule="evenodd"
                        />
                      </svg>
                      Uploaded
                    </span>
                  )}
                </div>
              ) : (
                <>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="h-9 w-9 text-slate-400">
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M3 16.5V6a1.5 1.5 0 0 1 1.5-1.5h15A1.5 1.5 0 0 1 21 6v10.5m-18 0A1.5 1.5 0 0 0 4.5 18h15a1.5 1.5 0 0 0 1.5-1.5m-18 0 5.47-5.47a1.5 1.5 0 0 1 2.12 0l2.66 2.66m7.75 2.81-3.66-3.66a1.5 1.5 0 0 0-2.12 0l-.97.97"
                    />
                  </svg>
                  <div>
                    <span className="text-sm font-semibold text-indigo-600">Click to upload</span>
                    <span className="text-sm text-slate-500"> or drag and drop</span>
                  </div>
                  <span className={hintClasses}>PNG or JPG, up to 10MB</span>
                </>
              )}
              <input
                id="coverPhoto"
                type="file"
                accept="image/*"
                onChange={handleCoverPhotoChange}
                disabled={isUploadingImage}
                className="hidden"
              />
            </label>

            {uploadError && <p className="text-xs font-medium text-rose-500">{uploadError}</p>}

            {coverPhotoPreview && !isUploadingImage && (
              <button
                type="button"
                onClick={handleRemovePhoto}
                className="self-start cursor-pointer text-xs font-medium text-rose-500 hover:text-rose-600"
              >
                Remove photo
              </button>
            )}
          </SectionCard>

          {submitError && <p className="text-sm font-medium text-rose-500">{submitError}</p>}

          <div className="flex justify-end">
            <Button type="submit" className="px-8 py-3 text-base" disabled={isUploadingImage || isSubmitting}>
              {isSubmitting ? 'Creating...' : 'Preview'}
            </Button>
          </div>
        </div>

        <div>
          <LivePreviewCard
            title={title}
            description={description}
            target={target}
            country={country}
            category={category}
            durationDays={durationDays}
            fundraisingFor={fundraisingFor}
            coverPhotoPreview={coverPhotoPreview}
          />
        </div>
      </form>
    </div>
  )
}

export default CreateCampaignPage
