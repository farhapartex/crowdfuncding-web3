import { useState } from 'react'
import Modal from './Modal'
import Button from './ui/Button'

function ArchiveCampaignModal({ onCancel, onConfirm, isArchiving, error }) {
  const [note, setNote] = useState('')

  function handleSubmit(e) {
    e.preventDefault()
    if (!note.trim() || isArchiving) return
    onConfirm(note.trim())
  }

  return (
    <Modal title="Archive campaign?" onClose={onCancel}>
      <form onSubmit={handleSubmit} className="flex flex-col gap-3">
        <p className="text-sm text-slate-600">
          Archiving hides this campaign from public browsing. It will still be visible to anyone
          with the direct link, but no new donations, comments, or replies will be accepted.
          Please explain why you're archiving it before the deadline.
        </p>
        <textarea
          value={note}
          onChange={(e) => setNote(e.target.value)}
          placeholder="Why are you archiving this campaign?"
          rows={3}
          autoFocus
          required
          disabled={isArchiving}
          className="w-full rounded-xl border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 placeholder:text-slate-400 shadow-sm transition focus:outline-none focus:ring-2 focus:ring-indigo-500/40 focus:border-indigo-500"
        />
        {error && <p className="text-xs font-medium text-rose-500">{error}</p>}
        <div className="mt-2 flex justify-end gap-2">
          <Button type="button" variant="secondary" onClick={onCancel} disabled={isArchiving}>
            Cancel
          </Button>
          <Button type="submit" variant="danger" disabled={isArchiving || !note.trim()}>
            {isArchiving ? 'Archiving...' : 'Archive Campaign'}
          </Button>
        </div>
      </form>
    </Modal>
  )
}

export default ArchiveCampaignModal
