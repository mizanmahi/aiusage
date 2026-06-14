import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { formatCost, formatDate, formatTokens } from '@/lib/format'
import type { UserSummary } from '@/types'
import { EmptyState } from './empty-state'

export function DeveloperSelect({
  users,
  selectedUserID,
  hasKey,
  isLoading,
  onSelect,
}: {
  users: UserSummary[]
  selectedUserID?: string
  hasKey: boolean
  isLoading: boolean
  onSelect: (userID: string) => void
}) {
  const selectedUser = users.find((user) => user.id === selectedUserID) ?? users[0]

  return (
    <Card>
      <CardHeader>
        <div>
          <CardTitle>Developers</CardTitle>
          <p className="text-xs text-muted-foreground">Select a teammate</p>
        </div>
        <Badge>{isLoading ? 'Syncing' : `${users.length} users`}</Badge>
      </CardHeader>
      <CardContent className="flex flex-col gap-3">
        {isLoading && <DeveloperSelectSkeleton />}
        {!isLoading && users.length > 0 && (
          <>
            <select
              className="h-9 w-full cursor-pointer rounded-md border border-input bg-background px-3 text-sm text-foreground shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
              value={selectedUser?.id ?? ''}
              onChange={(event) => onSelect(event.target.value)}
            >
              {users.map((user) => (
                <option key={user.id} value={user.id}>
                  {user.name || user.email}
                </option>
              ))}
            </select>
            {selectedUser && (
              <div className="grid gap-2 rounded-md border border-border bg-background p-3 text-sm">
                <p className="truncate font-medium text-foreground">{selectedUser.email}</p>
                <div className="grid grid-cols-3 gap-2 text-xs text-muted-foreground">
                  <span>{formatTokens(selectedUser.total_tokens)} tokens</span>
                  <span>{formatCost(selectedUser.total_cost_usd)}</span>
                  <span>{formatDate(selectedUser.last_active)}</span>
                </div>
              </div>
            )}
          </>
        )}
        {!isLoading && !users.length && (
          <EmptyState
            title={hasKey ? 'No developers found' : 'Admin token required'}
            description={hasKey ? 'Create a developer to start collecting usage.' : 'Load an admin token to fetch developers.'}
          />
        )}
      </CardContent>
    </Card>
  )
}

function DeveloperSelectSkeleton() {
  return (
    <>
      <Skeleton className="h-9 w-full" />
      <Skeleton className="h-20 w-full" />
    </>
  )
}
