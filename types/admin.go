package types

type UserSummary struct {
	ID          string  `json:"id"`
	Email       string  `json:"email"`
	Name        string  `json:"name"`
	TotalTokens int64   `json:"total_tokens"`
	TotalCost   float64 `json:"total_cost_usd"`
	LastActive  string  `json:"last_active"`
}

type ProjectSummary struct {
	Project     string  `json:"project"`
	UserID      string  `json:"user_id"`
	Tool        string  `json:"tool"`
	TotalTokens int64   `json:"total_tokens"`
	TotalCost   float64 `json:"total_cost_usd"`
	LastActive  string  `json:"last_active"`
}

type DailyPoint struct {
	Date        string  `json:"date"`
	TotalTokens int64   `json:"total_tokens"`
	TotalCost   float64 `json:"total_cost_usd"`
}

type APIResponse[T any] struct {
	Data T `json:"data"`
}

type APIErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
