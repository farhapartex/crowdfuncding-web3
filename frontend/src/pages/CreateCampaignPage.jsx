import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth0 } from '@auth0/auth0-react'
import { uploadAsset, createMyCampaign, fetchMyCampaign } from '../lib/api'
import Button from '../components/ui/Button'

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

function formatEthDisplay(value) {
  const parsed = Number(value)
  if (!value || Number.isNaN(parsed)) return '0'
  return parsed.toLocaleString(undefined, { maximumFractionDigits: 4 })
}

function Spinner() {
  return (
    <svg viewBox="0 0 24 24" fill="none" className="h-6 w-6 animate-spin text-white">
      <circle cx="12" cy="12" r="9" stroke="currentColor" strokeWidth="3" strokeOpacity="0.3" />
      <path d="M21 12a9 9 0 0 0-9-9" stroke="currentColor" strokeWidth="3" strokeLinecap="round" />
    </svg>
  )
}

function LivePreviewCard({ title, description, target, country, fundraisingFor, coverPhotoPreview }) {
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
        <span className="inline-flex w-fit items-center rounded-full bg-emerald-50 px-2.5 py-1 text-xs font-medium text-emerald-600">
          {country}
        </span>

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

        <div className="flex items-center gap-2 border-t border-slate-100 pt-3 text-xs text-slate-500">
          <svg viewBox="0 0 24 24" fill="currentColor" className="h-4 w-4 text-slate-400">
            <path d="M12 12a5 5 0 1 0 0-10 5 5 0 0 0 0 10Zm0 2c-4.42 0-8 2.24-8 5v1a1 1 0 0 0 1 1h14a1 1 0 0 0 1-1v-1c0-2.76-3.58-5-8-5Z" />
          </svg>
          Fundraising for {fundraisingFor.toLowerCase()}
        </div>
      </div>
    </div>
  )
}

function CreateCampaignPage() {
  const navigate = useNavigate()
  const { getAccessTokenSilently } = useAuth0()
  const [step, setStep] = useState('form')
  const [country, setCountry] = useState(COUNTRIES[0])
  const [title, setTitle] = useState('')
  const [target, setTarget] = useState('')
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
        title,
        description,
        targetEth: target,
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
    const previewCover = campaignPreview.assets?.find((asset) => asset.isCover) || campaignPreview.assets?.[0]

    return (
      <div className="mx-auto max-w-5xl">
        <button
          type="button"
          onClick={() => navigate('/my-campaigns')}
          className="mb-6 flex cursor-pointer items-center gap-1.5 text-sm font-medium text-slate-500 hover:text-slate-700"
        >
          <svg viewBox="0 0 20 20" fill="currentColor" className="h-4 w-4">
            <path
              fillRule="evenodd"
              d="M9.7 4.3a1 1 0 0 1 0 1.4L6.42 9h9.58a1 1 0 1 1 0 2H6.42l3.3 3.3a1 1 0 1 1-1.42 1.4l-5-5a1 1 0 0 1 0-1.4l5-5a1 1 0 0 1 1.4 0Z"
              clipRule="evenodd"
            />
          </svg>
          Back to Campaigns
        </button>

        <div className="mb-6 flex flex-wrap items-center justify-between gap-3">
          <div>
            <h1 className="text-2xl font-bold tracking-tight text-slate-900">Preview your campaign</h1>
            <p className="mt-1 text-sm text-slate-500">This is how it will look to potential supporters.</p>
          </div>
          <span className="inline-flex items-center gap-1.5 rounded-full bg-amber-50 px-3 py-1.5 text-xs font-medium text-amber-600">
            <svg viewBox="0 0 24 24" fill="currentColor" className="h-3.5 w-3.5">
              <path d="M12 2 2 7v6c0 5 4.5 8 10 9 5.5-1 10-4 10-9V7l-10-5Z" />
            </svg>
            {campaignPreview.status === 'draft' ? 'Draft — not published yet' : campaignPreview.status}
          </span>
        </div>

        <div className="grid grid-cols-1 gap-8 lg:grid-cols-3">
          <div className="flex flex-col gap-6 lg:col-span-2">
            <div className="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
              {previewCover ? (
                <img src={previewCover.url} alt="Cover" className="aspect-[16/9] w-full object-cover" />
              ) : (
                <div className="flex aspect-[16/9] w-full items-center justify-center bg-gradient-to-br from-indigo-50 via-white to-emerald-50 text-sm font-medium text-indigo-300">
                  No cover photo
                </div>
              )}
            </div>

            <div className="flex flex-wrap items-center gap-2">
              <span className="inline-flex items-center rounded-full bg-emerald-50 px-2.5 py-1 text-xs font-medium text-emerald-600">
                {campaignPreview.country}
              </span>
              <span className="inline-flex items-center rounded-full bg-slate-100 px-2.5 py-1 text-xs font-medium text-slate-600">
                For {campaignPreview.fundraisingFor.toLowerCase()}
              </span>
            </div>

            <h2 className="text-2xl font-bold text-slate-900">{campaignPreview.title || 'Untitled campaign'}</h2>
            <p className="whitespace-pre-line text-sm leading-relaxed text-slate-600">
              {campaignPreview.description}
            </p>
          </div>

          <div className="flex flex-col gap-4 self-start rounded-2xl border border-slate-200 bg-white p-6 shadow-sm lg:sticky lg:top-24">
            <div>
              <div className="h-2 w-full overflow-hidden rounded-full bg-slate-100">
                <div className="h-full w-0 rounded-full bg-indigo-600" />
              </div>
              <p className="mt-3 text-lg font-semibold text-slate-900">0 ETH raised</p>
              <p className="text-sm text-slate-500">of {formatEthDisplay(campaignPreview.targetEth)} ETH goal</p>
            </div>

            <div className="flex flex-col gap-2 border-t border-slate-100 pt-4">
              <Button onClick={() => navigate('/my-campaigns')} className="justify-center py-3 text-base">
                Publish
              </Button>
              <Button variant="secondary" onClick={() => setStep('form')} className="justify-center">
                Edit details
              </Button>
            </div>
          </div>
        </div>
      </div>
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
            fundraisingFor={fundraisingFor}
            coverPhotoPreview={coverPhotoPreview}
          />
        </div>
      </form>
    </div>
  )
}

export default CreateCampaignPage
