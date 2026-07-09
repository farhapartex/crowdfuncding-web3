import { useEffect, useRef, useState } from 'react'
import { Link } from 'react-router-dom'
import { useAuth0 } from '@auth0/auth0-react'
import { useCurrentUser } from '../auth/CurrentUserContext'
import Button from './ui/Button'

function AccountMenu() {
  const { isAuthenticated, isLoading, user, loginWithRedirect, logout } = useAuth0()
  const { currentUser } = useCurrentUser()
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

  if (isLoading) {
    return <span className="text-sm text-slate-400">Loading...</span>
  }

  if (!isAuthenticated) {
    return (
      <Button onClick={() => loginWithRedirect({ appState: { returnTo: '/my-campaigns' } })}>
        Log in
      </Button>
    )
  }

  const label = currentUser?.displayName || currentUser?.email || user?.name || user?.email || 'Account'
  const initial = label.charAt(0).toUpperCase()

  return (
    <div className="relative" ref={menuRef}>
      <button
        type="button"
        onClick={() => setIsOpen((open) => !open)}
        aria-label="Account menu"
        className="flex h-9 w-9 cursor-pointer items-center justify-center rounded-full bg-indigo-600 text-sm font-semibold text-white hover:bg-indigo-500"
      >
        {initial}
      </button>

      {isOpen && (
        <div className="absolute right-0 top-[calc(100%+8px)] z-50 flex w-48 flex-col gap-1 rounded-xl border border-slate-200 bg-white p-2 shadow-lg">
          <div className="truncate px-3 py-1.5 text-xs font-medium text-slate-400">{label}</div>

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
              setIsOpen(false)
              logout({ logoutParams: { returnTo: window.location.origin } })
            }}
            className="cursor-pointer rounded-lg px-3 py-2 text-left text-sm text-slate-700 hover:bg-slate-100"
          >
            Log out
          </button>
        </div>
      )}
    </div>
  )
}

export default AccountMenu
