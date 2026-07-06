const COLORS = {
  Active: 'bg-indigo-50 text-indigo-700 ring-1 ring-inset ring-indigo-200',
  Successful: 'bg-emerald-50 text-emerald-700 ring-1 ring-inset ring-emerald-200',
  Failed: 'bg-rose-50 text-rose-700 ring-1 ring-inset ring-rose-200',
}

function StatusBadge({ status }) {
  return (
    <span className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${COLORS[status]}`}>
      {status}
    </span>
  )
}

export default StatusBadge
