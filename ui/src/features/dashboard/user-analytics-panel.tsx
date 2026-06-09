import { useMemo, useState } from 'react'
import { ArrowUpDown } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { formatCost, formatTokens } from '@/lib/format'
import { cn } from '@/lib/utils'
import type { BreakdownGroup, ProviderFilter, UsageBreakdownRow, UsageSummaryStats, UserSummary } from '@/types'
import { useUserBreakdownQuery, useUserUsageSummaryQuery } from './queries'

type SortKey = 'group' | 'total_tokens'
type SortState = { key: SortKey; direction: 'asc' | 'desc' }

const emptyBreakdown: UsageBreakdownRow[] = []
const emptyStats: UsageSummaryStats = {
  provider: 'all',
  total_projects: 0,
  total_input_tokens: 0,
  total_output_tokens: 0,
  total_cached_tokens: 0,
  total_tokens: 0,
  total_cost_usd: 0,
}

export function UserAnalyticsPanel({ user, apiKey, enabled, authVersion }: { user?: UserSummary; apiKey: string; enabled: boolean; authVersion: number }) {
  const [tab, setTab] = useState<'breakdown' | 'summary'>('breakdown')
  const [groupBy, setGroupBy] = useState<BreakdownGroup>('day')
  const [breakdownProvider, setBreakdownProvider] = useState<ProviderFilter>('all')
  const [provider, setProvider] = useState<ProviderFilter>('all')
  const [from, setFrom] = useState('2026-01-01')
  const [to, setTo] = useState('2026-12-31')
  const [sort, setSort] = useState<SortState>({ key: 'group', direction: 'asc' })

  const userID = user?.id ?? ''
  const breakdownQuery = useUserBreakdownQuery(userID, { groupBy, provider: breakdownProvider, from, to }, { apiKey, enabled, authVersion })
  const summaryQuery = useUserUsageSummaryQuery(userID, { provider, from, to }, { apiKey, enabled, authVersion })
  const rows = useMemo(() => sortRows(breakdownQuery.data ?? emptyBreakdown, sort), [breakdownQuery.data, sort])
  const stats = summaryQuery.data ?? emptyStats

  if (!user) {
    return <EmptyPanel />
  }

  return (
    <Card>
      <CardHeader className="flex-col items-stretch gap-3 md:flex-row md:items-center">
        <div className="min-w-0">
          <CardTitle>{user.name || user.email}</CardTitle>
          <p className="truncate text-xs text-muted-foreground">{user.email}</p>
        </div>
        <div className="flex flex-wrap gap-2">
          <Button type="button" variant={tab === 'breakdown' ? 'default' : 'outline'} onClick={() => setTab('breakdown')}>
            Breakdown
          </Button>
          <Button type="button" variant={tab === 'summary' ? 'default' : 'outline'} onClick={() => setTab('summary')}>
            Summary
          </Button>
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        {tab === 'breakdown' ? (
          <BreakdownTab
            groupBy={groupBy}
            setGroupBy={setGroupBy}
            provider={breakdownProvider}
            setProvider={setBreakdownProvider}
            from={from}
            setFrom={setFrom}
            to={to}
            setTo={setTo}
            rows={rows}
            sort={sort}
            setSort={setSort}
            isLoading={breakdownQuery.isFetching}
            error={breakdownQuery.error instanceof Error ? breakdownQuery.error : null}
          />
        ) : (
          <SummaryTab provider={provider} setProvider={setProvider} stats={stats} isLoading={summaryQuery.isFetching} error={summaryQuery.error instanceof Error ? summaryQuery.error : null} />
        )}
      </CardContent>
    </Card>
  )
}

function BreakdownTab({
  groupBy,
  setGroupBy,
  provider,
  setProvider,
  from,
  setFrom,
  to,
  setTo,
  rows,
  sort,
  setSort,
  isLoading,
  error,
}: {
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
}) {
  const groupLabel = groupBy === 'month' ? 'Month' : groupBy === 'project' ? 'Project' : 'Date'

  return (
    <div className="space-y-3">
      <div className="flex flex-col gap-2 lg:flex-row lg:items-center lg:justify-between">
        <Badge>{isLoading ? 'Loading' : `${rows.length} rows`}</Badge>
        <div className="grid gap-2 sm:grid-cols-[120px_120px_1fr_1fr]">
          <select className="h-9 rounded-md border border-input bg-background px-3 text-sm text-foreground" value={groupBy} onChange={(event) => setGroupBy(event.target.value as BreakdownGroup)}>
            <option value="day">Days</option>
            <option value="month">Months</option>
            <option value="project">Projects</option>
          </select>
          <select className="h-9 rounded-md border border-input bg-background px-3 text-sm text-foreground" value={provider} onChange={(event) => setProvider(event.target.value as ProviderFilter)}>
            <option value="all">All</option>
            <option value="codex">Codex</option>
            <option value="claude">Claude</option>
          </select>
          <Input type="date" value={from} onChange={(event) => setFrom(event.target.value)} />
          <Input type="date" value={to} onChange={(event) => setTo(event.target.value)} />
        </div>
      </div>
      {error && <p className="text-xs font-medium text-foreground">{error.message}</p>}
      <div className="overflow-x-auto rounded-md border border-border">
        <table className="w-full min-w-[980px] border-collapse text-sm">
          <thead className="bg-muted text-xs text-muted-foreground">
            <tr>
              {groupBy === 'project' ? (
                <PlainHead label={groupLabel} />
              ) : (
                <SortableHead label={groupLabel} sortKey="group" sort={sort} setSort={setSort} />
              )}
              <PlainHead label="Agent" />
              <PlainHead label="Models" />
              <PlainHead label="Input" align="right" />
              <PlainHead label="Output" align="right" />
              <PlainHead label="Cache Create" align="right" />
              <PlainHead label="Cache Read" align="right" />
              <SortableHead label="Total Tokens" sortKey="total_tokens" sort={sort} setSort={setSort} align="right" />
              <PlainHead label="Cost (USD)" align="right" />
            </tr>
          </thead>
          <tbody>
            {rows.map((row, index) => (
              <tr className="border-t border-border" key={`${row.group}-${row.agent}-${index}`}>
                <td className="px-3 py-2 font-medium text-foreground">{row.agent === 'all' ? row.group : ''}</td>
                <td className="px-3 py-2 text-foreground">{row.agent === 'all' ? 'All' : `- ${capitalize(row.agent)}`}</td>
                <td className="max-w-52 truncate px-3 py-2 text-muted-foreground">{row.models.map((model) => `- ${model}`).join(', ')}</td>
                <NumberCell value={row.input_tokens} />
                <NumberCell value={row.output_tokens} />
                <NumberCell value={row.cache_creation_tokens} />
                <NumberCell value={row.cache_read_tokens} />
                <NumberCell value={row.total_tokens} />
                <td className="px-3 py-2 text-right font-medium text-foreground">{formatCost(row.total_cost_usd)}</td>
              </tr>
            ))}
            {!rows.length && (
              <tr>
                <td className="px-3 py-8 text-center text-muted-foreground" colSpan={9}>
                  No usage loaded
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}

function SummaryTab({ provider, setProvider, stats, isLoading, error }: { provider: ProviderFilter; setProvider: (value: ProviderFilter) => void; stats: UsageSummaryStats; isLoading: boolean; error: Error | null }) {
  return (
    <div className="space-y-3">
      <div className="flex flex-wrap items-center justify-between gap-2">
        <div className="flex gap-2">
          {(['all', 'codex', 'claude'] as ProviderFilter[]).map((value) => (
            <Button key={value} type="button" variant={provider === value ? 'default' : 'outline'} onClick={() => setProvider(value)}>
              {capitalize(value)}
            </Button>
          ))}
        </div>
        <Badge>{isLoading ? 'Loading' : capitalize(stats.provider)}</Badge>
      </div>
      {error && <p className="text-xs font-medium text-foreground">{error.message}</p>}
      <div className="grid gap-3 md:grid-cols-3">
        <SummaryMetric label="Total Projects" value={formatTokens(stats.total_projects)} />
        <SummaryMetric label="Total Input" value={formatTokens(stats.total_input_tokens)} />
        <SummaryMetric label="Total Output" value={formatTokens(stats.total_output_tokens)} />
        <SummaryMetric label="Total Cached" value={formatTokens(stats.total_cached_tokens)} />
        <SummaryMetric label="Total Tokens" value={formatTokens(stats.total_tokens)} />
        <SummaryMetric label="Total Cost (USD)" value={formatCost(stats.total_cost_usd)} />
      </div>
    </div>
  )
}

function SortableHead({ label, sortKey, sort, setSort, align = 'left' }: { label: string; sortKey: SortKey; sort: SortState; setSort: (value: SortState) => void; align?: 'left' | 'right' }) {
  const active = sort.key === sortKey
  const nextDirection = active && sort.direction === 'asc' ? 'desc' : 'asc'

  return (
    <th className={cn('px-2 py-2 font-semibold', align === 'right' && 'text-right')}>
      <button className={cn('inline-flex items-center gap-1', align === 'right' && 'justify-end')} type="button" onClick={() => setSort({ key: sortKey, direction: nextDirection })}>
        {label}
        <ArrowUpDown className={cn('size-3', active ? 'text-foreground' : 'text-muted-foreground')} />
      </button>
    </th>
  )
}

function PlainHead({ label, align = 'left' }: { label: string; align?: 'left' | 'right' }) {
  return <th className={cn('px-3 py-2 font-semibold', align === 'right' && 'text-right')}>{label}</th>
}

function NumberCell({ value }: { value: number }) {
  return <td className="px-3 py-2 text-right tabular-nums text-foreground">{formatTokens(value)}</td>
}

function SummaryMetric({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-md border border-border bg-background p-4">
      <span className="text-xs text-muted-foreground">{label}</span>
      <strong className="mt-3 block text-xl font-semibold text-foreground">{value}</strong>
    </div>
  )
}

function EmptyPanel() {
  return (
    <Card>
      <CardContent>
        <div className="grid min-h-56 place-items-center text-sm text-muted-foreground">No user selected</div>
      </CardContent>
    </Card>
  )
}

function sortRows(rows: UsageBreakdownRow[], sort: SortState) {
  return [...rows].sort((left, right) => {
    const direction = sort.direction === 'asc' ? 1 : -1
    if (sort.key === 'group') {
      const groupCompare = compareValue(left.group, right.group)
      if (groupCompare !== 0) return groupCompare * direction
      return agentRank(left.agent) - agentRank(right.agent)
    }
    return compareValue(sortValue(left, sort.key), sortValue(right, sort.key)) * direction
  })
}

function sortValue(row: UsageBreakdownRow, key: SortKey) {
  return row[key]
}

function compareValue(left: string | number, right: string | number) {
  if (typeof left === 'number' && typeof right === 'number') return left - right
  return String(left).localeCompare(String(right))
}

function agentRank(agent: string) {
  if (agent === 'all') return 0
  if (agent === 'codex') return 1
  if (agent === 'claude') return 2
  return 3
}

function capitalize(value: string) {
  return value.charAt(0).toUpperCase() + value.slice(1)
}
