import { useMemo, useState } from 'react'
import type { FormEvent } from 'react'
import { Activity, Database, FolderKanban, RefreshCw, Users } from 'lucide-react'
import { ThemeToggle } from '@/components/theme-toggle'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { formatCost, formatDate, formatTokens } from '@/lib/format'
import type { UserSummary } from '@/types'
import { CreateUserPanel } from './create-user-panel'
import { DashboardShell } from './components/dashboard-shell'
import { MetricGrid } from './components/metric-grid'
import { UserList } from './components/user-list'
import { summarizeUsers } from './lib/dashboard-format'
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
  const totals = useMemo(() => summarizeUsers(users), [users])

  function submitKey(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const nextKey = apiKey.trim()
    sessionStorage.setItem(storedAPIKey, nextKey)
    setAPIKey(nextKey)
    setActiveKey(nextKey)
    setAuthVersion((version) => version + 1)
  }

  return (
    <DashboardShell
      eyebrow="Internal AI telemetry"
      title="aiusage"
      description="Monitor Claude Code and Codex usage by developer, project, model, and spend."
      actions={<HeaderActions apiKey={apiKey} isFetching={usersQuery.isFetching} onAPIKeyChange={setAPIKey} onSubmit={submitKey} />}
    >
      <MetricGrid
        isLoading={usersQuery.isLoading}
        metrics={[
          { icon: Users, label: 'Developers', value: users.length.toString(), detail: statusMessage(hasKey, usersQuery.error) },
          { icon: Database, label: 'Total tokens', value: formatTokens(totals.tokens), detail: 'Input, output, cache, reasoning' },
          { icon: Activity, label: 'Estimated cost', value: formatCost(totals.cost), detail: 'Server-side model pricing' },
          { icon: RefreshCw, label: 'Last active', value: formatDate(totals.lastActive), detail: 'Latest tracked session' },
        ]}
      />

      <section className="grid gap-4 lg:grid-cols-[380px_minmax(0,1fr)]">
        <div className="flex flex-col gap-4">
          <CreateUserPanel
            enabled={hasKey}
            isCreating={createUser.isPending}
            error={createUser.error instanceof Error ? createUser.error : null}
            onCreate={(input) => createUser.mutateAsync(input)}
          />
          <UserList
            users={users}
            selectedUserID={selectedUser?.id}
            hasKey={hasKey}
            isLoading={usersQuery.isLoading}
            onSelect={setSelectedUserID}
          />
        </div>

        <UserAnalyticsPanel user={selectedUser} apiKey={activeKey} enabled={hasKey} authVersion={authVersion} />
      </section>
    </DashboardShell>
  )
}

function HeaderActions({
  apiKey,
  isFetching,
  onAPIKeyChange,
  onSubmit,
}: {
  apiKey: string
  isFetching: boolean
  onAPIKeyChange: (value: string) => void
  onSubmit: (event: FormEvent<HTMLFormElement>) => void
}) {
  return (
    <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
      <ThemeToggle />
      <form className="flex flex-col gap-2 sm:flex-row sm:items-center" onSubmit={onSubmit}>
        <Input
          className="sm:w-72"
          type="password"
          value={apiKey}
          onChange={(event) => onAPIKeyChange(event.target.value)}
          placeholder="Admin API key"
          autoComplete="current-password"
        />
        <Button type="submit" disabled={!apiKey.trim() || isFetching}>
          <FolderKanban data-icon="inline-start" />
          Load
        </Button>
      </form>
    </div>
  )
}

function statusMessage(hasKey: boolean, error: unknown) {
  if (!hasKey) return 'Waiting for key'
  if (error instanceof Error) return error.message
  return 'Live'
}
