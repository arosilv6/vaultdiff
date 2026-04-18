package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/user/vaultdiff/vault"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback [path] [version]",
	Short: "Roll back a secret to a previous version",
	Args:  cobra.ExactArgs(2),
	RunE:  runRollback,
}

func init() {
	rollbackCmd.Flags().String("mount", "secret", "KV v2 mount path")
	rootCmd.AddCommand(rollbackCmd)
}

func runRollback(cmd *cobra.Command, args []string) error {
	path := args[0]
	var version int
	if _, err := fmt.Sscanf(args[1], "%d", &version); err != nil {
		return fmt.Errorf("invalid version: %s", args[1])
	}

	mount, _ := cmd.Flags().GetString("mount")

	client, err := vault.NewClient()
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	rollbacker := vault.NewRollbacker(client, mount)
	result, err := rollbacker.Rollback(context.Background(), path, version)
	if err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Rolled back %s to version %d (new version: %d)\n",
		result.Path, result.ToVersion, result.FromVersion)
	return nil
}
