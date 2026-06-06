package types

import "time"

type Tool string

const (
	ToolClaude Tool = "claude"
	ToolCodex  Tool = "codex"
)

type UsageEvent struct {
	SessionID       string    `json:"session_id"`
	UserID          string    `json:"user_id"`
	Date            string    `json:"date"` // YYYY-MM-DD
	Tool            Tool      `json:"tool"`
	Project         string    `json:"project"`
	Cwd             string    `json:"cwd"`
	Model           string    `json:"model"`
	InputTokens     int64     `json:"input_tokens"`
	OutputTokens    int64     `json:"output_tokens"`
	CacheTokens     int64     `json:"cache_tokens"`
	ReasoningTokens int64     `json:"reasoning_tokens"`
	CostUSD         float64   `json:"cost_usd"`
	PushedAt        time.Time `json:"pushed_at"`
}

type PushPayload struct {
	// API keys travel in the Authorization header so ingest auth is not mixed into event data.
	Events []UsageEvent `json:"events"`
}

type PushResponse struct {
	Accepted int    `json:"accepted"`
	Skipped  int    `json:"skipped"`
	Message  string `json:"message"`
}
