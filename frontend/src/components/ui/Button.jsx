const VARIANTS = {
  primary: 'bg-indigo-600 text-white shadow-sm hover:bg-indigo-500',
  secondary: 'bg-white text-slate-700 border border-slate-300 hover:bg-slate-50',
  ghost: 'bg-transparent text-slate-600 hover:bg-slate-100',
  danger: 'bg-white text-rose-600 border border-rose-200 hover:bg-rose-50',
  link: 'bg-transparent text-indigo-600 hover:text-indigo-500 hover:underline p-0',
}

function Button({ variant = 'primary', className = '', ...props }) {
  const base =
    variant === 'link'
      ? 'inline-flex items-center text-sm font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed'
      : 'inline-flex items-center justify-center rounded-lg px-4 py-2 text-sm font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed'

  return <button className={`${base} ${VARIANTS[variant]} ${className}`} {...props} />
}

export default Button
