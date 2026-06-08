import type { APIErrorResponse, APIResponse, CreateUserInput, CreateUserResult, DailyPoint, ProjectSummary, UserSummary } from '@/types'

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
