package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	vaultclient "github.com/yourusername/vaultdiff/vault"
)

var pinMount string

func init() {
	pinCmd := &cobra.Command{
		Use:   "pin <path> <version>",
		Short: "Pin a secret version to prevent accidental overwrites",
		Args:  cobra.ExactArgs(2),
		RunE:  runPin,
	}
	pinCmd.Flags().StringVar(&pinMount, "mount", "secret", "KV mount path")

	unpinCmd := &cobra.Command{
		Use:   "unpin <path>",
		Short: "Remove the pinned version marker from a secret",
		Args:  cobra.ExactArgs(1),
		RunE:  runUnpin,
	}
	unpinCmd.Flags().StringVar(&pinMount, "mount", "secret", "KV mount path")

	rootCmd.AddCommand(pinCmd)
	rootCmd.AddCommand(unpinCmd)
}

func runPin(cmd *cobra.Command, args []string) error {
	path := args[0]
	version, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid version %q: %w", args[1], err)
	}

	client, err := vaultclient.NewClient()
	if err != nil {
		return fmt.Errorf("vault client error: %w", err)
	}

	pinner := vaultclient.NewPinner(client, pinMount)
	result, err := pinner.Pin(context.Background(), path, version)
	if err != nil {
		return fmt.Errorf("pin failed: %w", err)
	}

	fmt.Printf("Pinned %s @ version %d\n", result.Path, result.Version)
	return nil
}

func runUnpin(cmd *cobra.Command, args []string) error {
	path := args[0]

	client, err := vaultclient.NewClient()
	if err != nil {
		return fmt.Errorf("vault client error: %w", err)
	}

	pinner := vaultclient.NewPinner(client, pinMount)
	result, err := pinner.Unpin(context.Background(), path)
	if err != nil {
		return fmt.Errorf("unpin failed: %w", err)
	}

	fmt.Printf("Unpinned %s\n", result.Path)
	return nil
}
