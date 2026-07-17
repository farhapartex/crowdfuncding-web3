import { useCallback } from 'react'
import { useAuth0 } from '@auth0/auth0-react'

export function useAccessToken() {
  const { getAccessTokenSilently, loginWithRedirect } = useAuth0()

  return useCallback(async () => {
    try {
      return await getAccessTokenSilently()
    } catch (error) {
      await loginWithRedirect({
        appState: { returnTo: `${window.location.pathname}${window.location.search}` },
      })
      throw error
    }
  }, [getAccessTokenSilently, loginWithRedirect])
}
