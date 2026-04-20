package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultdiff/vault"
)

var (
	importMount  string
	importDryRun bool
)

func init() {
	importCmd := &cobra.Command{
		Use:   "import <file>",
		Short: "Import secrets from a JSON export file into Vault",
		Args:  cobra.ExactArgs(1),
		RunE:  runImport,
	}

	importCmd.Flags().StringVar(&importMount, "mount", "secret", "KV v2 mount path")
	importCmd.Flags().BoolVar(&importDryRun, "dry-run", false, "Simulate import without writing to Vault")

	rootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	client, err := vault.NewClient()
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	importer := vault.NewImporter(client, importMount)
	result, err := importer.ImportFile(context.Background(), filePath, importDryRun)
	if err != nil {
		return fmt.Errorf("import: %w", err)
	}

	if importDryRun {
		fmt.Printf("[dry-run] would import %d secret(s)\n", result.Imported)
		return nil
	}

	fmt.Printf("imported: %d  skipped: %d\n", result.Imported, result.Skipped)
	if len(result.Errors) > 0 {
		fmt.Fprintln(os.Stderr, "errors:")
		for _, e := range result.Errors {
			fmt.Fprintf(os.Stderr, "  - %s\n", e)
		}
	}
	return nil
}
