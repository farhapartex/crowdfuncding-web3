import { useState } from 'react'
import Modal from './Modal'

function XIcon() {
  return (
    <span className="flex h-12 w-12 items-center justify-center rounded-full bg-slate-900 text-lg font-bold text-white">
      X
    </span>
  )
}

function FacebookIcon() {
  return (
    <span className="flex h-12 w-12 items-center justify-center rounded-full bg-[#1877F2] text-lg font-bold text-white">
      f
    </span>
  )
}

function LinkedInIcon() {
  return (
    <span className="flex h-12 w-12 items-center justify-center rounded-full bg-[#0A66C2] text-sm font-bold text-white">
      in
    </span>
  )
}

function WhatsAppIcon() {
  return (
    <span className="flex h-12 w-12 items-center justify-center rounded-full bg-[#25D366] text-white">
      <svg viewBox="0 0 24 24" fill="currentColor" className="h-6 w-6">
        <path d="M12 2a10 10 0 0 0-8.6 15.1L2 22l5.05-1.36A10 10 0 1 0 12 2Zm0 18.2a8.16 8.16 0 0 1-4.16-1.14l-.3-.18-3 .81.8-2.93-.19-.3A8.2 8.2 0 1 1 20.2 12 8.21 8.21 0 0 1 12 20.2Zm4.5-6.13c-.24-.12-1.44-.71-1.66-.79s-.39-.12-.55.12-.63.79-.78.95-.28.18-.52.06a6.6 6.6 0 0 1-1.95-1.2 7.3 7.3 0 0 1-1.35-1.68c-.14-.24 0-.37.11-.49s.24-.28.35-.42a1.6 1.6 0 0 0 .24-.4.44.44 0 0 0 0-.42c-.06-.12-.55-1.32-.75-1.8s-.4-.41-.55-.41h-.47a.9.9 0 0 0-.65.3 2.73 2.73 0 0 0-.85 2 4.74 4.74 0 0 0 1 2.5 10.8 10.8 0 0 0 4.13 3.65c.58.25 1.03.4 1.38.51a3.32 3.32 0 0 0 1.53.1 2.5 2.5 0 0 0 1.64-1.16 2 2 0 0 0 .14-1.16c-.06-.11-.22-.17-.46-.29Z" />
      </svg>
    </span>
  )
}

function EmailIcon() {
  return (
    <span className="flex h-12 w-12 items-center justify-center rounded-full bg-slate-500 text-white">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className="h-5 w-5">
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M3.75 5.25h16.5v13.5H3.75V5.25Zm0 0 8.25 7.5 8.25-7.5"
        />
      </svg>
    </span>
  )
}

function CopyIcon() {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6" className="h-4 w-4">
      <rect x="9" y="9" width="11" height="11" rx="2" />
      <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
    </svg>
  )
}

function CheckIcon() {
  return (
    <svg viewBox="0 0 20 20" fill="currentColor" className="h-4 w-4">
      <path
        fillRule="evenodd"
        d="M16.7 5.3a1 1 0 0 1 0 1.4l-7.5 7.5a1 1 0 0 1-1.4 0l-3.5-3.5a1 1 0 1 1 1.4-1.4l2.8 2.8 6.8-6.8a1 1 0 0 1 1.4 0Z"
        clipRule="evenodd"
      />
    </svg>
  )
}

function buildShareTargets(url, title) {
  const encodedUrl = encodeURIComponent(url)
  const encodedTitle = encodeURIComponent(title)

  return [
    { key: 'x', label: 'X', Icon: XIcon, href: `https://twitter.com/intent/tweet?url=${encodedUrl}&text=${encodedTitle}` },
    { key: 'facebook', label: 'Facebook', Icon: FacebookIcon, href: `https://www.facebook.com/sharer/sharer.php?u=${encodedUrl}` },
    { key: 'whatsapp', label: 'WhatsApp', Icon: WhatsAppIcon, href: `https://wa.me/?text=${encodedTitle}%20${encodedUrl}` },
    { key: 'linkedin', label: 'LinkedIn', Icon: LinkedInIcon, href: `https://www.linkedin.com/sharing/share-offsite/?url=${encodedUrl}` },
    { key: 'email', label: 'Email', Icon: EmailIcon, href: `mailto:?subject=${encodedTitle}&body=${encodedUrl}` },
  ]
}

function ShareModal({ url, title, onClose, showToast }) {
  const [copied, setCopied] = useState(false)
  const targets = buildShareTargets(url, title)

  async function handleCopy() {
    try {
      await navigator.clipboard.writeText(url)
      setCopied(true)
      showToast?.('Link copied to clipboard!')
      setTimeout(() => setCopied(false), 1500)
    } catch {
      showToast?.('Could not copy link')
    }
  }

  return (
    <Modal title="Share this campaign" onClose={onClose}>
      <div className="grid grid-cols-5 gap-2">
        {targets.map(({ key, label, Icon, href }) => (
          <a
            key={key}
            href={href}
            target="_blank"
            rel="noopener noreferrer"
            className="flex flex-col items-center gap-2 rounded-xl p-2 text-center hover:bg-slate-50"
          >
            <Icon />
            <span className="text-xs font-medium text-slate-600">{label}</span>
          </a>
        ))}
      </div>

      <div className="mt-6 flex items-center gap-2 rounded-xl border border-slate-200 bg-slate-50 px-3 py-2.5">
        <input
          readOnly
          value={url}
          onFocus={(e) => e.target.select()}
          className="min-w-0 flex-1 truncate bg-transparent text-sm text-slate-600 outline-none"
        />
        <button
          type="button"
          onClick={handleCopy}
          className="flex shrink-0 cursor-pointer items-center gap-1.5 rounded-lg bg-indigo-600 px-3 py-1.5 text-xs font-medium text-white transition-colors hover:bg-indigo-500"
        >
          {copied ? <CheckIcon /> : <CopyIcon />}
          {copied ? 'Copied' : 'Copy link'}
        </button>
      </div>
    </Modal>
  )
}

export default ShareModal
