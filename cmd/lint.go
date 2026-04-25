package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	vaultpkg "github.com/yourusername/vaultdiff/vault"
)

var lintMount string

func init() {
	lintCmd := &cobra.Command{
		Use:   "lint <path>",
		Short: "Lint a secret path for common issues",
		Args:  cobra.ExactArgs(1),
		RunE:  runLint,
	}
	lintCmd.Flags().StringVar(&lintMount, "mount", "secret", "KV mount path")
	rootCmd.AddCommand(lintCmd)
}

func runLint(cmd *cobra.Command, args []string) error {
	path := args[0]

	client, err := vaultpkg.NewClient()
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	linter := vaultpkg.NewLinter(client, lintMount)
	result, err := linter.Lint(cmd.Context(), path)
	if err != nil {
		return fmt.Errorf("lint: %w", err)
	}

	if len(result.Issues) == 0 {
		fmt.Fprintf(os.Stdout, "✓ %s — no issues found\n", path)
		return nil
	}

	fmt.Fprintf(os.Stdout, "Issues in %s:\n", path)
	hasError := false
	for _, issue := range result.Issues {
		prefix := "⚠"
		if issue.Severity == "error" {
			prefix = "✗"
			hasError = true
		}
		if issue.Key != "" {
			fmt.Fprintf(os.Stdout, "  %s [%s] key=%q: %s\n", prefix, issue.Severity, issue.Key, issue.Message)
		} else {
			fmt.Fprintf(os.Stdout, "  %s [%s] %s\n", prefix, issue.Severity, issue.Message)
		}
	}

	if hasError {
		return fmt.Errorf("lint found errors in %s", path)
	}
	return nil
}
