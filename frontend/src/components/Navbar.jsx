import { NavLink } from 'react-router-dom'
import { useAuth0 } from '@auth0/auth0-react'
import ConnectWalletButton from './ConnectWalletButton'
import AccountMenu from './AccountMenu'
import { useCurrentUser } from '../auth/CurrentUserContext'

function Navbar({ account, onConnect }) {
  const { isAuthenticated } = useAuth0()
  const { currentUser } = useCurrentUser()
  const displayName = currentUser?.displayName || currentUser?.email
  return (
    <nav className="sticky top-0 z-40 border-b border-slate-200 bg-white/80 backdrop-blur-sm">
      <div className="mx-auto flex h-16 max-w-5xl items-center gap-8 px-6">
        <NavLink to="/" className="text-base font-semibold text-slate-900">
          🌱 Crowd Funding
        </NavLink>

        <div className="flex flex-1 gap-6">
          <NavLink
            to="/"
            end
            className={({ isActive }) =>
              `text-sm font-medium transition-colors ${isActive ? 'text-indigo-600' : 'text-slate-600 hover:text-slate-900'}`
            }
          >
            Home
          </NavLink>
          <NavLink
            to="/campaigns"
            className={({ isActive }) =>
              `text-sm font-medium transition-colors ${isActive ? 'text-indigo-600' : 'text-slate-600 hover:text-slate-900'}`
            }
          >
            Campaigns
          </NavLink>
          <NavLink
            to="/about"
            className={({ isActive }) =>
              `text-sm font-medium transition-colors ${isActive ? 'text-indigo-600' : 'text-slate-600 hover:text-slate-900'}`
            }
          >
            About
          </NavLink>
        </div>

        <div className="flex items-center gap-4">
          {isAuthenticated &&
            (account ? (
              <span className="hidden text-sm font-medium text-slate-700 sm:inline">{displayName}</span>
            ) : (
              <ConnectWalletButton onConnect={onConnect} />
            ))}

          <AccountMenu />
        </div>
      </div>
    </nav>
  )
}

export default Navbar
