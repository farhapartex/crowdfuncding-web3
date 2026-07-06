import Button from './ui/Button'

function Pagination({ offset, pageSize, total, onPrevious, onNext }) {
  if (total === 0) {
    return null
  }

  const from = offset + 1
  const to = Math.min(offset + pageSize, total)

  return (
    <div className="flex items-center justify-between px-1">
      <span className="text-sm text-slate-500">
        {from}-{to} of {total}
      </span>
      <div className="flex gap-2">
        <Button variant="secondary" onClick={onPrevious} disabled={offset === 0}>
          Previous
        </Button>
        <Button variant="secondary" onClick={onNext} disabled={offset + pageSize >= total}>
          Next
        </Button>
      </div>
    </div>
  )
}

export default Pagination
