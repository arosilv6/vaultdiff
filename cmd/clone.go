package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/audit"
	"github.com/yourusername/vaultdiff/vault"
)

var cloneCmd = &cobra.Command{
	Use:   "clone <source-path> <dest-path>",
	Short: "Clone a secret path to a new location",
	Long: `Clone copies all versions of a secret from one KV path to another.

By default, all versions are cloned. Use --version to clone a specific version.
The destination path must not already exist unless --force is specified.`,
	Args: cobra.ExactArgs(2),
	RunE: runClone,
}

func init() {
	cloneCmd.Flags().IntP("version", "v", 0, "specific version to clone (0 = all versions)")
	cloneCmd.Flags().StringP("mount", "m", "secret", "KV secrets engine mount path")
	cloneCmd.Flags().BoolP("force", "f", false, "overwrite destination if it already exists")
	cloneCmd.Flags().StringP("audit-log", "a", "", "path to audit log file (default: stdout)")

	rootCmd.AddCommand(cloneCmd)
}

func runClone(cmd *cobra.Command, args []string) error {
	srcPath := args[0]
	dstPath := args[1]

	version, err := cmd.Flags().GetInt("version")
	if err != nil {
		return fmt.Errorf("invalid version flag: %w", err)
	}

	mount, err := cmd.Flags().GetString("mount")
	if err != nil {
		return fmt.Errorf("invalid mount flag: %w", err)
	}

	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		return fmt.Errorf("invalid force flag: %w", err)
	}

	auditPath, err := cmd.Flags().GetString("audit-log")
	if err != nil {
		return fmt.Errorf("invalid audit-log flag: %w", err)
	}

	logger, err := audit.NewLogger(auditPath)
	if err != nil {
		return fmt.Errorf("failed to initialize audit logger: %w", err)
	}

	client, err := vault.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	cloner := vault.NewCloner(client, mount)

	result, err := cloner.Clone(cmd.Context(), srcPath, dstPath, version, force)
	if err != nil {
		return fmt.Errorf("clone failed: %w", err)
	}

	// Audit the clone operation.
	logger.Record(audit.Entry{
		Operation: "clone",
		Path:      srcPath,
		Metadata: map[string]string{
			"destination":      dstPath,
			"versions_cloned":  fmt.Sprintf("%d", result.VersionsCloned),
			"mount":            mount,
			"force":            fmt.Sprintf("%v", force),
		},
	})

	fmt.Fprintf(os.Stdout, "Cloned %q → %q (%d version(s) copied)\n",
		srcPath, dstPath, result.VersionsCloned)

	return nil
}
