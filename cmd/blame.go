package cmd

import (
	"context"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"

	vaultpkg "vaultdiff/vault"
)

var blameCmd = &cobra.Command{
	Use:   "blame <path>",
	Short: "Show authorship history for each version of a secret",
	Args:  cobra.ExactArgs(1),
	RunE:  runBlame,
}

func init() {
	blameCmd.Flags().String("mount", "secret", "KV v2 mount path")
	rootCmd.AddCommand(blameCmd)
}

func runBlame(cmd *cobra.Command, args []string) error {
	path := args[0]
	mount, _ := cmd.Flags().GetString("mount")

	client, err := vaultpkg.NewClient()
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	blamer := vaultpkg.NewBlamer(client, mount)
	entries, err := blamer.Blame(context.Background(), path)
	if err != nil {
		return fmt.Errorf("blame %s: %w", path, err)
	}

	if len(entries) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "no versions found for %s\n", path)
		return nil
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Version < entries[j].Version
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "VERSION\tCREATED\tAUTHOR\tOPERATION\tDELETED")
	for _, e := range entries {
		deletedStr := ""
		if e.Deleted {
			deletedStr = "yes"
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
			e.Version,
			e.CreatedTime.Format("2006-01-02 15:04:05"),
			e.CreatedBy,
			e.Operation,
			deletedStr,
		)
	}
	return w.Flush()
}
