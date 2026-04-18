package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	vaultpkg "vaultdiff/vault"
)

var (
	searchMount   string
	searchPattern string
)

func init() {
	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search for secrets containing a key matching a pattern",
		RunE:  runSearch,
	}
	searchCmd.Flags().StringVar(&searchMount, "mount", "secret", "KV v2 mount path")
	searchCmd.Flags().StringVar(&searchPattern, "pattern", "", "Key pattern to search for (required)")
	_ = searchCmd.MarkFlagRequired("pattern")
	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	client, err := vaultpkg.NewClient()
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	searcher := vaultpkg.NewSearcher(client)
	results, err := searcher.FindByKey(cmd.Context(), searchMount, searchPattern)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if len(results) == 0 {
		fmt.Fprintln(os.Stdout, "No secrets found matching pattern:", searchPattern)
		return nil
	}

	fmt.Fprintf(os.Stdout, "%-50s  %s\n", "PATH", "VERSION")
	fmt.Fprintf(os.Stdout, "%-50s  %s\n", "----", "-------")
	for _, r := range results {
		fmt.Fprintf(os.Stdout, "%-50s  %d\n", r.Path, r.Version)
	}
	return nil
}
