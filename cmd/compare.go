package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/user/vaultdiff/diff"
	"github.com/user/vaultdiff/vault"
)

var compareMount string

func init() {
	compareCmd := &cobra.Command{
		Use:   "compare <path> <versionA> <versionB>",
		Short: "Compare two versions of a Vault secret",
		Args:  cobra.ExactArgs(3),
		RunE:  runCompare,
	}
	compareCmd.Flags().StringVar(&compareMount, "mount", "secret", "KV mount path")
	rootCmd.AddCommand(compareCmd)
}

func runCompare(cmd *cobra.Command, args []string) error {
	path := args[0]
	versionA, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid versionA: %w", err)
	}
	versionB, err := strconv.Atoi(args[2])
	if err != nil {
		return fmt.Errorf("invalid versionB: %w", err)
	}

	client, err := vault.NewClient()
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	comparer := vault.NewComparer(client, compareMount)
	result, err := comparer.FetchVersions(path, versionA, versionB)
	if err != nil {
		return fmt.Errorf("fetching versions: %w", err)
	}

	changes := diff.Compare(result.DataA, result.DataB)
	if len(changes) == 0 {
		fmt.Fprintf(os.Stdout, "No differences between version %d and version %d of %q\n",
			versionA, versionB, path)
		return nil
	}

	fmt.Fprintf(os.Stdout, "Diff for %q (v%d → v%d):\n", path, versionA, versionB)
	diff.Print(os.Stdout, changes)
	return nil
}
