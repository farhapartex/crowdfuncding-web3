const CATEGORIES = [
  'All Categories',
  'Technology',
  'Education',
  'Health & Medical',
  'Community',
  'Arts & Creative',
  'Charity',
  'Other',
]

function SearchIcon() {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className="h-4.5 w-4.5">
      <circle cx="11" cy="11" r="7" />
      <path strokeLinecap="round" d="m20 20-3.2-3.2" />
    </svg>
  )
}

function ChevronIcon() {
  return (
    <svg viewBox="0 0 20 20" fill="currentColor" className="h-4 w-4">
      <path
        fillRule="evenodd"
        d="M5.3 7.3a1 1 0 0 1 1.4 0L10 10.6l3.3-3.3a1 1 0 1 1 1.4 1.4l-4 4a1 1 0 0 1-1.4 0l-4-4a1 1 0 0 1 0-1.4Z"
        clipRule="evenodd"
      />
    </svg>
  )
}

function CampaignSearchBar({ search, onSearchChange, category, onCategoryChange }) {
  return (
    <div className="mx-auto flex w-full max-w-3xl flex-col gap-3 sm:flex-row">
      <div className="relative flex-1">
        <span className="pointer-events-none absolute inset-y-0 left-4 flex items-center text-slate-400">
          <SearchIcon />
        </span>
        <input
          type="text"
          value={search}
          onChange={(e) => onSearchChange(e.target.value)}
          placeholder="Search campaigns..."
          className="w-full rounded-xl border border-slate-200 bg-white py-3 pl-11 pr-4 text-sm text-slate-900 placeholder:text-slate-400 shadow-sm transition focus:outline-none focus:ring-2 focus:ring-indigo-500/40 focus:border-indigo-500"
        />
      </div>

      <div className="relative sm:w-56">
        <select
          value={category}
          onChange={(e) => onCategoryChange(e.target.value)}
          className="w-full appearance-none rounded-xl border border-slate-200 bg-white py-3 pl-4 pr-10 text-sm font-medium text-slate-700 shadow-sm transition focus:outline-none focus:ring-2 focus:ring-indigo-500/40 focus:border-indigo-500"
        >
          {CATEGORIES.map((option) => (
            <option key={option} value={option}>
              {option}
            </option>
          ))}
        </select>
        <span className="pointer-events-none absolute inset-y-0 right-3.5 flex items-center text-slate-400">
          <ChevronIcon />
        </span>
      </div>
    </div>
  )
}

export default CampaignSearchBar
