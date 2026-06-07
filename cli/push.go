package main

import (
	"fmt"
	"io"
	"time"

	"github.com/mizanmahi/aiusage/cli/internal/claude"
	"github.com/mizanmahi/aiusage/cli/internal/codex"
	"github.com/mizanmahi/aiusage/cli/internal/config"
	"github.com/mizanmahi/aiusage/cli/internal/state"
	"github.com/mizanmahi/aiusage/types"
	"github.com/spf13/cobra"
)

var pushDryRun bool

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push pending Claude Code and Codex usage sessions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPush(cmd.OutOrStdout(), cmd.ErrOrStderr(), pushDryRun)
	},
}

func init() {
	pushCmd.Flags().BoolVar(&pushDryRun, "dry-run", false, "Preview pending sessions without sending data")
	rootCmd.AddCommand(pushCmd)
}

func runPush(out, errOut io.Writer, dryRun bool) error {
	if !dryRun {
		return fmt.Errorf("push is not implemented yet; use --dry-run to preview pending sessions")
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config not found; run 'aiusage init' first: %w", err)
	}

	st, err := state.Load()
	if err != nil {
		return err
	}

	events, claudeCount, codexCount := collectPendingEvents(cfg, st.LastPushedAt, errOut)
	printDryRun(out, cfg.ServerURL, st.LastPushedAt, events, claudeCount, codexCount)
	return nil
}

func collectPendingEvents(cfg *config.Config, since time.Time, errOut io.Writer) ([]types.UsageEvent, int, int) {
	var events []types.UsageEvent

	claudeSessions, err := claude.ReadSessions(cfg.ClaudePath, since)
	if err != nil {
		fmt.Fprintf(errOut, "warning: could not read Claude sessions: %v\n", err)
	}
	for _, session := range claudeSessions {
		events = append(events, session.ToUsageEvent(""))
	}

	codexSessions, err := codex.ReadSessions(cfg.CodexPath, since)
	if err != nil {
		fmt.Fprintf(errOut, "warning: could not read Codex sessions: %v\n", err)
	}
	for _, session := range codexSessions {
		events = append(events, session.ToUsageEvent(""))
	}

	return events, len(claudeSessions), len(codexSessions)
}

func printDryRun(out io.Writer, serverURL string, since time.Time, events []types.UsageEvent, claudeCount, codexCount int) {
	fmt.Fprintln(out, "Dry run: no data sent.")
	fmt.Fprintf(out, "Server URL: %s\n", serverURL)
	if since.IsZero() {
		fmt.Fprintln(out, "Last push: never")
	} else {
		fmt.Fprintf(out, "Last push: %s\n", since.Format(time.RFC3339))
	}
	fmt.Fprintf(out, "Pending sessions: %d (%d Claude, %d Codex)\n", len(events), claudeCount, codexCount)

	if len(events) == 0 {
		return
	}

	fmt.Fprintln(out, "Sessions:")
	for _, event := range events {
		fmt.Fprintf(out, "- %s %s project=%s date=%s model=%s input=%d output=%d cache=%d reasoning=%d\n",
			event.Tool,
			event.SessionID,
			event.Project,
			event.Date,
			event.Model,
			event.InputTokens,
			event.OutputTokens,
			event.CacheTokens,
			event.ReasoningTokens,
		)
	}
}
