import { Auth0Provider } from '@auth0/auth0-react'
import { useNavigate } from 'react-router-dom'

function Auth0ProviderWithNavigate({ children }) {
  const navigate = useNavigate()

  const domain = import.meta.env.VITE_AUTH0_DOMAIN
  const clientId = import.meta.env.VITE_AUTH0_CLIENT_ID
  const audience = import.meta.env.VITE_AUTH0_AUDIENCE

  function onRedirectCallback(appState) {
    navigate(appState?.returnTo || '/my-campaigns')
  }

  return (
    <Auth0Provider
      domain={domain}
      clientId={clientId}
      authorizationParams={{
        redirect_uri: `${window.location.origin}/auth/callback`,
        audience,
        scope: 'openid profile email offline_access',
      }}
      useRefreshTokens
      cacheLocation="localstorage"
      onRedirectCallback={onRedirectCallback}
    >
      {children}
    </Auth0Provider>
  )
}

export default Auth0ProviderWithNavigate
