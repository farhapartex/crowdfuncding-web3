import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import './index.css'
import App from './App.jsx'
import Auth0ProviderWithNavigate from './auth/Auth0ProviderWithNavigate.jsx'
import { CurrentUserProvider } from './auth/CurrentUserContext.jsx'

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <BrowserRouter>
      <Auth0ProviderWithNavigate>
        <CurrentUserProvider>
          <App />
        </CurrentUserProvider>
      </Auth0ProviderWithNavigate>
    </BrowserRouter>
  </StrictMode>,
)
