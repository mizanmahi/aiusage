import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { formatCost, formatTokens } from '@/lib/format'
import type { ProviderFilter, UsageSummaryStats } from '@/types'
import { capitalize } from '../lib/dashboard-format'

export function SummaryTab({ provider, setProvider, stats, isLoading, error }: { provider: ProviderFilter; setProvider: (value: ProviderFilter) => void; stats: UsageSummaryStats; isLoading: boolean; error: Error | null }) {
  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap items-center justify-between gap-2">
        <div className="flex gap-2">
          {(['all', 'codex', 'claude'] as ProviderFilter[]).map((value) => (
            <Button key={value} type="button" variant={provider === value ? 'default' : 'outline'} onClick={() => setProvider(value)}>
              {capitalize(value)}
            </Button>
          ))}
        </div>
      </div>
      {error && <p className="text-xs font-medium text-foreground">{error.message}</p>}
      <div className="grid gap-3 md:grid-cols-3">
        <SummaryMetric isLoading={isLoading} label="Total Projects" value={formatTokens(stats.total_projects)} />
        <SummaryMetric isLoading={isLoading} label="Total Input" value={formatTokens(stats.total_input_tokens)} />
        <SummaryMetric isLoading={isLoading} label="Total Output" value={formatTokens(stats.total_output_tokens)} />
        <SummaryMetric isLoading={isLoading} label="Total Cached" value={formatTokens(stats.total_cached_tokens)} />
        <SummaryMetric isLoading={isLoading} label="Total Tokens" value={formatTokens(stats.total_tokens)} />
        <SummaryMetric isLoading={isLoading} label="Total Cost" value={formatCost(stats.total_cost_usd)} />
      </div>
    </div>
  )
}

function SummaryMetric({ label, value, isLoading }: { label: string; value: string; isLoading: boolean }) {
  return (
    <div className="rounded-md border border-border bg-background p-4">
      <span className="text-xs text-muted-foreground">{label}</span>
      {isLoading ? <Skeleton className="mt-3 h-7 w-24" /> : <strong className="mt-3 block text-xl font-semibold text-foreground">{value}</strong>}
    </div>
  )
}
