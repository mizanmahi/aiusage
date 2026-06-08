export type APIResponse<T> = {
  data: T
}

export type APIErrorResponse = {
  error: {
    code: string
    message: string
  }
}

export type UserSummary = {
  id: string
  email: string
  name: string
  total_tokens: number
  total_cost_usd: number
  last_active: string
}

export type ProjectSummary = {
  project: string
  user_id: string
  tool: string
  total_tokens: number
  total_cost_usd: number
  last_active: string
}

export type DailyPoint = {
  date: string
  total_tokens: number
  total_cost_usd: number
}
