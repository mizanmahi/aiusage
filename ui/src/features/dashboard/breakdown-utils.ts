import type { ProviderFilter, UsageBreakdownRow } from '@/types'

export type SortKey = 'group' | 'total_tokens'
export type SortState = { key: SortKey; direction: 'asc' | 'desc' }

export function rowsForProvider(rows: UsageBreakdownRow[], provider: ProviderFilter) {
  if (provider === 'all') return rows

  const details = rows.filter((row) => row.agent === provider)
  const groups = new Map<string, UsageBreakdownRow[]>()
  for (const row of details) {
    groups.set(row.group, [...(groups.get(row.group) ?? []), row])
  }

  return [...groups.entries()].flatMap(([group, groupRows]) => [rollupRow(group, groupRows), ...groupRows])
}

export function sortRows(rows: UsageBreakdownRow[], sort: SortState) {
  return [...rows].sort((left, right) => {
    const direction = sort.direction === 'asc' ? 1 : -1
    if (sort.key === 'group') {
      const groupCompare = compareValue(left.group, right.group)
      if (groupCompare !== 0) return groupCompare * direction
      return agentRank(left.agent) - agentRank(right.agent)
    }
    return compareValue(left.total_tokens, right.total_tokens) * direction
  })
}

function rollupRow(group: string, rows: UsageBreakdownRow[]): UsageBreakdownRow {
  return {
    group,
    agent: 'all',
    models: [...new Set(rows.flatMap((row) => row.models))].sort(),
    input_tokens: sum(rows, 'input_tokens'),
    output_tokens: sum(rows, 'output_tokens'),
    cache_creation_tokens: sum(rows, 'cache_creation_tokens'),
    cache_read_tokens: sum(rows, 'cache_read_tokens'),
    reasoning_tokens: sum(rows, 'reasoning_tokens'),
    total_tokens: sum(rows, 'total_tokens'),
    total_cost_usd: rows.reduce((total, row) => total + row.total_cost_usd, 0),
    last_active: rows.reduce((latest, row) => (row.last_active > latest ? row.last_active : latest), ''),
  }
}

function sum(rows: UsageBreakdownRow[], key: 'input_tokens' | 'output_tokens' | 'cache_creation_tokens' | 'cache_read_tokens' | 'reasoning_tokens' | 'total_tokens') {
  return rows.reduce((total, row) => total + row[key], 0)
}

function compareValue(left: string | number, right: string | number) {
  if (typeof left === 'number' && typeof right === 'number') return left - right
  return String(left).localeCompare(String(right))
}

function agentRank(agent: string) {
  if (agent === 'all') return 0
  if (agent === 'codex') return 1
  if (agent === 'claude') return 2
  return 3
}
