import { useMemo, useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import type { BreakdownGroup, ProviderFilter, UsageBreakdownRow, UsageSummaryStats, UserSummary } from '@/types'
import { rowsForProvider, sortRows, type SortState } from './breakdown-utils'
import { BreakdownTab } from './components/breakdown-tab'
import { EmptyState } from './components/empty-state'
import { PanelTabs } from './components/panel-tabs'
import { SummaryTab } from './components/summary-tab'
import { useUserBreakdownQuery, useUserUsageSummaryQuery } from './queries'

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
  const breakdownFrom = groupBy === 'month' ? '' : from
  const breakdownTo = groupBy === 'month' ? '' : to
  const breakdownQuery = useUserBreakdownQuery(userID, { groupBy, provider: breakdownProvider, from: breakdownFrom, to: breakdownTo }, { apiKey, enabled, authVersion })
  const summaryQuery = useUserUsageSummaryQuery(userID, { provider, from, to }, { apiKey, enabled, authVersion })
  const rows = useMemo(() => sortRows(rowsForProvider(breakdownQuery.data ?? emptyBreakdown, breakdownProvider), sort), [breakdownQuery.data, breakdownProvider, sort])

  if (!user) {
    return <EmptyAnalyticsPanel />
  }

  return (
    <Card>
      <CardHeader className="flex-col items-stretch gap-3 md:flex-row md:items-center">
        <div className="min-w-0">
          <CardTitle>{user.name || user.email}</CardTitle>
          <p className="truncate text-xs text-muted-foreground">{user.email}</p>
        </div>
        <PanelTabs value={tab} onChange={setTab} />
      </CardHeader>
      <CardContent className="p-4">
        {tab === 'breakdown' ? (
          <BreakdownTab
            groupBy={groupBy}
            setGroupBy={(value) => {
              setGroupBy(value)
              setSort(value === 'project' ? { key: 'total_tokens', direction: 'desc' } : { key: 'group', direction: 'asc' })
            }}
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
          <SummaryTab provider={provider} setProvider={setProvider} stats={summaryQuery.data ?? emptyStats} isLoading={summaryQuery.isFetching} error={summaryQuery.error instanceof Error ? summaryQuery.error : null} />
        )}
      </CardContent>
    </Card>
  )
}

function EmptyAnalyticsPanel() {
  return (
    <Card>
      <CardContent className="p-4">
        <EmptyState title="No developer selected" description="Load users and choose someone to inspect their usage breakdown." />
      </CardContent>
    </Card>
  )
}
