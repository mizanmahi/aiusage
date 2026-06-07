package claude

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mizanmahi/aiusage/types"
)

type assistantLine struct {
	Type      string `json:"type"`
	SessionID string `json:"sessionId"`
	Timestamp string `json:"timestamp"`
	Message   struct {
		Model string `json:"model"`
		Usage struct {
			InputTokens              int64 `json:"input_tokens"`
			OutputTokens             int64 `json:"output_tokens"`
			CacheCreationInputTokens int64 `json:"cache_creation_input_tokens"`
			CacheReadInputTokens     int64 `json:"cache_read_input_tokens"`
		} `json:"usage"`
	} `json:"message"`
}

type Session struct {
	ID              string
	Project         string
	Cwd             string
	Date            string
	Model           string
	InputTokens     int64
	OutputTokens    int64
	CacheTokens     int64
	ReasoningTokens int64
}

func ReadSessions(claudeHome string, since time.Time) ([]Session, error) {
	projectsDir := filepath.Join(claudeHome, "projects")
	projectDirs, err := os.ReadDir(projectsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var sessions []Session
	for _, projectDir := range projectDirs {
		if !projectDir.IsDir() {
			continue
		}

		project := decodePath(projectDir.Name())
		cwd := encodedToPath(projectDir.Name())
		sessionDir := filepath.Join(projectsDir, projectDir.Name())

		files, err := os.ReadDir(sessionDir)
		if err != nil {
			continue
		}

		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".jsonl") {
				continue
			}

			info, err := file.Info()
			if err != nil {
				continue
			}
			if info.ModTime().Before(since) {
				continue
			}

			path := filepath.Join(sessionDir, file.Name())
			session, err := parseSessionFile(path, project, cwd)
			if err != nil || session == nil {
				continue
			}

			sessions = append(sessions, *session)
		}
	}

	return sessions, nil
}

func parseSessionFile(path, project, cwd string) (*Session, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	session := Session{
		Project: project,
		Cwd:     cwd,
	}

	scanner := bufio.NewScanner(file)
	// Claude messages can contain large tool outputs or prompt context.
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
		if peek.Type != "assistant" {
			continue
		}

		var assistant assistantLine
		if err := json.Unmarshal(line, &assistant); err != nil {
			continue
		}
		applyAssistantLine(assistant, &session)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if session.ID == "" {
		return nil, nil
	}

	return &session, nil
}

func applyAssistantLine(line assistantLine, session *Session) {
	if session.ID == "" {
		session.ID = line.SessionID
	}
	if session.Date == "" && len(line.Timestamp) >= len("2006-01-02") {
		session.Date = line.Timestamp[:len("2006-01-02")]
	}
	if session.Model == "" {
		session.Model = line.Message.Model
	}

	session.InputTokens += line.Message.Usage.InputTokens
	session.OutputTokens += line.Message.Usage.OutputTokens
	// Claude exposes cache writes and reads separately; aiusage stores them together.
	session.CacheTokens += line.Message.Usage.CacheCreationInputTokens + line.Message.Usage.CacheReadInputTokens
}

func decodePath(encoded string) string {
	trimmed := strings.Trim(encoded, "-")
	if trimmed == "" {
		return encoded
	}

	parts := strings.Split(trimmed, "-")
	return parts[len(parts)-1]
}

func encodedToPath(encoded string) string {
	trimmed := strings.TrimPrefix(encoded, "-")
	if trimmed == "" {
		return "/"
	}

	return "/" + strings.ReplaceAll(trimmed, "-", "/")
}

func (s Session) ToUsageEvent(userID string) types.UsageEvent {
	return types.UsageEvent{
		SessionID:       s.ID,
		UserID:          userID,
		Date:            s.Date,
		Tool:            types.ToolClaude,
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
