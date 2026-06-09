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
  is_admin: boolean
  total_tokens: number
  total_cost_usd: number
  last_active: string
}

export type CreateUserInput = {
  email: string
  name: string
  is_admin: boolean
}

export type CreateUserResult = {
  user: UserSummary
  api_key: string
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

export type BreakdownGroup = 'day' | 'month' | 'project'
export type ProviderFilter = 'all' | 'codex' | 'claude'

export type UsageBreakdownRow = {
  group: string
  agent: string
  models: string[]
  input_tokens: number
  output_tokens: number
  cache_creation_tokens: number
  cache_read_tokens: number
  reasoning_tokens: number
  total_tokens: number
  total_cost_usd: number
  last_active: string
}

export type UsageSummaryStats = {
  provider: ProviderFilter
  total_projects: number
  total_input_tokens: number
  total_output_tokens: number
  total_cached_tokens: number
  total_tokens: number
  total_cost_usd: number
}
