import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'
import type { BreakdownGroup, ProviderFilter, UsageBreakdownRow } from '@/types'
import type { SortState } from '../breakdown-utils'
import { BreakdownTable } from './breakdown-table'

type BreakdownTabProps = {
  groupBy: BreakdownGroup
  setGroupBy: (value: BreakdownGroup) => void
  provider: ProviderFilter
  setProvider: (value: ProviderFilter) => void
  from: string
  setFrom: (value: string) => void
  to: string
  setTo: (value: string) => void
  rows: UsageBreakdownRow[]
  sort: SortState
  setSort: (value: SortState) => void
  isLoading: boolean
  error: Error | null
}

export function BreakdownTab(props: BreakdownTabProps) {
  const { groupBy, setGroupBy, provider, setProvider, from, setFrom, to, setTo, rows, sort, setSort, isLoading, error } = props

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-col gap-3 xl:flex-row xl:items-center xl:justify-between">
        <Badge>{isLoading ? 'Loading usage' : `${rows.length} rows`}</Badge>
        <div className={cn('grid gap-2', groupBy === 'month' ? 'sm:grid-cols-[140px_140px]' : 'sm:grid-cols-[140px_140px_1fr_1fr]')}>
          <select className="h-9 rounded-md border border-input bg-background px-3 text-sm text-foreground" value={groupBy} onChange={(event) => setGroupBy(event.target.value as BreakdownGroup)}>
            <option value="day">Days</option>
            <option value="month">Months</option>
            <option value="project">Projects</option>
          </select>
          <select className="h-9 rounded-md border border-input bg-background px-3 text-sm text-foreground" value={provider} onChange={(event) => setProvider(event.target.value as ProviderFilter)}>
            <option value="all">All tools</option>
            <option value="codex">Codex</option>
            <option value="claude">Claude</option>
          </select>
          {groupBy !== 'month' && (
            <>
              <Input type="date" value={from} onChange={(event) => setFrom(event.target.value)} />
              <Input type="date" value={to} onChange={(event) => setTo(event.target.value)} />
            </>
          )}
        </div>
      </div>
      {error && <p className="text-xs font-medium text-foreground">{error.message}</p>}
      <BreakdownTable groupBy={groupBy} rows={rows} sort={sort} setSort={setSort} isLoading={isLoading} />
    </div>
  )
}
