import { useState } from 'react'
import { Link } from 'react-router-dom'
import { formatEthDisplay } from '../utils/format'
import Button from './ui/Button'
import TabButton from './ui/TabButton'

const SEED_COMMENTS = [
  { id: 1, author: 'Alex Morgan', text: 'This is such a great initiative, happy to support!', postedAt: '2 days ago' },
  { id: 2, author: 'Priya Singh', text: 'Following this closely, good luck reaching the goal.', postedAt: '5 hours ago' },
]

function TrashIcon() {
  return (
    <svg viewBox="0 0 24 24" fill="currentColor" className="h-4 w-4">
      <path d="M9 3a1 1 0 0 0-1 1v1H4.5a1 1 0 1 0 0 2H5v12a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V7h.5a1 1 0 1 0 0-2H16V4a1 1 0 0 0-1-1H9Zm1 2h4v1h-4V5ZM8 7h8v12H8V7Zm2 2a1 1 0 0 0-1 1v7a1 1 0 1 0 2 0v-7a1 1 0 0 0-1-1Zm4 0a1 1 0 0 0-1 1v7a1 1 0 1 0 2 0v-7a1 1 0 0 0-1-1Z" />
    </svg>
  )
}

function CampaignPreview({
  campaign,
  onBack,
  onPublish,
  onEdit,
  onDelete,
  publishLabel = 'Publish',
  isPublishing = false,
  publishError = null,
}) {
  const previewCover = campaign.assets?.find((asset) => asset.isCover) || campaign.assets?.[0]
  const isPublished = campaign.status === 'published'
  const [activeTab, setActiveTab] = useState('story')
  const [comments, setComments] = useState(SEED_COMMENTS)
  const [commentText, setCommentText] = useState('')
  const [showCommentForm, setShowCommentForm] = useState(false)

  function handlePostComment(e) {
    e.preventDefault()
    if (!commentText.trim()) return

    setComments((prev) => [{ id: Date.now(), author: 'You', text: commentText.trim(), postedAt: 'Just now' }, ...prev])
    setCommentText('')
    setShowCommentForm(false)
  }

  return (
    <div className="mx-auto max-w-5xl">
      <button
        type="button"
        onClick={onBack}
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
        <div className="flex items-center gap-3">
          <span
            className={`inline-flex items-center gap-1.5 rounded-full px-3 py-1.5 text-xs font-medium ${
              isPublished ? 'bg-emerald-50 text-emerald-600' : 'bg-amber-50 text-amber-600'
            }`}
          >
            <svg viewBox="0 0 24 24" fill="currentColor" className="h-3.5 w-3.5">
              <path d="M12 2 2 7v6c0 5 4.5 8 10 9 5.5-1 10-4 10-9V7l-10-5Z" />
            </svg>
            {isPublished ? 'Published' : 'Draft — not published yet'}
          </span>

          {onDelete && campaign.status === 'draft' && (
            <div className="group relative">
              <button
                type="button"
                onClick={onDelete}
                aria-label="Delete campaign"
                className="flex cursor-pointer items-center justify-center rounded-lg p-2 text-red-600 hover:bg-red-50"
              >
                <TrashIcon />
              </button>
              <span className="pointer-events-none absolute left-1/2 top-full z-10 mt-1 -translate-x-1/2 whitespace-nowrap rounded-md bg-slate-900 px-2 py-1 text-xs font-medium text-white opacity-0 transition-opacity group-hover:opacity-100">
                Delete
              </span>
            </div>
          )}
        </div>
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
              {campaign.country}
            </span>
            <span className="inline-flex items-center rounded-full bg-indigo-50 px-2.5 py-1 text-xs font-medium text-indigo-600">
              {campaign.category}
            </span>
            <span className="inline-flex items-center rounded-full bg-slate-100 px-2.5 py-1 text-xs font-medium text-slate-600">
              For {campaign.fundraisingFor.toLowerCase()}
            </span>
          </div>

          <h2 className="text-2xl font-bold text-slate-900">{campaign.title || 'Untitled campaign'}</h2>

          <div className="flex gap-6 border-b border-slate-200">
            <TabButton active={activeTab === 'story'} onClick={() => setActiveTab('story')}>
              Story
            </TabButton>
            <TabButton active={activeTab === 'comments'} onClick={() => setActiveTab('comments')}>
              Comments ({comments.length})
            </TabButton>
          </div>

          {activeTab === 'story' ? (
            <p className="whitespace-pre-line text-sm leading-relaxed text-slate-600">{campaign.description}</p>
          ) : (
            <div className="flex flex-col gap-4">
              <div className="flex items-center justify-end">
                {!showCommentForm && (
                  <Button variant="secondary" onClick={() => setShowCommentForm(true)}>
                    Add Comment
                  </Button>
                )}
              </div>

              {showCommentForm && (
                <form onSubmit={handlePostComment} className="flex flex-col gap-2">
                  <textarea
                    value={commentText}
                    onChange={(e) => setCommentText(e.target.value)}
                    placeholder="Leave a comment of support..."
                    rows={3}
                    autoFocus
                    className="w-full rounded-xl border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 placeholder:text-slate-400 shadow-sm transition focus:outline-none focus:ring-2 focus:ring-indigo-500/40 focus:border-indigo-500"
                  />
                  <div className="flex justify-end gap-2">
                    <Button type="button" variant="secondary" onClick={() => setShowCommentForm(false)}>
                      Cancel
                    </Button>
                    <Button type="submit">Post Comment</Button>
                  </div>
                </form>
              )}

              <div className="flex flex-col gap-4">
                {comments.map((comment) => (
                  <div key={comment.id} className="flex gap-3">
                    <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-slate-100 text-xs font-medium text-slate-600">
                      {comment.author.charAt(0).toUpperCase()}
                    </div>
                    <div className="flex flex-col gap-0.5">
                      <div className="flex items-center gap-2">
                        <span className="text-sm font-medium text-slate-900">{comment.author}</span>
                        <span className="text-xs text-slate-400">{comment.postedAt}</span>
                      </div>
                      <p className="text-sm text-slate-600">{comment.text}</p>
                      <button
                        type="button"
                        className="mt-1 self-start cursor-pointer text-xs font-medium text-slate-500 hover:text-indigo-600"
                      >
                        Reply
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>

        <div className="flex flex-col gap-4 self-start rounded-2xl border border-slate-200 bg-white p-6 shadow-sm lg:sticky lg:top-24">
          <div>
            <div className="h-2 w-full overflow-hidden rounded-full bg-slate-100">
              <div className="h-full w-0 rounded-full bg-indigo-600" />
            </div>
            <p className="mt-3 text-lg font-semibold text-slate-900">0 ETH raised</p>
            <p className="text-sm text-slate-500">of {formatEthDisplay(campaign.targetEth)} ETH goal</p>
            <p className="mt-2 text-xs text-slate-400">Runs for {campaign.durationDays} days once published</p>
          </div>

          <div className="flex flex-col gap-2 border-t border-slate-100 pt-4">
            {isPublished ? (
              <Link
                to={`/campaigns/${campaign.onChainCampaignId}`}
                className="inline-flex cursor-pointer items-center justify-center rounded-lg bg-indigo-600 px-4 py-3 text-base font-medium text-white shadow-sm transition-colors hover:bg-indigo-500"
              >
                View Live Campaign
              </Link>
            ) : (
              <>
                <Button
                  onClick={onPublish}
                  disabled={isPublishing}
                  className="justify-center py-3 text-base"
                >
                  {publishLabel}
                </Button>
                {publishError && <p className="text-xs font-medium text-rose-500">{publishError}</p>}
                {onEdit && (
                  <Button variant="secondary" onClick={onEdit} disabled={isPublishing} className="justify-center">
                    Edit details
                  </Button>
                )}
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

export default CampaignPreview
