function Pagination({ offset, pageSize, total, onPrevious, onNext }) {
  if (total === 0) {
    return null
  }

  const from = offset + 1
  const to = Math.min(offset + pageSize, total)

  return (
    <div className="pagination">
      <span className="pagination-info">
        {from}-{to} of {total}
      </span>
      <div className="pagination-controls">
        <button type="button" onClick={onPrevious} disabled={offset === 0}>
          Previous
        </button>
        <button type="button" onClick={onNext} disabled={offset + pageSize >= total}>
          Next
        </button>
      </div>
    </div>
  )
}

export default Pagination
