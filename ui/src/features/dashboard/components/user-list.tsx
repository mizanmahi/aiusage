import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { formatCost, formatDate, formatTokens } from '@/lib/format'
import { cn } from '@/lib/utils'
import type { UserSummary } from '@/types'
import { EmptyState } from './empty-state'

export function UserList({
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
  return (
    <Card>
      <CardHeader>
        <div>
          <CardTitle>Developers</CardTitle>
          <p className="text-xs text-muted-foreground">Team usage leaderboard</p>
        </div>
        <Badge>{isLoading ? 'Syncing' : `${users.length} users`}</Badge>
      </CardHeader>
      <CardContent className="flex flex-col gap-2">
        {isLoading && <UserListSkeleton />}
        {!isLoading &&
          users.map((user) => (
            <UserRow key={user.id} active={user.id === selectedUserID} user={user} onSelect={() => onSelect(user.id)} />
          ))}
        {!isLoading && !users.length && (
          <EmptyState
            title={hasKey ? 'No users found' : 'Admin key required'}
            description={hasKey ? 'Create a developer to start collecting usage.' : 'Enter an admin API key to load dashboard data.'}
          />
        )}
      </CardContent>
    </Card>
  )
}

function UserRow({ user, active, onSelect }: { user: UserSummary; active: boolean; onSelect: () => void }) {
  return (
    <button
      type="button"
      className={cn(
        'grid min-h-20 w-full grid-cols-[minmax(0,1fr)_auto] items-center gap-3 rounded-md border p-3 text-left transition-colors',
        active ? 'border-primary bg-primary/10 shadow-sm' : 'border-border bg-background hover:bg-muted/70',
      )}
      onClick={onSelect}
    >
      <span className="min-w-0">
        <strong className="block truncate text-sm text-foreground">{user.name || user.email}</strong>
        <small className="block truncate text-xs text-muted-foreground">{user.email}</small>
        <small className="mt-2 block text-xs text-muted-foreground">Last active {formatDate(user.last_active)}</small>
      </span>
      <span className="text-right">
        <strong className="block text-sm text-foreground">{formatTokens(user.total_tokens)}</strong>
        <small className="block text-xs text-muted-foreground">{formatCost(user.total_cost_usd)}</small>
      </span>
    </button>
  )
}

function UserListSkeleton() {
  return (
    <>
      {[0, 1, 2, 3].map((item) => (
        <div className="rounded-md border border-border bg-background p-3" key={item}>
          <Skeleton className="h-4 w-40" />
          <Skeleton className="mt-2 h-3 w-56" />
          <Skeleton className="mt-3 h-3 w-28" />
        </div>
      ))}
    </>
  )
}
