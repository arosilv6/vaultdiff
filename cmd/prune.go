package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/vault"
)

var (
	pruneMount  string
	pruneKeep   int
	pruneDryRun bool
)

func init() {
	pruneCmd := &cobra.Command{
		Use:   "prune <secret-path>",
		Short: "Remove old secret versions beyond a keep threshold",
		Args:  cobra.ExactArgs(1),
		RunE:  runPrune,
	}
	pruneCmd.Flags().StringVar(&pruneMount, "mount", "secret", "KV v2 mount path")
	pruneCmd.Flags().IntVar(&pruneKeep, "keep", 5, "Number of recent versions to keep")
	pruneCmd.Flags().BoolVar(&pruneDryRun, "dry-run", false, "Preview changes without deleting")
	rootCmd.AddCommand(pruneCmd)
}

func runPrune(cmd *cobra.Command, args []string) error {
	secretPath := args[0]

	client, err := vault.NewClient()
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	pruner := vault.NewPruner(client, pruneMount)
	result, err := pruner.Prune(context.Background(), secretPath, pruneKeep, pruneDryRun)
	if err != nil {
		return fmt.Errorf("pruning versions: %w", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	if pruneDryRun {
		fmt.Fprintln(w, "[dry-run] No changes applied.")
	}

	fmt.Fprintf(w, "Path:\t%s\n", result.Path)
	fmt.Fprintf(w, "Versions kept:\t%d\n", result.VersionsKept)

	if len(result.VersionsPruned) == 0 {
		fmt.Fprintln(w, "Versions pruned:\tnone")
	} else {
		fmt.Fprintf(w, "Versions pruned:\t%v\n", result.VersionsPruned)
	}

	return nil
}
