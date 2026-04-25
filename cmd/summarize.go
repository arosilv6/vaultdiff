package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/vault"
)

var (
	summarizeMountPath string
	summarizeTopN      int
	summarizeShowEmpty bool
)

func init() {
	summarizeCmd := &cobra.Command{
		Use:   "summarize [path]",
		Short: "Summarize secret activity and version statistics at a given path",
		Long: `Summarize displays an overview of secret version counts, modification
frequency, and staleness for all secrets found under the given KV path.

Examples:
  vaultdiff summarize secret/myapp
  vaultdiff summarize secret/myapp --top 10
  vaultdiff summarize secret/ --show-empty`,
		Args: cobra.ExactArgs(1),
		RunE: runSummarize,
	}

	summarizeCmd.Flags().StringVar(&summarizeMountPath, "mount", "secret", "KV v2 mount path")
	summarizeCmd.Flags().IntVar(&summarizeTopN, "top", 0, "Limit output to the top N secrets by version count (0 = all)")
	summarizeCmd.Flags().BoolVar(&summarizeShowEmpty, "show-empty", false, "Include secrets with no active versions")

	rootCmd.AddCommand(summarizeCmd)
}

func runSummarize(cmd *cobra.Command, args []string) error {
	path := args[0]

	client, err := vault.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	summarizer := vault.NewSummarizer(client, summarizeMountPath)

	entries, err := summarizer.Summarize(cmd.Context(), path)
	if err != nil {
		return fmt.Errorf("summarize failed: %w", err)
	}

	if !summarizeShowEmpty {
		filtered := entries[:0]
		for _, e := range entries {
			if e.VersionCount > 0 {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}

	if summarizeTopN > 0 && summarizeTopN < len(entries) {
		entries = entries[:summarizeTopN]
	}

	if len(entries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No secrets found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PATH\tVERSIONS\tCURRENT\tDESTROYED\tLAST MODIFIED")
	fmt.Fprintln(w, "----\t--------\t-------\t---------\t-------------")

	for _, e := range entries {
		lastMod := "—"
		if !e.LastModified.IsZero() {
			lastMod = e.LastModified.Format("2006-01-02 15:04:05")
		}
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%s\n",
			e.Path,
			e.VersionCount,
			e.CurrentVersion,
			e.DestroyedCount,
			lastMod,
		)
	}

	return w.Flush()
}
