import type {
  APIErrorResponse,
  APIResponse,
  BreakdownGroup,
  CreateUserInput,
  CreateUserResult,
  DailyPoint,
  ProjectSummary,
  ProviderFilter,
  UsageBreakdownRow,
  UsageSummaryStats,
  UserSummary,
} from '@/types'

const BASE_URL = import.meta.env.VITE_API_BASE_URL ?? ''

export type APIOptions = {
  apiKey: string
}

export async function getUsers(options: APIOptions): Promise<UserSummary[]> {
  return request<UserSummary[]>('/admin/users', options)
}

export async function createUser(input: CreateUserInput, options: APIOptions): Promise<CreateUserResult> {
  return request<CreateUserResult>('/admin/users', options, {
    method: 'POST',
    body: JSON.stringify(input),
  })
}

export async function getUserProjects(userID: string, options: APIOptions): Promise<ProjectSummary[]> {
  return request<ProjectSummary[]>(`/admin/users/${encodeURIComponent(userID)}`, options)
}

export async function getDailySummary(from: string, to: string, options: APIOptions): Promise<DailyPoint[]> {
  const params = new URLSearchParams()
  if (from) params.set('from', from)
  if (to) params.set('to', to)

  const query = params.toString()
  return request<DailyPoint[]>(`/admin/summary${query ? `?${query}` : ''}`, options)
}

export async function getUserBreakdown(
  userID: string,
  filters: { groupBy: BreakdownGroup; from: string; to: string },
  options: APIOptions,
): Promise<UsageBreakdownRow[]> {
  const params = new URLSearchParams({ group_by: filters.groupBy })
  if (filters.from) params.set('from', filters.from)
  if (filters.to) params.set('to', filters.to)

  return request<UsageBreakdownRow[]>(`/admin/users/${encodeURIComponent(userID)}/breakdown?${params.toString()}`, options)
}

export async function getUserUsageSummary(
  userID: string,
  filters: { provider: ProviderFilter; from: string; to: string },
  options: APIOptions,
): Promise<UsageSummaryStats> {
  const params = new URLSearchParams({ provider: filters.provider })
  if (filters.from) params.set('from', filters.from)
  if (filters.to) params.set('to', filters.to)

  return request<UsageSummaryStats>(`/admin/users/${encodeURIComponent(userID)}/summary?${params.toString()}`, options)
}

async function request<T>(path: string, options: APIOptions, init: RequestInit = {}): Promise<T> {
  const response = await fetch(`${BASE_URL}${path}`, {
    ...init,
    headers: {
      ...init.headers,
      Authorization: `Bearer ${options.apiKey}`,
      ...(init.body ? { 'Content-Type': 'application/json' } : {}),
    },
  })

  const payload = (await response.json()) as APIResponse<T> | APIErrorResponse
  if (!response.ok || 'error' in payload) {
    const message = 'error' in payload ? payload.error.message : `request failed with ${response.status}`
    throw new Error(message)
  }

  return payload.data
}
