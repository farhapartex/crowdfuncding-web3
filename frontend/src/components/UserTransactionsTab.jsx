import { useEffect, useState } from 'react'
import { formatEther } from 'ethers'
import { fetchWalletTransactions } from '../lib/api'
import { shortenAddress, formatEth } from '../utils/format'
import Button from './ui/Button'
import Pagination from './Pagination'

const PAGE_SIZE = 10

function TransactionTypeBadge({ type }) {
  const styles = {
    withdraw: 'bg-indigo-50 text-indigo-600',
    refund: 'bg-amber-50 text-amber-600',
    contribution: 'bg-emerald-50 text-emerald-600',
  }
  const labels = {
    withdraw: 'Withdraw',
    refund: 'Refund',
    contribution: 'Contribution',
  }

  return (
    <span className={`inline-flex items-center rounded-full px-2.5 py-1 text-xs font-medium ${styles[type] || 'bg-slate-100 text-slate-600'}`}>
      {labels[type] || type}
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
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M12 3v12m0 0-4-4m4 4 4-4M4 17v2a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-2"
      />
    </svg>
  )
}

function downloadTransactionsCSV(transactions) {
  const header = ['Type', 'Campaign', 'Amount (ETH)', 'Date', 'Transaction Hash']
  const rows = transactions.map((tx) => [
    tx.type,
    `Campaign #${tx.campaignId}`,
    formatEther(tx.amount),
    new Date(tx.blockTimestamp).toISOString(),
    tx.txHash,
  ])

  const csvContent = [header, ...rows].map((row) => row.join(',')).join('\n')
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)

  const link = document.createElement('a')
  link.href = url
  link.download = 'my-transactions.csv'
  link.click()

  URL.revokeObjectURL(url)
}

function UserTransactionsTab({ address }) {
  const [transactions, setTransactions] = useState([])
  const [total, setTotal] = useState(0)
  const [offset, setOffset] = useState(0)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState(null)
  const [copiedTxId, setCopiedTxId] = useState(null)

  useEffect(() => {
    if (address) refresh(0)
  }, [address])

  async function refresh(targetOffset) {
    setIsLoading(true)
    setError(null)
    try {
      const { items, total: totalCount } = await fetchWalletTransactions(address, {
        offset: targetOffset,
        limit: PAGE_SIZE,
      })
      setTransactions(items)
      setTotal(totalCount)
      setOffset(targetOffset)
    } catch (err) {
      setError(err.message)
    } finally {
      setIsLoading(false)
    }
  }

  function handleCopy(tx) {
    navigator.clipboard.writeText(tx.txHash)
    setCopiedTxId(tx.id)
    setTimeout(() => setCopiedTxId((current) => (current === tx.id ? null : current)), 1500)
  }

  if (!address) {
    return (
      <div className="rounded-xl border border-dashed border-slate-300 bg-white px-6 py-16 text-center">
        <p className="text-sm text-slate-500">Connect your wallet to see your transactions.</p>
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-3">
      <div className="flex items-center justify-between">
        <h2 className="text-sm font-semibold text-slate-900">My Transactions</h2>
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

      {error && <p className="text-sm text-rose-600">{error}</p>}

      {isLoading ? (
        <p className="text-sm text-slate-500">Loading transactions...</p>
      ) : transactions.length === 0 ? (
        <div className="rounded-xl border border-dashed border-slate-300 bg-white px-6 py-16 text-center">
          <p className="text-sm text-slate-500">You haven't made any transactions yet.</p>
        </div>
      ) : (
        <>
          <div className="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm">
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-slate-200 bg-slate-50">
                    <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wide text-slate-500">
                      Type
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wide text-slate-500">
                      Campaign
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
                      <td className="px-4 py-3 text-slate-900">Campaign #{tx.campaignId}</td>
                      <td className="px-4 py-3 text-slate-600">{formatEth(tx.amount)}</td>
                      <td className="px-4 py-3 text-slate-500">{new Date(tx.blockTimestamp).toLocaleString()}</td>
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

          <Pagination
            offset={offset}
            pageSize={PAGE_SIZE}
            total={total}
            onPrevious={() => refresh(Math.max(0, offset - PAGE_SIZE))}
            onNext={() => refresh(offset + PAGE_SIZE)}
          />
        </>
      )}
    </div>
  )
}

export default UserTransactionsTab
