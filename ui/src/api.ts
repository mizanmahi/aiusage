import type { APIErrorResponse, APIResponse, DailyPoint, ProjectSummary, UserSummary } from './types'

const BASE_URL = import.meta.env.VITE_API_BASE_URL ?? ''

type RequestOptions = {
  apiKey: string
}

export async function getUsers(options: RequestOptions): Promise<UserSummary[]> {
  return request<UserSummary[]>('/admin/users', options)
}

export async function getUserProjects(userID: string, options: RequestOptions): Promise<ProjectSummary[]> {
  return request<ProjectSummary[]>(`/admin/users/${encodeURIComponent(userID)}`, options)
}

export async function getDailySummary(
  from: string,
  to: string,
  options: RequestOptions,
): Promise<DailyPoint[]> {
  const params = new URLSearchParams()
  if (from) params.set('from', from)
  if (to) params.set('to', to)

  const query = params.toString()
  return request<DailyPoint[]>(`/admin/summary${query ? `?${query}` : ''}`, options)
}

async function request<T>(path: string, options: RequestOptions): Promise<T> {
  const response = await fetch(`${BASE_URL}${path}`, {
    headers: {
      Authorization: `Bearer ${options.apiKey}`,
    },
  })

  const payload = (await response.json()) as APIResponse<T> | APIErrorResponse
  if (!response.ok || 'error' in payload) {
    const message = 'error' in payload ? payload.error.message : `request failed with ${response.status}`
    throw new Error(message)
  }

  return payload.data
}
