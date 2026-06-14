import { ArrowUpDown } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { formatCost, formatTokens } from '@/lib/format'
import { cn } from '@/lib/utils'
import type { BreakdownGroup, UsageBreakdownRow } from '@/types'
import type { SortKey, SortState } from '../breakdown-utils'
import { capitalize } from '../lib/dashboard-format'

export function BreakdownTable({ groupBy, rows, sort, setSort, isLoading }: { groupBy: BreakdownGroup; rows: UsageBreakdownRow[]; sort: SortState; setSort: (value: SortState) => void; isLoading: boolean }) {
  const groupLabel = groupBy === 'month' ? 'Month' : groupBy === 'project' ? 'Project' : 'Date'

  return (
    <div className="overflow-x-auto rounded-md border border-border bg-background">
      <Table className="min-w-[980px]">
        <TableHeader className="bg-muted/70">
          <TableRow>
            {groupBy === 'project' ? <PlainHead label={groupLabel} /> : <SortableHead label={groupLabel} sortKey="group" sort={sort} setSort={setSort} />}
            <PlainHead label="Agent" />
            <PlainHead label="Models" />
            <PlainHead label="Input" align="right" />
            <PlainHead label="Output" align="right" />
            <PlainHead label="Cache Create" align="right" />
            <PlainHead label="Cache Read" align="right" />
            <SortableHead label="Total Tokens" sortKey="total_tokens" sort={sort} setSort={setSort} align="right" />
            <PlainHead label="Cost" align="right" />
          </TableRow>
        </TableHeader>
        <TableBody>
          {isLoading ? <BreakdownSkeleton /> : rows.map((row, index) => <BreakdownRow key={`${row.group}-${row.agent}-${index}`} row={row} />)}
          {!isLoading && !rows.length && (
            <TableRow>
              <TableCell className="py-10 text-center text-muted-foreground" colSpan={9}>
                No usage loaded
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  )
}

function BreakdownRow({ row }: { row: UsageBreakdownRow }) {
  return (
    <TableRow>
      <TableCell className="font-medium text-foreground">{row.agent === 'all' ? row.group : ''}</TableCell>
      <TableCell className="text-foreground">{row.agent === 'all' ? 'All' : capitalize(row.agent)}</TableCell>
      <TableCell className="max-w-56 truncate text-muted-foreground">{row.models.join(', ')}</TableCell>
      <NumberCell value={row.input_tokens} />
      <NumberCell value={row.output_tokens} />
      <NumberCell value={row.cache_creation_tokens} />
      <NumberCell value={row.cache_read_tokens} />
      <NumberCell value={row.total_tokens} />
      <TableCell className="text-right font-medium text-foreground">{formatCost(row.total_cost_usd)}</TableCell>
    </TableRow>
  )
}

function SortableHead({ label, sortKey, sort, setSort, align = 'left' }: { label: string; sortKey: SortKey; sort: SortState; setSort: (value: SortState) => void; align?: 'left' | 'right' }) {
  const active = sort.key === sortKey
  const nextDirection = active && sort.direction === 'asc' ? 'desc' : 'asc'
  return (
    <TableHead className={cn(align === 'right' && 'text-right')}>
      <Button type="button" variant="ghost" className={cn('h-7 px-2', align === 'right' && 'ml-auto')} onClick={() => setSort({ key: sortKey, direction: nextDirection })}>
        {label}
        <ArrowUpDown data-icon="inline-end" className={cn(active && 'text-foreground')} />
      </Button>
    </TableHead>
  )
}

function PlainHead({ label, align = 'left' }: { label: string; align?: 'left' | 'right' }) {
  return <TableHead className={cn(align === 'right' && 'text-right')}>{label}</TableHead>
}

function NumberCell({ value }: { value: number }) {
  return <TableCell className="text-right tabular-nums text-foreground">{formatTokens(value)}</TableCell>
}

function BreakdownSkeleton() {
  return (
    <>
      {[0, 1, 2, 3, 4].map((item) => (
        <TableRow key={item}>
          {Array.from({ length: 9 }).map((_, index) => (
            <TableCell key={index}>
              <Skeleton className="h-4 w-full" />
            </TableCell>
          ))}
        </TableRow>
      ))}
    </>
  )
}
