const MINUTE = 60
const HOUR = 60 * MINUTE
const DAY = 24 * HOUR

export function formatCommentTimestamp(unixSeconds) {
  const date = new Date(unixSeconds * 1000)
  const ageSeconds = Math.floor(Date.now() / 1000) - unixSeconds

  if (ageSeconds >= DAY) {
    return (
      date.toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' }) +
      ' at ' +
      date.toLocaleTimeString(undefined, { hour: 'numeric', minute: '2-digit' })
    )
  }
  if (ageSeconds >= HOUR) {
    const hours = Math.floor(ageSeconds / HOUR)
    return `${hours} hour${hours === 1 ? '' : 's'} ago`
  }
  if (ageSeconds >= MINUTE) {
    const minutes = Math.floor(ageSeconds / MINUTE)
    return `${minutes} minute${minutes === 1 ? '' : 's'} ago`
  }
  return 'Just now'
}

export function groupComments(comments) {
  const rootComments = comments.filter((comment) => !comment.parentId)
  const repliesByParent = comments.reduce((acc, comment) => {
    if (!comment.parentId) return acc
    acc[comment.parentId] = acc[comment.parentId] || []
    acc[comment.parentId].push(comment)
    return acc
  }, {})
  Object.values(repliesByParent).forEach((replies) => replies.sort((a, b) => a.createdAt - b.createdAt))

  return { rootComments, repliesByParent }
}
