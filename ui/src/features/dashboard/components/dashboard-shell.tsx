import { type ReactNode } from 'react'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { cn } from '@/lib/utils'
import { DashboardNav } from './dashboard-nav'
import { PagePatternBackground } from './page-pattern-background'

type DashboardShellProps = {
  eyebrow: string
  title: string
  description: string
  navActions: ReactNode
  isPowered: boolean
  children: ReactNode
}

export function DashboardShell({ eyebrow, title, description, navActions, isPowered, children }: DashboardShellProps) {
  return (
    <main className={cn('power-shell relative isolate min-h-screen overflow-hidden text-foreground', isPowered && 'is-powered')}>
      <PagePatternBackground />
      <div className="relative z-10 mx-auto flex w-full max-w-[1180px] flex-col gap-5 px-4 py-4 md:px-6 md:py-6">
        <DashboardNav actions={navActions} />
        <header className="power-core rounded-lg border border-border bg-card/95 shadow-sm backdrop-blur">
          <div className="flex min-h-56 flex-col items-start justify-center gap-4 p-5 md:p-8">
            <div className="min-w-0">
              <Badge>{eyebrow}</Badge>
              <h1 className="mt-3 text-3xl font-semibold leading-tight text-foreground md:text-5xl">{title}</h1>
              <p className="mt-3 max-w-2xl text-sm text-muted-foreground md:text-base">{description}</p>
            </div>
          </div>
          <Separator />
          <div className="flex flex-wrap gap-2 px-4 py-3">
            {['codex', 'claude', 'tokens', 'cost', 'projects'].map((label) => (
              <Badge className="bg-background" key={label}>
                {label}
              </Badge>
            ))}
          </div>
        </header>
        <div className="power-content flex flex-col gap-5">{children}</div>
      </div>
    </main>
  )
}
