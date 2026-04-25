package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/vault"
)

var (
	historyMount   string
	historyPath    string
	historyLimit   int
	historyShowAll bool
)

func init() {
	historyCmd := &cobra.Command{
		Use:   "history",
		Short: "Show the version history of a Vault secret",
		Long: `Display a chronological history of all versions for a given secret path,
including creation time, deletion status, and destruction status.`,
		RunE: runHistory,
	}

	historyCmd.Flags().StringVar(&historyMount, "mount", "secret", "KV v2 mount path")
	historyCmd.Flags().StringVar(&historyPath, "path", "", "Secret path to retrieve history for (required)")
	historyCmd.Flags().IntVar(&historyLimit, "limit", 0, "Maximum number of versions to display (0 = all)")
	historyCmd.Flags().BoolVar(&historyShowAll, "all", false, "Include destroyed and deleted versions")

	_ = historyCmd.MarkFlagRequired("path")

	rootCmd.AddCommand(historyCmd)
}

func runHistory(cmd *cobra.Command, args []string) error {
	client, err := vault.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	historian := vault.NewHistorian(client, historyMount)

	entries, err := historian.History(historyPath)
	if err != nil {
		return fmt.Errorf("failed to retrieve history for %q: %w", historyPath, err)
	}

	if len(entries) == 0 {
		fmt.Printf("No version history found for path: %s\n", historyPath)
		return nil
	}

	// Apply limit if specified
	displayed := entries
	if historyLimit > 0 && historyLimit < len(entries) {
		displayed = entries[:historyLimit]
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "VERSION\tCREATED\tDELETED\tDESTROYED\tMETADATA")
	fmt.Fprintln(w, "-------\t-------\t-------\t---------\t--------")

	for _, entry := range displayed {
		// Skip deleted/destroyed unless --all is set
		if !historyShowAll && (entry.Deleted || entry.Destroyed) {
			continue
		}

		deletedStr := "-"
		if entry.Deleted {
			deletedStr = "yes"
		}

		destroyedStr := "-"
		if entry.Destroyed {
			destroyedStr = "yes"
		}

		metaStr := "-"
		if len(entry.CustomMetadata) > 0 {
			metaStr = fmt.Sprintf("%d key(s)", len(entry.CustomMetadata))
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
			entry.Version,
			entry.CreatedAt.Format("2006-01-02 15:04:05"),
			deletedStr,
			destroyedStr,
			metaStr,
		)
	}

	return w.Flush()
}
