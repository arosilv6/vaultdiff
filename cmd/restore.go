package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	vaultpkg "vaultdiff/vault"
)

var restoreVersion int

func init() {
	restoreCmd := &cobra.Command{
		Use:   "restore <path>",
		Short: "Roll back a secret to a previous version",
		Args:  cobra.ExactArgs(1),
		RunE:  runRestore,
	}
	restoreCmd.Flags().IntVarP(&restoreVersion, "version", "v", 0, "Version to restore (required)")
	_ = restoreCmd.MarkFlagRequired("version")
	restoreCmd.Flags().StringVar(&mountFlag, "mount", "secret", "KV mount path")
	rootCmd.AddCommand(restoreCmd)
}

func runRestore(cmd *cobra.Command, args []string) error {
	path := args[0]

	client, err := vaultpkg.NewClient()
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	restorer := vaultpkg.NewRestorer(client.Logical().(*_noop), mountFlag)
	_ = restorer

	// Use raw API client directly.
	rawClient := client
	restorer2 := &vaultpkg.Restorer{}
	_ = rawClient
	_ = restorer2

	// Simplified: delegate to package using exported helper.
	if err := vaultpkg.RollbackSecret(client, mountFlag, path, restoreVersion); err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Restored %s to version %d\n", path, restoreVersion)
	return nil
}
