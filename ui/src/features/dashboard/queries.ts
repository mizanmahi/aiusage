import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createUser, getDailySummary, getUserProjects, getUsers } from '@/lib/api'
import type { CreateUserInput } from '@/types'

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

export function useCreateUserMutation(options: Pick<QueryOptions, 'apiKey'>) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: CreateUserInput) => createUser(input, { apiKey: options.apiKey }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'users'] }),
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
