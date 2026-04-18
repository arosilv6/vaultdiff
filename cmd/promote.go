package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/vault"
)

var promoteMount string

func init() {
	promoteCmd := &cobra.Command{
		Use:   "promote <path> <version>",
		Short: "Promote an older secret version to the latest",
		Args:  cobra.ExactArgs(2),
		RunE:  runPromote,
	}
	promoteCmd.Flags().StringVar(&promoteMount, "mount", "secret", "KV v2 mount path")
	rootCmd.AddCommand(promoteCmd)
}

func runPromote(cmd *cobra.Command, args []string) error {
	path := args[0]
	version, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid version %q: %w", args[1], err)
	}

	client, err := vault.NewClient()
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	promoter := vault.NewPromoter(client, promoteMount)
	data, err := promoter.Promote(path, version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	fmt.Printf("Promoted version %d of %s/%s to latest.\n", version, promoteMount, path)
	fmt.Printf("Keys written: %d\n", len(data))
	return nil
}
