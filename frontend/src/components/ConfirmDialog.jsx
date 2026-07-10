import Modal from './Modal'
import Button from './ui/Button'

function ConfirmDialog({ title, message, confirmLabel = 'Confirm', isConfirming, onConfirm, onCancel }) {
  return (
    <Modal title={title} onClose={onCancel}>
      <p className="text-sm text-slate-600">{message}</p>

      <div className="mt-6 flex justify-end gap-2">
        <Button variant="secondary" onClick={onCancel} disabled={isConfirming}>
          Cancel
        </Button>
        <Button variant="danger" onClick={onConfirm} disabled={isConfirming}>
          {isConfirming ? 'Deleting...' : confirmLabel}
        </Button>
      </div>
    </Modal>
  )
}

export default ConfirmDialog
