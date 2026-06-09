import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createUser, getDailySummary, getUserBreakdown, getUserProjects, getUserUsageSummary, getUsers } from '@/lib/api'
import type { BreakdownGroup, CreateUserInput, ProviderFilter } from '@/types'

type QueryOptions = {
  apiKey: string
  enabled: boolean
  authVersion?: number
}

export function useUsersQuery(options: QueryOptions) {
  return useQuery({
    queryKey: ['admin', 'users', options.authVersion],
    queryFn: () => getUsers({ apiKey: options.apiKey }),
    enabled: options.enabled,
  })
}

export function useCreateUserMutation(options: Pick<QueryOptions, 'apiKey'>) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: CreateUserInput) => createUser(input, { apiKey: options.apiKey }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'users'] }),
  })
}

export function useProjectsQuery(userID: string, options: QueryOptions) {
  return useQuery({
    queryKey: ['admin', 'users', userID, 'projects', options.authVersion],
    queryFn: () => getUserProjects(userID, { apiKey: options.apiKey }),
    enabled: options.enabled && Boolean(userID),
  })
}

export function useDailySummaryQuery(from: string, to: string, options: QueryOptions) {
  return useQuery({
    queryKey: ['admin', 'summary', from, to, options.authVersion],
    queryFn: () => getDailySummary(from, to, { apiKey: options.apiKey }),
    enabled: options.enabled,
  })
}

export function useUserBreakdownQuery(userID: string, filters: { groupBy: BreakdownGroup; provider: ProviderFilter; from: string; to: string }, options: QueryOptions) {
  return useQuery({
    queryKey: ['admin', 'users', userID, 'breakdown', filters.groupBy, filters.provider, filters.from, filters.to, options.authVersion],
    queryFn: () => getUserBreakdown(userID, filters, { apiKey: options.apiKey }),
    enabled: options.enabled && Boolean(userID),
  })
}

export function useUserUsageSummaryQuery(userID: string, filters: { provider: ProviderFilter; from: string; to: string }, options: QueryOptions) {
  return useQuery({
    queryKey: ['admin', 'users', userID, 'usage-summary', filters.provider, filters.from, filters.to, options.authVersion],
    queryFn: () => getUserUsageSummary(userID, filters, { apiKey: options.apiKey }),
    enabled: options.enabled && Boolean(userID),
  })
}
