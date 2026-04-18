package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	vaultpkg "vaultdiff/vault"
)

var snapshotOutput string

func init() {
	snapshotCmd := &cobra.Command{
		Use:   "snapshot [base-path]",
		Short: "Capture a snapshot of all secrets under a path",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runSnapshot,
	}
	snapshotCmd.Flags().StringVarP(&snapshotOutput, "output", "o", "", "Write snapshot JSON to file (default: stdout)")
	snapshotCmd.Flags().StringVar(&mountFlag, "mount", "secret", "KV v2 mount path")
	rootCmd.AddCommand(snapshotCmd)
}

func runSnapshot(cmd *cobra.Command, args []string) error {
	basePath := ""
	if len(args) > 0 {
		basePath = args[0]
	}

	client, err := vaultpkg.NewClient()
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	s := vaultpkg.NewSnapshotter(client, mountFlag)
	entries, err := s.Snapshot(basePath)
	if err != nil {
		return fmt.Errorf("snapshot: %w", err)
	}

	out, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding snapshot: %w", err)
	}

	if snapshotOutput != "" {
		if err := os.WriteFile(snapshotOutput, out, 0600); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Snapshot written to %s (%d entries)\n", snapshotOutput, len(entries))
	} else {
		fmt.Fprintln(cmd.OutOrStdout(), string(out))
	}
	return nil
}
