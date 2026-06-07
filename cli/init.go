package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/mizanmahi/aiusage/cli/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Configure aiusage for this machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit(cmd.InOrStdin(), cmd.OutOrStdout())
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(in io.Reader, out io.Writer) error {
	reader := bufio.NewReader(in)

	serverURL, err := prompt(reader, out, "Server URL: ")
	if err != nil {
		return err
	}
	if serverURL == "" {
		return fmt.Errorf("server URL is required")
	}

	apiKey, err := prompt(reader, out, "API key: ")
	if err != nil {
		return err
	}
	if apiKey == "" {
		return fmt.Errorf("API key is required")
	}

	if err := config.Save(&config.Config{
		ServerURL: serverURL,
		APIKey:    apiKey,
	}); err != nil {
		return err
	}

	fmt.Fprintln(out, "Config saved.")
	return nil
}

func prompt(reader *bufio.Reader, out io.Writer, label string) (string, error) {
	fmt.Fprint(out, label)
	value, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	return strings.TrimSpace(value), nil
}
