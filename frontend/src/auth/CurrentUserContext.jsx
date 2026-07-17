import { createContext, useCallback, useContext, useEffect, useState } from 'react'
import { useAuth0 } from '@auth0/auth0-react'
import { fetchAuth0Me, syncAuth0User } from '../lib/api'
import { useAccessToken } from '../hooks/useAccessToken'

const CurrentUserContext = createContext({
  currentUser: null,
  isLoading: false,
  error: null,
  refresh: () => {},
})

async function loadCurrentUser(accessToken) {
  try {
    return await fetchAuth0Me(accessToken)
  } catch (err) {
    if (err.status !== 404) throw err
    return syncAuth0User(accessToken)
  }
}

export function CurrentUserProvider({ children }) {
  const { isAuthenticated } = useAuth0()
  const getAccessToken = useAccessToken()
  const [currentUser, setCurrentUser] = useState(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState(null)

  const refresh = useCallback(async () => {
    if (!isAuthenticated) {
      setCurrentUser(null)
      return
    }

    setIsLoading(true)
    setError(null)
    try {
      const accessToken = await getAccessToken()
      const me = await loadCurrentUser(accessToken)
      setCurrentUser(me)
    } catch (err) {
      setError(err.message)
    } finally {
      setIsLoading(false)
    }
  }, [isAuthenticated, getAccessToken])

  useEffect(() => {
    refresh()
  }, [refresh])

  return (
    <CurrentUserContext.Provider value={{ currentUser, isLoading, error, refresh }}>
      {children}
    </CurrentUserContext.Provider>
  )
}

export function useCurrentUser() {
  return useContext(CurrentUserContext)
}
