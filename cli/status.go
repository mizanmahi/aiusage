package main

import (
	"fmt"
	"io"
	"time"

	"github.com/mizanmahi/aiusage/cli/internal/claude"
	"github.com/mizanmahi/aiusage/cli/internal/codex"
	"github.com/mizanmahi/aiusage/cli/internal/config"
	"github.com/mizanmahi/aiusage/cli/internal/state"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show aiusage configuration and pending session count",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runStatus(cmd.OutOrStdout(), cmd.ErrOrStderr())
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(out, errOut io.Writer) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config not found; run 'aiusage init' first: %w", err)
	}

	st, err := state.Load()
	if err != nil {
		return err
	}

	since := st.LastPushedAt
	claudeSessions, err := claude.ReadSessions(cfg.ClaudePath, since)
	if err != nil {
		fmt.Fprintf(errOut, "warning: could not read Claude sessions: %v\n", err)
	}

	codexSessions, err := codex.ReadSessions(cfg.CodexPath, since)
	if err != nil {
		fmt.Fprintf(errOut, "warning: could not read Codex sessions: %v\n", err)
	}

	fmt.Fprintf(out, "Server URL: %s\n", cfg.ServerURL)
	if since.IsZero() {
		fmt.Fprintln(out, "Last push: never")
	} else {
		fmt.Fprintf(out, "Last push: %s\n", since.Format(time.RFC3339))
	}
	fmt.Fprintf(out, "Pending sessions: %d (%d Claude, %d Codex)\n",
		len(claudeSessions)+len(codexSessions),
		len(claudeSessions),
		len(codexSessions),
	)

	return nil
}
