import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { fetchCampaignComments } from '../lib/api'
import { formatEthDisplay, formatEth, formatDate } from '../utils/format'
import { formatCommentTimestamp, groupComments } from '../utils/comments'
import ArchiveCampaignModal from './ArchiveCampaignModal'
import Button from './ui/Button'
import TabButton from './ui/TabButton'

function computeProgressPercent(amountRaised, goal) {
  if (!goal || goal === '0') return 0
  const percent = (BigInt(amountRaised) * 10000n) / BigInt(goal)
  return Math.min(100, Number(percent) / 100)
}

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
  onWithdraw,
  onArchive,
  publishLabel = 'Publish',
  isPublishing = false,
  publishError = null,
  withdrawLabel = 'Withdraw Funds',
  isWithdrawing = false,
  withdrawError = null,
  isArchiving = false,
  archiveError = null,
}) {
  const previewCover = campaign.assets?.find((asset) => asset.isCover) || campaign.assets?.[0]
  const isPublished = campaign.status === 'published'
  const isArchived = campaign.status === 'archived'
  const chainDataAvailable = (isPublished || isArchived) && campaign.onChainAvailable !== false
  const progress = chainDataAvailable ? computeProgressPercent(campaign.amountRaised, campaign.goal) : 0
  const canWithdraw =
    chainDataAvailable &&
    !campaign.withdrawn &&
    onWithdraw &&
    BigInt(campaign.amountRaised || '0') >= BigInt(campaign.goal || '0')
  const [activeTab, setActiveTab] = useState('story')
  const [comments, setComments] = useState([])
  const [expandedReplies, setExpandedReplies] = useState({})
  const [showArchiveModal, setShowArchiveModal] = useState(false)

  async function handleConfirmArchive(note) {
    try {
      await onArchive(note)
      setShowArchiveModal(false)
    } catch {
      // error is surfaced via archiveError, keep the modal open so the user can retry
    }
  }

  useEffect(() => {
    if (!campaign.id) return

    fetchCampaignComments(campaign.id)
      .then(({ items }) => setComments(items))
      .catch(() => {})
  }, [campaign.id])

  function toggleReplies(commentId) {
    setExpandedReplies((prev) => ({ ...prev, [commentId]: !prev[commentId] }))
  }

  const { rootComments, repliesByParent } = groupComments(comments)

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
              isArchived
                ? 'bg-slate-100 text-slate-600'
                : isPublished
                  ? 'bg-emerald-50 text-emerald-600'
                  : 'bg-amber-50 text-amber-600'
            }`}
          >
            <svg viewBox="0 0 24 24" fill="currentColor" className="h-3.5 w-3.5">
              <path d="M12 2 2 7v6c0 5 4.5 8 10 9 5.5-1 10-4 10-9V7l-10-5Z" />
            </svg>
            {isArchived ? 'Archived' : isPublished ? 'Published' : 'Draft — not published yet'}
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
              {rootComments.length === 0 && (
                <p className="text-sm text-slate-500">No comments yet.</p>
              )}
              {rootComments.map((comment) => {
                const replies = repliesByParent[comment.id] || []
                const isExpanded = Boolean(expandedReplies[comment.id])

                return (
                  <div key={comment.id} className="flex flex-col gap-3">
                    <div className="flex gap-3">
                      <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-slate-100 text-xs font-medium text-slate-600">
                        {comment.authorName.charAt(0).toUpperCase()}
                      </div>
                      <div className="flex flex-col gap-0.5">
                        <div className="flex items-center gap-2">
                          <span className="text-sm font-medium text-slate-900">{comment.authorName}</span>
                          <span className="text-xs text-slate-400">{formatCommentTimestamp(comment.createdAt)}</span>
                        </div>
                        <p className="text-sm text-slate-600">{comment.text}</p>
                      </div>
                    </div>

                    {replies.length > 0 && (
                      <div className="ml-11 flex flex-col gap-3">
                        {!isExpanded ? (
                          <button
                            type="button"
                            onClick={() => toggleReplies(comment.id)}
                            className="self-start cursor-pointer text-xs font-medium text-slate-500 hover:text-indigo-600"
                          >
                            {replies.length} {replies.length === 1 ? 'reply' : 'replies'}
                          </button>
                        ) : (
                          <>
                            {replies.map((reply) => (
                              <div key={reply.id} className="flex gap-2.5">
                                <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-slate-100 text-xs font-medium text-slate-600">
                                  {reply.authorName.charAt(0).toUpperCase()}
                                </div>
                                <div className="flex flex-col gap-0.5">
                                  <div className="flex items-center gap-2">
                                    <span className="text-sm font-medium text-slate-900">{reply.authorName}</span>
                                    <span className="text-xs text-slate-400">{formatCommentTimestamp(reply.createdAt)}</span>
                                  </div>
                                  <p className="text-sm text-slate-600">{reply.text}</p>
                                </div>
                              </div>
                            ))}
                            <button
                              type="button"
                              onClick={() => toggleReplies(comment.id)}
                              className="self-start cursor-pointer text-xs font-medium text-slate-500 hover:text-indigo-600"
                            >
                              Hide replies
                            </button>
                          </>
                        )}
                      </div>
                    )}
                  </div>
                )
              })}
            </div>
          )}
        </div>

        <div className="flex flex-col gap-4 self-start rounded-2xl border border-slate-200 bg-white p-6 shadow-sm lg:sticky lg:top-24">
          <div>
            <div className="h-2 w-full overflow-hidden rounded-full bg-slate-100">
              <div className="h-full rounded-full bg-indigo-600" style={{ width: `${progress}%` }} />
            </div>
            {chainDataAvailable ? (
              <>
                <p className="mt-3 text-lg font-semibold text-slate-900">{formatEth(campaign.amountRaised)} raised</p>
                <p className="text-sm text-slate-500">of {formatEth(campaign.goal)} goal</p>
                <p className="mt-2 text-xs text-slate-400">Ends {formatDate(campaign.deadline)}</p>
              </>
            ) : isPublished || isArchived ? (
              <p className="mt-3 text-sm font-medium text-amber-600">
                Live blockchain data is temporarily unavailable. Please refresh shortly.
              </p>
            ) : (
              <>
                <p className="mt-3 text-lg font-semibold text-slate-900">0 ETH raised</p>
                <p className="text-sm text-slate-500">of {formatEthDisplay(campaign.targetEth)} ETH goal</p>
                <p className="mt-2 text-xs text-slate-400">Runs for {campaign.durationDays} days once published</p>
              </>
            )}
          </div>

          <div className="flex flex-col gap-2 border-t border-slate-100 pt-4">
            {isPublished || isArchived ? (
              <>
                <Link
                  to={`/campaigns/${campaign.id}`}
                  className="inline-flex cursor-pointer items-center justify-center rounded-lg bg-indigo-600 px-4 py-3 text-base font-medium text-white shadow-sm transition-colors hover:bg-indigo-500"
                >
                  View Live Campaign
                </Link>

                {campaign.withdrawn ? (
                  <span className="inline-flex items-center justify-center gap-1.5 rounded-lg bg-emerald-50 px-4 py-3 text-sm font-medium text-emerald-600">
                    Funds withdrawn
                  </span>
                ) : (
                  canWithdraw && (
                    <>
                      <Button
                        variant="secondary"
                        onClick={onWithdraw}
                        disabled={isWithdrawing}
                        className="justify-center py-3 text-base"
                      >
                        {withdrawLabel}
                      </Button>
                      {withdrawError && <p className="text-xs font-medium text-rose-500">{withdrawError}</p>}
                    </>
                  )
                )}

                {isArchived ? (
                  <div className="rounded-lg bg-slate-50 px-4 py-3 text-xs text-slate-500">
                    <p className="font-medium text-slate-700">
                      Archived{campaign.archivedAt ? ` on ${new Date(campaign.archivedAt).toLocaleString()}` : ''}
                    </p>
                    {campaign.archiveNote && <p className="mt-1">{campaign.archiveNote}</p>}
                  </div>
                ) : (
                  onArchive && (
                    <>
                      <Button
                        variant="secondary"
                        onClick={() => setShowArchiveModal(true)}
                        className="justify-center"
                      >
                        Archive Campaign
                      </Button>
                      {archiveError && !showArchiveModal && (
                        <p className="text-xs font-medium text-rose-500">{archiveError}</p>
                      )}
                    </>
                  )
                )}
              </>
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

      {showArchiveModal && (
        <ArchiveCampaignModal
          isArchiving={isArchiving}
          error={archiveError}
          onCancel={() => setShowArchiveModal(false)}
          onConfirm={handleConfirmArchive}
        />
      )}
    </div>
  )
}

export default CampaignPreview
