import type { LucideIcon } from 'lucide-react'
import { Card, CardContent } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'

type Metric = {
  icon: LucideIcon
  label: string
  value: string
  detail: string
}

export function MetricGrid({ metrics, isLoading }: { metrics: Metric[]; isLoading: boolean }) {
  return (
    <section className="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
      {metrics.map((metric) => (
        <MetricCard isLoading={isLoading} key={metric.label} metric={metric} />
      ))}
    </section>
  )
}

function MetricCard({ metric, isLoading }: { metric: Metric; isLoading: boolean }) {
  const Icon = metric.icon

  return (
    <Card>
      <CardContent className="p-4">
        <div className="flex items-start justify-between gap-3">
          <div className="min-w-0">
            <p className="text-sm font-medium text-muted-foreground">{metric.label}</p>
            {isLoading ? <Skeleton className="mt-3 h-8 w-28" /> : <strong className="mt-2 block text-2xl font-semibold text-foreground">{metric.value}</strong>}
          </div>
          <div className="grid size-10 place-items-center rounded-md border border-border bg-muted text-muted-foreground">
            <Icon />
          </div>
        </div>
        {isLoading ? <Skeleton className="mt-4 h-4 w-36" /> : <p className="mt-3 truncate text-xs text-muted-foreground">{metric.detail}</p>}
      </CardContent>
    </Card>
  )
}
