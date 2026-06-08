import { useQuery } from '@tanstack/react-query'
import { getDailySummary, getUserProjects, getUsers } from '@/lib/api'

type QueryOptions = {
  apiKey: string
  enabled: boolean
}

export function useUsersQuery(options: QueryOptions) {
  return useQuery({
    queryKey: ['admin', 'users'],
    queryFn: () => getUsers({ apiKey: options.apiKey }),
    enabled: options.enabled,
  })
}

export function useProjectsQuery(userID: string, options: QueryOptions) {
  return useQuery({
    queryKey: ['admin', 'users', userID, 'projects'],
    queryFn: () => getUserProjects(userID, { apiKey: options.apiKey }),
    enabled: options.enabled && Boolean(userID),
  })
}

export function useDailySummaryQuery(from: string, to: string, options: QueryOptions) {
  return useQuery({
    queryKey: ['admin', 'summary', from, to],
    queryFn: () => getDailySummary(from, to, { apiKey: options.apiKey }),
    enabled: options.enabled,
  })
}
