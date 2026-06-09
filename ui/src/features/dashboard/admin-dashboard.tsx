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
import type { UserSummary } from '@/types'
import { CreateUserPanel } from './create-user-panel'
import { useCreateUserMutation, useUsersQuery } from './queries'
import { UserAnalyticsPanel } from './user-analytics-panel'

const storedAPIKey = 'aiusage.admin.apiKey'
const emptyUsers: UserSummary[] = []

export function AdminDashboard() {
  const [apiKey, setAPIKey] = useState(() => sessionStorage.getItem(storedAPIKey) ?? '')
  const [activeKey, setActiveKey] = useState(apiKey)
  const [authVersion, setAuthVersion] = useState(0)
  const [selectedUserID, setSelectedUserID] = useState('')

  const hasKey = Boolean(activeKey)
  const usersQuery = useUsersQuery({ apiKey: activeKey, enabled: hasKey, authVersion })
  const users = usersQuery.data ?? emptyUsers
  const selectedUser = users.find((user) => user.id === selectedUserID) ?? users[0]
  const createUser = useCreateUserMutation({ apiKey: activeKey })

  const totals = useMemo(() => summarize(users), [users])
  const message = statusMessage(hasKey, usersQuery.error instanceof Error ? usersQuery.error : null)
  const isFetching = usersQuery.isFetching

  function submitKey(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const nextKey = apiKey.trim()
    sessionStorage.setItem(storedAPIKey, nextKey)
    setAPIKey(nextKey)
    setActiveKey(nextKey)
    setAuthVersion((version) => version + 1)
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
              <Button type="submit" disabled={!apiKey.trim() || isFetching}>
                <KeyRound className="size-4" />
                Load
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
          <div className="space-y-4">
            <CreateUserPanel
              enabled={hasKey}
              isCreating={createUser.isPending}
              error={createUser.error instanceof Error ? createUser.error : null}
              onCreate={(input) => createUser.mutateAsync(input)}
            />
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
          </div>

          <UserAnalyticsPanel user={selectedUser} apiKey={activeKey} enabled={hasKey} authVersion={authVersion} />
        </section>
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
