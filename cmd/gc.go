package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/user/vaultdiff/vault"
)

func init() {
	var mount string
	var maxVersions int

	gcCmd := &cobra.Command{
		Use:   "gc <path>",
		Short: "Garbage collect destroyed secret versions",
		Long: `Permanently purge destroyed versions at a KV v2 path.
Versions that are already marked destroyed and fall within the
retention window (1..maxVersions) are deleted from Vault storage.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGC(args[0], mount, maxVersions)
		},
	}

	gcCmd.Flags().StringVar(&mount, "mount", "secret", "KV v2 mount path")
	gcCmd.Flags().IntVar(&maxVersions, "max-versions", 10, "Maximum version number to consider for purge")

	rootCmd.AddCommand(gcCmd)
}

func runGC(path, mount string, maxVersions int) error {
	client, err := vault.NewClient()
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	gc := vault.NewGarbageCollector(client, mount)
	result, err := gc.Collect(rootCmd.Context(), path, maxVersions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gc error: %v\n", err)
		return err
	}

	if result.Count == 0 {
		fmt.Printf("No destroyed versions to purge at %s\n", path)
		return nil
	}

	purged := make([]string, len(result.VersionsPurged))
	for i, v := range result.VersionsPurged {
		purged[i] = fmt.Sprintf("%d", v)
	}
	fmt.Printf("Purged %d destroyed version(s) at %s: [%s]\n",
		result.Count, result.Path, strings.Join(purged, ", "))
	return nil
}
