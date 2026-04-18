package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/user/vaultdiff/audit"
	"github.com/user/vaultdiff/diff"
	"github.com/user/vaultdiff/vault"
)

var auditLogFile string

var auditCmd = &cobra.Command{
	Use:   "audit <path> <versionA> <versionB>",
	Short: "Diff two secret versions and record the result to an audit log",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		vA, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid versionA: %w", err)
		}
		vB, err := strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("invalid versionB: %w", err)
		}

		client, err := vault.NewClient()
		if err != nil {
			return err
		}

		secretA, err := client.GetSecretVersion(path, vA)
		if err != nil {
			return err
		}
		secretB, err := client.GetSecretVersion(path, vB)
		if err != nil {
			return err
		}

		changes := diff.Compare(secretA, secretB)
		diff.Print(changes)

		logger, err := audit.NewLogger(auditLogFile)
		if err != nil {
			return err
		}
		defer logger.Close()

		if err := logger.Record(path, vA, vB, changes); err != nil {
			fmt.Fprintf(os.Stderr, "warning: audit log write failed: %v\n", err)
		}
		return nil
	},
}

func init() {
	auditCmd.Flags().StringVar(&auditLogFile, "log", "", "Path to audit log file (default: stdout)")
	rootCmd.AddCommand(auditCmd)
}
