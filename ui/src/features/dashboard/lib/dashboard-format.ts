import type { UserSummary } from '@/types'

export function summarizeUsers(users: UserSummary[]) {
  return users.reduce(
    (total, user) => ({
      tokens: total.tokens + user.total_tokens,
      cost: total.cost + user.total_cost_usd,
      lastActive: maxDate(total.lastActive, user.last_active),
    }),
    { tokens: 0, cost: 0, lastActive: '' },
  )
}

export function capitalize(value: string) {
  return value.charAt(0).toUpperCase() + value.slice(1)
}

function maxDate(left: string, right: string) {
  if (!left) return right
  if (!right) return left
  return left > right ? left : right
}
