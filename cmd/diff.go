package cmd

import (
	"fmt"

spf13/cobra"
	"github.com/user/vaultdiff/vault"
	"github.com/user/vaultdiff/diff"
)

var (
	flagVersionA int
	flagVersionB int
	flagMount    string
)

var DiffCmd = &cobra.Command{
	Use:   "diff <secret-path>",
	Short: "Diff two versions of a Vault secret",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		client, err := vault.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create vault client: %w", err)
		}

		secretA, err := client.GetSecretVersion(flagMount, path, flagVersionA)
		if err != nil {
			return fmt.Errorf("failed to fetch version %d: %w", flagVersionA, err)
		}

		secretB, err := client.GetSecretVersion(flagMount, path, flagVersionB)
		if err != nil {
			return fmt.Errorf("failed to fetch version %d: %w", flagVersionB, err)
		}

		result := diff.Compare(secretA, secretB)
		diff.Print(os.Stdout, result)
		return nil
	},
}

func init() {
	DiffCmd.Flags().IntVar(&flagVersionA, "version-a", 1, "First version to compare")
	DiffCmd.Flags().IntVar(&flagVersionB, "version-b", 2, "Second version to compare")
	DiffCmd.Flags().StringVar(&flagMount, "mount", "secret", "KV v2 mount path")
}
