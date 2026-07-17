import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useAuth0 } from '@auth0/auth0-react'
import { parseEther } from 'ethers'
import {
  fetchCampaign,
  fetchContributors,
  fetchPublicProfile,
  fetchCampaignComments,
  postCampaignComment,
  postCommentReply,
} from '../lib/api'
import { getCrowdFundingContract } from '../lib/crowdFundingContract'
import { useAccessToken } from '../hooks/useAccessToken'
import { shortenAddress, formatEth, formatDate } from '../utils/format'
import { formatCommentTimestamp, groupComments } from '../utils/comments'
import StatusBadge from '../components/ui/StatusBadge'
import ContributeForm from '../components/ContributeForm'
import ShareModal from '../components/ShareModal'
import Button from '../components/ui/Button'
import TabButton from '../components/ui/TabButton'

function MessageIcon() {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6" className="h-4 w-4">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M21 12c0 4.142-4.03 7.5-9 7.5a9.7 9.7 0 0 1-2.9-.44L3 20l1.2-3.6A7.2 7.2 0 0 1 3 12c0-4.142 4.03-7.5 9-7.5s9 3.358 9 7.5Z"
      />
    </svg>
  )
}

function SendIcon() {
  return (
    <svg viewBox="0 0 24 24" fill="currentColor" className="h-4 w-4">
      <path d="M2.01 21 23 12 2.01 3 2 10l15 2-15 2z" />
    </svg>
  )
}

function ShareIcon() {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6" className="h-4 w-4">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M8.7 10.3 15.3 7m-6.6 6.7 6.6 3.3M18 6a2 2 0 1 1-4 0 2 2 0 0 1 4 0Zm0 12a2 2 0 1 1-4 0 2 2 0 0 1 4 0ZM8 12a2 2 0 1 1-4 0 2 2 0 0 1 4 0Z"
      />
    </svg>
  )
}

function computeProgressPercent(amountRaised, goal) {
  if (goal === '0') return 0
  const percent = (BigInt(amountRaised) * 10000n) / BigInt(goal)
  return Math.min(100, Number(percent) / 100)
}

function CampaignDetailsPage({ provider, account, onConnectWallet, setError, showToast }) {
  const { id } = useParams()
  const { isAuthenticated } = useAuth0()
  const getAccessToken = useAccessToken()
  const [campaign, setCampaign] = useState(null)
  const [ownerDisplayName, setOwnerDisplayName] = useState('')
  const [contributors, setContributors] = useState([])
  const [isContributing, setIsContributing] = useState(false)
  const [comments, setComments] = useState([])
  const [commentText, setCommentText] = useState('')
  const [showCommentForm, setShowCommentForm] = useState(false)
  const [isPostingComment, setIsPostingComment] = useState(false)
  const [commentError, setCommentError] = useState(null)
  const [replyingToId, setReplyingToId] = useState(null)
  const [replyText, setReplyText] = useState('')
  const [isPostingReply, setIsPostingReply] = useState(false)
  const [replyError, setReplyError] = useState(null)
  const [expandedReplies, setExpandedReplies] = useState({})
  const [activeTab, setActiveTab] = useState('story')
  const [showShareModal, setShowShareModal] = useState(false)

  useEffect(() => {
    fetchCampaign(id)
      .then(setCampaign)
      .catch((err) => setError(err.message))
    fetchContributors(id)
      .then(setContributors)
      .catch(() => {})
    fetchCampaignComments(id)
      .then(({ items }) => setComments(items))
      .catch(() => {})
  }, [id])

  useEffect(() => {
    if (!campaign) return

    fetchPublicProfile(campaign.owner)
      .then(({ displayName }) => setOwnerDisplayName(displayName))
      .catch(() => setOwnerDisplayName(''))
  }, [campaign])

  async function handleContribute(amountEth) {
    setError(null)
    setIsContributing(true)

    try {
      const signer = await provider.getSigner()
      const crowdFunding = getCrowdFundingContract(signer)

      const amountInWei = parseEther(amountEth)

      const tx = await crowdFunding.contribute(campaign.onChainCampaignId, { value: amountInWei })
      await tx.wait()

      const [updatedCampaign, updatedContributors] = await Promise.all([fetchCampaign(id), fetchContributors(id)])
      setCampaign(updatedCampaign)
      setContributors(updatedContributors)
      showToast('Thank you for your contribution!')
    } catch (err) {
      setError(err.shortMessage || err.message)
      throw err
    } finally {
      setIsContributing(false)
    }
  }

  async function handlePostComment(e) {
    e.preventDefault()
    if (!commentText.trim() || isPostingComment) return

    setCommentError(null)
    setIsPostingComment(true)

    try {
      const accessToken = await getAccessToken()
      const comment = await postCampaignComment(accessToken, id, commentText.trim())
      setComments((prev) => [comment, ...prev])
      setCommentText('')
      setShowCommentForm(false)
    } catch (err) {
      setCommentError(err.message)
    } finally {
      setIsPostingComment(false)
    }
  }

  function toggleReply(commentId) {
    setReplyError(null)
    setReplyText('')
    setReplyingToId((current) => (current === commentId ? null : commentId))
  }

  async function handleSubmitReply(e, parentId) {
    e.preventDefault()
    if (!replyText.trim() || isPostingReply) return

    setReplyError(null)
    setIsPostingReply(true)

    try {
      const accessToken = await getAccessToken()
      const reply = await postCommentReply(accessToken, parentId, replyText.trim())
      setComments((prev) => [reply, ...prev])
      setReplyText('')
      setReplyingToId(null)
      setExpandedReplies((prev) => ({ ...prev, [parentId]: true }))
    } catch (err) {
      setReplyError(err.message)
    } finally {
      setIsPostingReply(false)
    }
  }

  function toggleReplies(commentId) {
    setExpandedReplies((prev) => ({ ...prev, [commentId]: !prev[commentId] }))
  }

  if (!campaign) {
    return <p className="text-sm text-slate-500">Loading campaign...</p>
  }

  const canContribute = Date.now() / 1000 < Number(campaign.deadline)
  const progress = computeProgressPercent(campaign.amountRaised, campaign.goal)

  const { rootComments, repliesByParent } = groupComments(comments)

  return (
    <div className="flex flex-col gap-6">
      <div>
        <Link to="/campaigns" className="text-sm text-indigo-600 hover:text-indigo-500">
          &larr; Back to campaigns
        </Link>
        <h1 className="mt-2 text-2xl font-semibold text-slate-900">{campaign.title}</h1>
      </div>

      <div className="grid grid-cols-1 gap-8 lg:grid-cols-3">
        <div className="flex flex-col gap-6 lg:col-span-2">
          {campaign.coverUrl ? (
            <img
              src={campaign.coverUrl}
              alt={campaign.title}
              className="aspect-video w-full rounded-xl object-cover"
            />
          ) : (
            <div className="flex aspect-video w-full items-center justify-center rounded-xl bg-gradient-to-br from-indigo-50 to-indigo-100">
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="1.5"
                className="h-14 w-14 text-indigo-300"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M3 16.5V6.75A2.25 2.25 0 0 1 5.25 4.5h13.5A2.25 2.25 0 0 1 21 6.75v9.75M3 16.5l4.72-4.72a1.5 1.5 0 0 1 2.12 0l3.66 3.66a1.5 1.5 0 0 0 2.12 0l1.66-1.66a1.5 1.5 0 0 1 2.12 0L21 16.5M3 16.5V18a2.25 2.25 0 0 0 2.25 2.25h13.5A2.25 2.25 0 0 0 21 18v-1.5"
                />
              </svg>
            </div>
          )}

          <div className="flex items-center justify-between gap-3">
            <div className="flex items-center gap-3">
              <div className="flex h-9 w-9 items-center justify-center rounded-full bg-slate-200 text-sm font-medium text-slate-600">
                {(ownerDisplayName || campaign.owner).charAt(0).toUpperCase()}
              </div>
              <div>
                <p className="text-sm font-medium text-slate-900">
                  {ownerDisplayName || shortenAddress(campaign.owner)}
                </p>
                <p className="text-xs text-slate-500">Organizer</p>
              </div>
            </div>

            <Button variant="secondary" onClick={() => showToast('Messaging is coming soon!')} className="gap-1.5">
              <MessageIcon />
              Message
            </Button>
          </div>

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
              {isAuthenticated && (
                <div className="flex items-center justify-end">
                  {!showCommentForm && (
                    <Button variant="secondary" onClick={() => setShowCommentForm(true)}>
                      Add Comment
                    </Button>
                  )}
                </div>
              )}

              {isAuthenticated && showCommentForm && (
                <form onSubmit={handlePostComment} className="flex flex-col gap-2">
                  <textarea
                    value={commentText}
                    onChange={(e) => setCommentText(e.target.value)}
                    placeholder="Leave a comment of support..."
                    rows={3}
                    autoFocus
                    disabled={isPostingComment}
                    className="w-full rounded-xl border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 placeholder:text-slate-400 shadow-sm transition focus:outline-none focus:ring-2 focus:ring-indigo-500/40 focus:border-indigo-500"
                  />
                  {commentError && <p className="text-xs font-medium text-rose-500">{commentError}</p>}
                  <div className="flex justify-end gap-2">
                    <Button
                      type="button"
                      variant="secondary"
                      onClick={() => setShowCommentForm(false)}
                      disabled={isPostingComment}
                    >
                      Cancel
                    </Button>
                    <Button type="submit" disabled={isPostingComment}>
                      {isPostingComment ? 'Posting...' : 'Post Comment'}
                    </Button>
                  </div>
                </form>
              )}

              <div className="flex flex-col gap-4">
                {rootComments.length === 0 && (
                  <p className="text-sm text-slate-500">Be the first to leave a comment.</p>
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
                        {isAuthenticated && (
                          <button
                            type="button"
                            onClick={() => toggleReply(comment.id)}
                            className="mt-1 self-start cursor-pointer text-xs font-medium text-slate-500 hover:text-indigo-600"
                          >
                            Reply
                          </button>
                        )}
                      </div>
                    </div>

                    <div className="ml-11 flex flex-col gap-3">
                      {replies.length > 0 && !isExpanded && (
                        <button
                          type="button"
                          onClick={() => toggleReplies(comment.id)}
                          className="self-start cursor-pointer text-xs font-medium text-slate-500 hover:text-indigo-600"
                        >
                          {replies.length} {replies.length === 1 ? 'reply' : 'replies'}
                        </button>
                      )}

                      {isExpanded &&
                        replies.map((reply) => (
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

                      {isExpanded && replies.length > 0 && (
                        <button
                          type="button"
                          onClick={() => toggleReplies(comment.id)}
                          className="self-start cursor-pointer text-xs font-medium text-slate-500 hover:text-indigo-600"
                        >
                          Hide replies
                        </button>
                      )}

                      {replyingToId === comment.id && (
                        <form onSubmit={(e) => handleSubmitReply(e, comment.id)} className="flex flex-col gap-1">
                          <div className="relative max-w-md">
                            <input
                              type="text"
                              value={replyText}
                              onChange={(e) => setReplyText(e.target.value)}
                              placeholder="Write a reply..."
                              autoFocus
                              disabled={isPostingReply}
                              className="w-full rounded-full border border-slate-200 bg-white py-2 pl-4 pr-10 text-sm text-slate-900 placeholder:text-slate-400 shadow-sm outline-none transition focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/40"
                            />
                            <button
                              type="submit"
                              disabled={isPostingReply || !replyText.trim()}
                              aria-label="Send reply"
                              className="absolute right-1.5 top-1/2 flex h-7 w-7 -translate-y-1/2 cursor-pointer items-center justify-center rounded-full text-indigo-600 hover:bg-indigo-50 disabled:cursor-not-allowed disabled:opacity-40"
                            >
                              <SendIcon />
                            </button>
                          </div>
                          {replyError && <p className="text-xs font-medium text-rose-500">{replyError}</p>}
                        </form>
                      )}
                    </div>
                  </div>
                  )
                })}
              </div>
            </div>
          )}
        </div>

        <div className="flex flex-col gap-4 self-start rounded-xl border border-slate-200 bg-white p-6 shadow-sm lg:sticky lg:top-24">
          <div>
            <div className="h-2 w-full overflow-hidden rounded-full bg-slate-100">
              <div className="h-full rounded-full bg-indigo-600" style={{ width: `${progress}%` }} />
            </div>
            <p className="mt-3 text-lg font-semibold text-slate-900">{formatEth(campaign.amountRaised)} raised</p>
            <p className="text-sm text-slate-500">of {formatEth(campaign.goal)} goal</p>
          </div>

          <Button variant="secondary" onClick={() => setShowShareModal(true)} className="justify-center gap-1.5">
            <ShareIcon />
            Share
          </Button>

          <div className="flex items-center justify-between text-sm text-slate-500">
            <span>
              {contributors.length} contributor{contributors.length === 1 ? '' : 's'}
            </span>
            <StatusBadge status={campaign.status} />
          </div>

          <p className="text-sm text-slate-500">Ends {formatDate(campaign.deadline)}</p>

          <div>
            {isAuthenticated &&
              (canContribute ? (
                account ? (
                  <ContributeForm onContribute={handleContribute} isContributing={isContributing} />
                ) : (
                  <Button onClick={onConnectWallet}>Connect Wallet to Contribute</Button>
                )
              ) : (
                <p className="text-sm text-slate-500">This campaign is no longer accepting contributions.</p>
              ))}
          </div>

          {contributors.length > 0 && (
            <div className="flex flex-col gap-3 border-t border-slate-200 pt-4">
              <h3 className="text-sm font-medium text-slate-900">Recent Supporters</h3>
              {contributors.slice(0, 6).map((contributor) => (
                <div key={contributor.address} className="flex items-center gap-3">
                  <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-slate-100 text-xs font-medium text-slate-600">
                    {(contributor.displayName || contributor.address).charAt(0).toUpperCase()}
                  </div>
                  <div className="flex flex-1 items-center justify-between gap-2">
                    <span className="truncate text-sm text-slate-700">
                      {contributor.displayName || shortenAddress(contributor.address)}
                    </span>
                    <span className="shrink-0 text-sm text-slate-500">{formatEth(contributor.amount)}</span>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {showShareModal && (
        <ShareModal
          url={window.location.href}
          title={campaign.title}
          onClose={() => setShowShareModal(false)}
          showToast={showToast}
        />
      )}
    </div>
  )
}

export default CampaignDetailsPage
