function Toast({ message, onDismiss }) {
  return (
    <div className="flex items-center gap-3 rounded-lg bg-slate-900 px-4 py-3 text-sm text-white shadow-lg">
      <span>{message}</span>
      <button type="button" onClick={onDismiss} aria-label="Dismiss" className="text-slate-400 hover:text-white">
        &times;
      </button>
    </div>
  )
}

export default Toast
