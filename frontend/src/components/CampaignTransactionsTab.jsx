import { useState } from 'react'
import { formatEther } from 'ethers'
import { shortenAddress, formatEth, formatDate } from '../utils/format'
import Button from './ui/Button'

const MOCK_TRANSACTIONS = [
  {
    id: 1,
    type: 'contribution',
    address: '0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984',
    amount: '1500000000000000000',
    txHash: '0xa1b2c3d4e5f60718293a4b5c6d7e8f9012345678901234567890abcdef12345',
    timestamp: Math.floor(Date.now() / 1000) - 2 * 24 * 60 * 60,
  },
  {
    id: 2,
    type: 'contribution',
    address: '0x5B38Da6a701c568545dCfcB03FcB875f56beddC4',
    amount: '500000000000000000',
    txHash: '0xb2c3d4e5f60718293a4b5c6d7e8f9012345678901234567890abcdef123456a',
    timestamp: Math.floor(Date.now() / 1000) - 1 * 24 * 60 * 60,
  },
  {
    id: 3,
    type: 'withdraw',
    address: '0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266',
    amount: '2000000000000000000',
    txHash: '0xc3d4e5f60718293a4b5c6d7e8f9012345678901234567890abcdef123456ab',
    timestamp: Math.floor(Date.now() / 1000) - 3 * 60 * 60,
  },
]

function TransactionTypeBadge({ type }) {
  const isWithdraw = type === 'withdraw'
  return (
    <span
      className={`inline-flex items-center rounded-full px-2.5 py-1 text-xs font-medium ${
        isWithdraw ? 'bg-indigo-50 text-indigo-600' : 'bg-emerald-50 text-emerald-600'
      }`}
    >
      {isWithdraw ? 'Withdraw' : 'Contribution'}
    </span>
  )
}

function CopyIcon() {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6" className="h-3.5 w-3.5">
      <rect x="9" y="9" width="11" height="11" rx="2" />
      <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
    </svg>
  )
}

function CheckIcon() {
  return (
    <svg viewBox="0 0 20 20" fill="currentColor" className="h-3.5 w-3.5 text-emerald-500">
      <path
        fillRule="evenodd"
        d="M16.7 5.3a1 1 0 0 1 0 1.4l-7.5 7.5a1 1 0 0 1-1.4 0l-3.5-3.5a1 1 0 1 1 1.4-1.4l2.8 2.8 6.8-6.8a1 1 0 0 1 1.4 0Z"
        clipRule="evenodd"
      />
    </svg>
  )
}

function DownloadIcon() {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6" className="h-4 w-4">
      <path strokeLinecap="round" strokeLinejoin="round" d="M12 3v12m0 0-4-4m4 4 4-4M4 17v2a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-2" />
    </svg>
  )
}

function downloadTransactionsCSV(transactions) {
  const header = ['Type', 'Address', 'Amount (ETH)', 'Date', 'Transaction Hash']
  const rows = transactions.map((tx) => [
    tx.type,
    tx.address,
    formatEther(tx.amount),
    new Date(tx.timestamp * 1000).toISOString(),
    tx.txHash,
  ])

  const csvContent = [header, ...rows].map((row) => row.join(',')).join('\n')
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)

  const link = document.createElement('a')
  link.href = url
  link.download = 'transactions.csv'
  link.click()

  URL.revokeObjectURL(url)
}

function CampaignTransactionsTab() {
  const transactions = MOCK_TRANSACTIONS
  const [copiedTxId, setCopiedTxId] = useState(null)

  function handleCopy(tx) {
    navigator.clipboard.writeText(tx.txHash)
    setCopiedTxId(tx.id)
    setTimeout(() => setCopiedTxId((current) => (current === tx.id ? null : current)), 1500)
  }

  return (
    <div className="flex flex-col gap-3">
      <div className="flex items-center justify-between">
        <h2 className="text-sm font-semibold text-slate-900">All Transactions</h2>
        <Button
          variant="secondary"
          onClick={() => downloadTransactionsCSV(transactions)}
          disabled={transactions.length === 0}
          className="gap-1.5"
        >
          <DownloadIcon />
          Download
        </Button>
      </div>

      {transactions.length === 0 ? (
        <div className="rounded-xl border border-dashed border-slate-300 bg-white px-6 py-16 text-center">
          <p className="text-sm text-slate-500">No transactions yet.</p>
        </div>
      ) : (
        <div className="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-slate-200 bg-slate-50">
                  <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wide text-slate-500">
                    Type
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wide text-slate-500">
                    Address
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wide text-slate-500">
                    Amount
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wide text-slate-500">
                    Date
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wide text-slate-500">
                    Transaction
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {transactions.map((tx) => (
                  <tr key={tx.id}>
                    <td className="px-4 py-3">
                      <TransactionTypeBadge type={tx.type} />
                    </td>
                    <td className="px-4 py-3 font-mono text-slate-900">{shortenAddress(tx.address)}</td>
                    <td className="px-4 py-3 text-slate-600">{formatEth(tx.amount)}</td>
                    <td className="px-4 py-3 text-slate-500">{formatDate(tx.timestamp)}</td>
                    <td className="px-4 py-3">
                      <div className="flex items-center gap-2">
                        <span className="font-mono text-xs text-indigo-600">{shortenAddress(tx.txHash)}</span>
                        <button
                          type="button"
                          onClick={() => handleCopy(tx)}
                          aria-label="Copy transaction hash"
                          className="cursor-pointer text-slate-400 hover:text-slate-600"
                        >
                          {copiedTxId === tx.id ? <CheckIcon /> : <CopyIcon />}
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  )
}

export default CampaignTransactionsTab
