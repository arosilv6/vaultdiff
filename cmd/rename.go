package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/user/vaultdiff/vault"
)

var renameMount string

func init() {
	renameCmd := &cobra.Command{
		Use:   "rename <src-path> <dst-path>",
		Short: "Rename (move) a secret from one path to another",
		Args:  cobra.ExactArgs(2),
		RunE:  runRename,
	}

	renameCmd.Flags().StringVar(&renameMount, "mount", "secret", "KV v2 mount path")
	rootCmd.AddCommand(renameCmd)
}

func runRename(cmd *cobra.Command, args []string) error {
	srcPath := args[0]
	dstPath := args[1]

	client, err := vault.NewClient()
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	renamer := vault.NewRenamer(client, renameMount)
	if err := renamer.Rename(srcPath, dstPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	fmt.Printf("Renamed %q -> %q\n", srcPath, dstPath)
	return nil
}
