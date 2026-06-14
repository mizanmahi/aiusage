import { type ReactNode } from 'react'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'

type DashboardShellProps = {
  eyebrow: string
  title: string
  description: string
  actions: ReactNode
  children: ReactNode
}

export function DashboardShell({ eyebrow, title, description, actions, children }: DashboardShellProps) {
  return (
    <main className="min-h-screen text-foreground">
      <div className="mx-auto flex w-full max-w-[1500px] flex-col gap-5 px-4 py-4 md:px-6 md:py-6">
        <header className="rounded-lg border border-border bg-card/95 shadow-sm backdrop-blur">
          <div className="flex flex-col gap-4 p-4 lg:flex-row lg:items-center lg:justify-between">
            <div className="min-w-0">
              <Badge>{eyebrow}</Badge>
              <h1 className="mt-3 text-3xl font-semibold leading-tight text-foreground md:text-4xl">{title}</h1>
              <p className="mt-2 max-w-2xl text-sm text-muted-foreground">{description}</p>
            </div>
            {actions}
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
        {children}
      </div>
    </main>
  )
}
