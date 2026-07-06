function AuthStatus({ sessionAddress, isSigningIn, onSignIn, onSignOut, onEditProfile }) {
  if (sessionAddress) {
    return (
      <>
        <button type="button" onClick={onEditProfile}>
          Edit Profile
        </button>
        <button type="button" onClick={onSignOut}>
          Sign Out
        </button>
      </>
    )
  }

  return (
    <button type="button" onClick={onSignIn} disabled={isSigningIn}>
      {isSigningIn ? 'Signing in...' : 'Sign In'}
    </button>
  )
}

export default AuthStatus
