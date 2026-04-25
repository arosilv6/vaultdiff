package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/vault"
)

var templateMount string

func init() {
	templateCmd := &cobra.Command{
		Use:   "template [template-string]",
		Short: "Render a template string by substituting Vault secret values",
		Long: `Render a template string replacing {{ secret "path" "key" }} directives
with live values fetched from Vault.

Example:
  vaultdiff template 'DB_PASS={{ secret "db/creds" "password" }}'`,
		Args: cobra.ExactArgs(1),
		RunE: runTemplate,
	}

	templateCmd.Flags().StringVar(&templateMount, "mount", "secret", "KV mount path")
	rootCmd.AddCommand(templateCmd)
}

func runTemplate(cmd *cobra.Command, args []string) error {
	client, err := vault.NewClient()
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	tr := vault.NewTemplater(client, templateMount)
	result, err := tr.Render(cmd.Context(), args[0])
	if err != nil {
		return fmt.Errorf("rendering template: %w", err)
	}

	fmt.Fprintln(os.Stdout, result.Rendered)

	if len(result.SecretsUsed) > 0 {
		fmt.Fprintf(os.Stderr, "# secrets used: %v\n", result.SecretsUsed)
	}

	return nil
}
