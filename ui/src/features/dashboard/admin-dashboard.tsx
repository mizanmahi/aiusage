import { useMemo, useState } from 'react'
import type { FormEvent } from 'react'
import type { LucideIcon } from 'lucide-react'
import { Database, FolderKanban, KeyRound, RefreshCw, Users } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { ThemeToggle } from '@/components/theme-toggle'
import { formatCost, formatDate, formatTokens } from '@/lib/format'
import { cn } from '@/lib/utils'
import type { ProjectSummary, UserSummary } from '@/types'
import { useDailySummaryQuery, useProjectsQuery, useUsersQuery } from './queries'

const storedAPIKey = 'aiusage.admin.apiKey'
const emptyUsers: UserSummary[] = []
const emptyProjects: ProjectSummary[] = []
const emptyDaily: Array<{ date: string; total_tokens: number }> = []

export function AdminDashboard() {
  const [apiKey, setAPIKey] = useState(() => sessionStorage.getItem(storedAPIKey) ?? '')
  const [activeKey, setActiveKey] = useState(apiKey)
  const [selectedUserID, setSelectedUserID] = useState('')
  const [from, setFrom] = useState('2026-01-01')
  const [to, setTo] = useState('2026-12-31')

  const hasKey = Boolean(activeKey)
  const usersQuery = useUsersQuery({ apiKey: activeKey, enabled: hasKey })
  const users = usersQuery.data ?? emptyUsers
  const selectedUser = users.find((user) => user.id === selectedUserID) ?? users[0]
  const projectsQuery = useProjectsQuery(selectedUser?.id ?? '', { apiKey: activeKey, enabled: hasKey })
  const dailyQuery = useDailySummaryQuery(from, to, { apiKey: activeKey, enabled: hasKey })

  const projects = projectsQuery.data ?? emptyProjects
  const daily = dailyQuery.data ?? emptyDaily
  const totals = useMemo(() => summarize(users), [users])
  const queryError = usersQuery.error ?? projectsQuery.error ?? dailyQuery.error
  const message = statusMessage(hasKey, queryError instanceof Error ? queryError : null)
  const isFetching = usersQuery.isFetching || projectsQuery.isFetching || dailyQuery.isFetching

  function submitKey(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    sessionStorage.setItem(storedAPIKey, apiKey)
    setActiveKey(apiKey)
  }

  return (
    <main className="min-h-screen bg-background text-foreground">
      <div className="mx-auto flex w-full max-w-[1440px] flex-col gap-4 p-4 md:p-6">
        <header className="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
          <div>
            <p className="text-xs font-bold uppercase text-muted-foreground">aiusage</p>
            <h1 className="mt-1 text-3xl font-semibold tracking-normal text-foreground">Admin Dashboard</h1>
          </div>
          <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
            <ThemeToggle />
            <form className="flex flex-col gap-2 sm:flex-row sm:items-end" onSubmit={submitKey}>
              <label className="grid gap-1 text-xs font-semibold text-muted-foreground">
                API key
                <Input
                  className="w-full sm:w-64"
                  type="password"
                  value={apiKey}
                  onChange={(event) => setAPIKey(event.target.value)}
                  autoComplete="current-password"
                />
              </label>
              <Button type="submit" disabled={!apiKey || isFetching}>
                <KeyRound className="size-4" />
                Refresh
              </Button>
            </form>
          </div>
        </header>

        <section className="grid gap-3 md:grid-cols-4">
          <Metric icon={Users} label="Users" value={users.length.toString()} />
          <Metric icon={Database} label="Tokens" value={formatTokens(totals.tokens)} />
          <Metric icon={RefreshCw} label="Cost" value={formatCost(totals.cost)} />
          <Metric icon={FolderKanban} label="Last active" value={formatDate(totals.lastActive)} />
        </section>

        <section className="grid gap-4 lg:grid-cols-[360px_minmax(0,1fr)]">
          <Card>
            <CardHeader>
              <CardTitle>Users</CardTitle>
              <Badge>{message}</Badge>
            </CardHeader>
            <CardContent className="space-y-2">
              {users.map((user) => (
                <UserRow
                  key={user.id}
                  user={user}
                  active={user.id === selectedUser?.id}
                  onSelect={() => setSelectedUserID(user.id)}
                />
              ))}
              {!users.length && <EmptyState label={hasKey ? 'No users found' : 'Enter an admin API key'} />}
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>{selectedUser ? selectedUser.name || selectedUser.email : 'Projects'}</CardTitle>
              <Badge>{projects.length} projects</Badge>
            </CardHeader>
            <CardContent className="space-y-1">
              <ProjectList projects={projects} />
            </CardContent>
          </Card>
        </section>

        <Card>
          <CardHeader className="flex-col items-stretch gap-3 py-3 md:flex-row md:items-center">
            <CardTitle>Daily Summary</CardTitle>
            <div className="flex flex-col gap-2 sm:flex-row">
              <Input type="date" value={from} onChange={(event) => setFrom(event.target.value)} />
              <Input type="date" value={to} onChange={(event) => setTo(event.target.value)} />
              <Button type="button" variant="outline" onClick={() => void dailyQuery.refetch()} disabled={!hasKey || isFetching}>
                Apply
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <DailyBars points={daily} />
          </CardContent>
        </Card>
      </div>
    </main>
  )
}

function Metric({ icon: Icon, label, value }: { icon: LucideIcon; label: string; value: string }) {
  return (
    <Card className="p-4">
      <div className="flex items-center justify-between gap-3">
        <span className="text-sm text-muted-foreground">{label}</span>
        <Icon className="size-4 text-muted-foreground" />
      </div>
      <strong className="mt-4 block text-2xl font-semibold text-foreground">{value}</strong>
    </Card>
  )
}

function UserRow({ user, active, onSelect }: { user: UserSummary; active: boolean; onSelect: () => void }) {
  return (
    <button
      type="button"
      className={cn(
        'grid h-14 w-full grid-cols-[minmax(0,1fr)_auto] items-center gap-3 rounded-md border px-3 text-left transition-colors',
        active ? 'border-primary bg-primary/10' : 'border-border hover:bg-muted',
      )}
      onClick={onSelect}
    >
      <span className="min-w-0">
        <strong className="block truncate text-sm text-foreground">{user.name || user.email}</strong>
        <small className="block truncate text-xs text-muted-foreground">{user.email}</small>
      </span>
      <strong className="text-xs text-foreground">{formatTokens(user.total_tokens)}</strong>
    </button>
  )
}

function ProjectList({ projects }: { projects: ProjectSummary[] }) {
  const maxTokens = Math.max(1, ...projects.map((project) => project.total_tokens))

  if (!projects.length) {
    return <EmptyState label="No project usage loaded" />
  }

  return projects.map((project) => (
    <div
      className="grid min-h-16 grid-cols-1 gap-3 border-b border-border py-3 last:border-b-0 md:grid-cols-[minmax(150px,1.3fr)_minmax(120px,1fr)_90px]"
      key={`${project.user_id}-${project.tool}-${project.project}`}
    >
      <div className="min-w-0">
        <strong className="block truncate text-sm text-foreground">{project.project}</strong>
        <span className="text-xs text-muted-foreground">
          {project.tool} · {formatDate(project.last_active)}
        </span>
      </div>
      <div className="h-2 self-center overflow-hidden rounded-full bg-muted">
        <div className="h-full rounded-full bg-primary" style={{ width: `${(project.total_tokens / maxTokens) * 100}%` }} />
      </div>
      <div className="text-left md:text-right">
        <strong className="block text-sm text-foreground">{formatTokens(project.total_tokens)}</strong>
        <span className="text-xs text-muted-foreground">{formatCost(project.total_cost_usd)}</span>
      </div>
    </div>
  ))
}

function DailyBars({ points }: { points: Array<{ date: string; total_tokens: number }> }) {
  const maxTokens = Math.max(1, ...points.map((point) => point.total_tokens))

  if (!points.length) {
    return <EmptyState label="No daily data loaded" />
  }

  return (
    <div className="space-y-2">
      {points.map((point) => (
        <div className="grid items-center gap-3 md:grid-cols-[104px_minmax(120px,1fr)_78px]" key={point.date}>
          <span className="text-xs text-muted-foreground">{point.date}</span>
          <div className="h-2 overflow-hidden rounded-full bg-muted">
            <div className="h-full rounded-full bg-primary" style={{ width: `${(point.total_tokens / maxTokens) * 100}%` }} />
          </div>
          <strong className="text-xs text-foreground md:text-right">{formatTokens(point.total_tokens)}</strong>
        </div>
      ))}
    </div>
  )
}

function EmptyState({ label }: { label: string }) {
  return <div className="grid min-h-24 place-items-center rounded-md border border-dashed border-border text-sm text-muted-foreground">{label}</div>
}

function summarize(users: UserSummary[]) {
  return users.reduce(
    (total, user) => ({
      tokens: total.tokens + user.total_tokens,
      cost: total.cost + user.total_cost_usd,
      lastActive: maxDate(total.lastActive, user.last_active),
    }),
    { tokens: 0, cost: 0, lastActive: '' },
  )
}

function statusMessage(hasKey: boolean, error: Error | null) {
  if (!hasKey) return 'Waiting'
  if (error) return error.message
  return 'Live'
}

function maxDate(left: string, right: string) {
  if (!left) return right
  if (!right) return left
  return left > right ? left : right
}
