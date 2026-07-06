import Toast from './Toast'

function ToastContainer({ toasts, onDismiss }) {
  if (toasts.length === 0) {
    return null
  }

  return (
    <div className="fixed bottom-4 left-4 z-50 flex flex-col gap-2">
      {toasts.map((toast) => (
        <Toast key={toast.id} message={toast.message} onDismiss={() => onDismiss(toast.id)} />
      ))}
    </div>
  )
}

export default ToastContainer
