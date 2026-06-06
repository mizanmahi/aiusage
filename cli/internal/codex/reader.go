package codex

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mizanmahi/aiusage/types"
)

type sessionMeta struct {
	Payload struct {
		ID  string `json:"id"`
		Cwd string `json:"cwd"`
	} `json:"payload"`
}

type eventMsg struct {
	Payload struct {
		Type    string `json:"type"`
		TurnCtx struct {
			Model string `json:"model"`
		} `json:"turn_context"`
		InputTokens     int64 `json:"input_tokens"`
		OutputTokens    int64 `json:"output_tokens"`
		CachedTokens    int64 `json:"cached_tokens"`
		ReasoningTokens int64 `json:"reasoning_tokens"`
		Info            struct {
			TotalTokenUsage tokenUsage `json:"total_token_usage"`
			LastTokenUsage  tokenUsage `json:"last_token_usage"`
		} `json:"info"`
	} `json:"payload"`
}

type tokenUsage struct {
	InputTokens           int64 `json:"input_tokens"`
	CachedInputTokens     int64 `json:"cached_input_tokens"`
	OutputTokens          int64 `json:"output_tokens"`
	ReasoningOutputTokens int64 `json:"reasoning_output_tokens"`
}

type Session struct {
	ID              string
	Cwd             string
	Project         string
	Date            string
	Model           string
	InputTokens     int64
	OutputTokens    int64
	CacheTokens     int64
	ReasoningTokens int64
}

func ReadSessions(codexHome string, since time.Time) ([]Session, error) {
	sessionsDir := filepath.Join(codexHome, "sessions")
	var sessions []Session

	err := filepath.WalkDir(sessionsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".jsonl") {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}
		if info.ModTime().Before(since) {
			return nil
		}

		session, err := parseSessionFile(path)
		if err != nil || session == nil {
			return nil
		}

		sessions = append(sessions, *session)
		return nil
	})
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	return sessions, nil
}

func parseSessionFile(path string) (*Session, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var session Session
	var metaSeen bool

	scanner := bufio.NewScanner(file)
	// Codex JSONL lines can be large because some events include prompt/context data.
	scanner.Buffer(make([]byte, 1024*1024), 4*1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var peek struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(line, &peek); err != nil {
			continue
		}

		switch peek.Type {
		case "session_meta":
			if metaSeen {
				continue
			}
			if err := applySessionMeta(line, path, &session); err != nil {
				continue
			}
			metaSeen = true
		case "event_msg":
			applyTokenCount(line, &session)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if !metaSeen || session.ID == "" {
		return nil, nil
	}

	return &session, nil
}

func applySessionMeta(line []byte, path string, session *Session) error {
	var meta sessionMeta
	if err := json.Unmarshal(line, &meta); err != nil {
		return err
	}

	session.ID = meta.Payload.ID
	session.Cwd = meta.Payload.Cwd
	session.Project = filepath.Base(meta.Payload.Cwd)
	session.Date = dateFromSessionPath(path)
	return nil
}

func applyTokenCount(line []byte, session *Session) {
	var event eventMsg
	if err := json.Unmarshal(line, &event); err != nil {
		return
	}
	if event.Payload.Type != "token_count" {
		return
	}

	switch {
	case hasTokenUsage(event.Payload.Info.LastTokenUsage):
		// Current Codex writes per-turn deltas under last_token_usage.
		addTokenUsage(session, event.Payload.Info.LastTokenUsage)
	case hasTokenUsage(event.Payload.Info.TotalTokenUsage):
		// total_token_usage is cumulative, so replace with the latest value.
		setTokenUsage(session, event.Payload.Info.TotalTokenUsage)
	default:
		// Older/assumed Codex shapes put token fields directly under payload.
		session.InputTokens += event.Payload.InputTokens
		session.OutputTokens += event.Payload.OutputTokens
		session.CacheTokens += event.Payload.CachedTokens
		session.ReasoningTokens += event.Payload.ReasoningTokens
	}

	if session.Model == "" {
		session.Model = event.Payload.TurnCtx.Model
	}
}

func hasTokenUsage(usage tokenUsage) bool {
	return usage.InputTokens != 0 ||
		usage.CachedInputTokens != 0 ||
		usage.OutputTokens != 0 ||
		usage.ReasoningOutputTokens != 0
}

func addTokenUsage(session *Session, usage tokenUsage) {
	session.InputTokens += usage.InputTokens
	session.OutputTokens += usage.OutputTokens
	session.CacheTokens += usage.CachedInputTokens
	session.ReasoningTokens += usage.ReasoningOutputTokens
}

func setTokenUsage(session *Session, usage tokenUsage) {
	session.InputTokens = usage.InputTokens
	session.OutputTokens = usage.OutputTokens
	session.CacheTokens = usage.CachedInputTokens
	session.ReasoningTokens = usage.ReasoningOutputTokens
}

func dateFromSessionPath(path string) string {
	parts := strings.Split(filepath.ToSlash(path), "/")
	if len(parts) < 4 {
		return ""
	}

	n := len(parts)
	return parts[n-4] + "-" + parts[n-3] + "-" + parts[n-2]
}

func (s Session) ToUsageEvent(userID string) types.UsageEvent {
	return types.UsageEvent{
		SessionID:       s.ID,
		UserID:          userID,
		Date:            s.Date,
		Tool:            types.ToolCodex,
		Project:         s.Project,
		Cwd:             s.Cwd,
		Model:           s.Model,
		InputTokens:     s.InputTokens,
		OutputTokens:    s.OutputTokens,
		CacheTokens:     s.CacheTokens,
		ReasoningTokens: s.ReasoningTokens,
		PushedAt:        time.Now(),
	}
}
