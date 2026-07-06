import { useEffect, useRef, useState } from 'react'
import { Link } from 'react-router-dom'

function ProfileIcon() {
  return (
    <svg viewBox="0 0 24 24" fill="currentColor" className="h-5 w-5">
      <path d="M12 12a5 5 0 1 0 0-10 5 5 0 0 0 0 10Zm0 2c-4.42 0-8 2.24-8 5v1a1 1 0 0 0 1 1h14a1 1 0 0 0 1-1v-1c0-2.76-3.58-5-8-5Z" />
    </svg>
  )
}

function AccountMenu({ onSignOut }) {
  const [isOpen, setIsOpen] = useState(false)
  const menuRef = useRef(null)

  useEffect(() => {
    function handleClickOutside(e) {
      if (menuRef.current && !menuRef.current.contains(e.target)) {
        setIsOpen(false)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  return (
    <div className="relative" ref={menuRef}>
      <button
        type="button"
        onClick={() => setIsOpen((open) => !open)}
        aria-label="Account menu"
        className="flex h-9 w-9 items-center justify-center rounded-full bg-indigo-600 text-white hover:bg-indigo-500"
      >
        <ProfileIcon />
      </button>

      {isOpen && (
        <div className="absolute right-0 top-[calc(100%+8px)] z-50 flex w-44 flex-col gap-1 rounded-xl border border-slate-200 bg-white p-2 shadow-lg">
          <Link
            to="/profile"
            onClick={() => setIsOpen(false)}
            className="rounded-lg px-3 py-2 text-sm text-slate-700 hover:bg-slate-100"
          >
            Profile
          </Link>
          <Link
            to="/my-campaigns"
            onClick={() => setIsOpen(false)}
            className="rounded-lg px-3 py-2 text-sm text-slate-700 hover:bg-slate-100"
          >
            Campaign
          </Link>
          <button
            type="button"
            onClick={() => {
              onSignOut()
              setIsOpen(false)
            }}
            className="rounded-lg px-3 py-2 text-left text-sm text-slate-700 hover:bg-slate-100"
          >
            Logout
          </button>
        </div>
      )}
    </div>
  )
}

export default AccountMenu
