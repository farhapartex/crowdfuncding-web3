import { useEffect, useRef } from 'react'
import { useAuth0 } from '@auth0/auth0-react'
import { syncAuth0User } from '../lib/api'
import Button from './ui/Button'

function Auth0AccountButton() {
  const { isAuthenticated, isLoading, user, loginWithRedirect, logout, getAccessTokenSilently } = useAuth0()
  const hasSynced = useRef(false)

  useEffect(() => {
    if (!isAuthenticated || hasSynced.current) return
    hasSynced.current = true

    getAccessTokenSilently()
      .then((accessToken) => syncAuth0User(accessToken))
      .catch((err) => {
        console.error('failed to sync auth0 user', err)
        hasSynced.current = false
      })
  }, [isAuthenticated, getAccessTokenSilently])

  if (isLoading) {
    return <span className="text-sm text-slate-400">Loading...</span>
  }

  if (isAuthenticated) {
    return (
      <div className="flex items-center gap-3">
        <span className="text-sm font-medium text-slate-700">{user?.name || user?.email}</span>
        <Button
          variant="secondary"
          onClick={() => logout({ logoutParams: { returnTo: window.location.origin } })}
        >
          Log out
        </Button>
      </div>
    )
  }

  return (
    <Button onClick={() => loginWithRedirect({ appState: { returnTo: window.location.pathname } })}>
      Log in
    </Button>
  )
}

export default Auth0AccountButton
