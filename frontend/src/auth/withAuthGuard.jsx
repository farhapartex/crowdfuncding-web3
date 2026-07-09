import { withAuthenticationRequired } from '@auth0/auth0-react'

function CheckingSession() {
  return (
    <div className="flex justify-center py-16">
      <p className="text-sm text-slate-500">Checking your session...</p>
    </div>
  )
}

function withAuthGuard(Component) {
  return withAuthenticationRequired(Component, {
    onRedirecting: CheckingSession,
  })
}

export default withAuthGuard
