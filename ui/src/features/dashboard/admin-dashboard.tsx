import { useMemo, useState } from 'react'
import { Activity, Database, RefreshCw, Users } from 'lucide-react'
import { formatCost, formatDate, formatTokens } from '@/lib/format'
import type { UserSummary } from '@/types'
import { CreateUserPanel } from './create-user-panel'
import { AdminTokenDialog } from './components/admin-token-dialog'
import { DashboardShell } from './components/dashboard-shell'
import { DeveloperSelect } from './components/developer-select'
import { MetricGrid } from './components/metric-grid'
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

  function loadAdminToken(nextKey: string) {
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
      isPowered={hasKey}
      navActions={<AdminTokenDialog token={apiKey} isLoading={usersQuery.isFetching} isPowered={hasKey} onLoad={loadAdminToken} />}
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

      <section className="grid gap-4 lg:grid-cols-2">
        <CreateUserPanel
          enabled={hasKey}
          isCreating={createUser.isPending}
          error={createUser.error instanceof Error ? createUser.error : null}
          onCreate={(input) => createUser.mutateAsync(input)}
        />
        <DeveloperSelect
          users={users}
          selectedUserID={selectedUser?.id}
          hasKey={hasKey}
          isLoading={usersQuery.isLoading}
          onSelect={setSelectedUserID}
        />
      </section>
      <section>
        <UserAnalyticsPanel user={selectedUser} apiKey={activeKey} enabled={hasKey} authVersion={authVersion} />
      </section>
    </DashboardShell>
  )
}

function statusMessage(hasKey: boolean, error: unknown) {
  if (!hasKey) return 'Waiting for key'
  if (error instanceof Error) return error.message
  return 'Live'
}
