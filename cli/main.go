package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cliVersion = "dev"

var rootCmd = &cobra.Command{
	Use:   "aiusage",
	Short: "Push Claude Code and Codex usage analytics to your server",
	Long:  "aiusage is a command-line tool for pushing Claude Code and Codex usage analytics to your server. It allows you to track and analyze your usage of these AI models, providing insights into your interactions and helping you optimize your usage.",
}

func main() {
	rootCmd.Version = cliVersion
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
