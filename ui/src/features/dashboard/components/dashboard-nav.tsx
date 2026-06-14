import { type ReactNode } from 'react'
import { ThemeToggle } from '@/components/theme-toggle'

export function DashboardNav({ actions }: { actions: ReactNode }) {
  return (
    <nav className="flex h-14 items-center justify-between rounded-lg border border-border bg-card/95 px-4 shadow-sm backdrop-blur">
      <a className="cursor-pointer text-sm font-semibold text-foreground" href="/">
        aiusage
      </a>
      <div className="flex items-center gap-2">
        {actions}
        <ThemeToggle />
      </div>
    </nav>
  )
}
