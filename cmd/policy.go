package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	vaultpkg "github.com/vaultdiff/vault"
)

var policyPath string

func init() {
	policyCmd := &cobra.Command{
		Use:   "policy",
		Short: "Check token capabilities for a secret path",
		RunE:  runPolicy,
	}

	policyCmd.Flags().StringVar(&policyPath, "path", "", "Vault secret path to check (required)")
	_ = policyCmd.MarkFlagRequired("path")

	rootCmd.AddCommand(policyCmd)
}

func runPolicy(cmd *cobra.Command, args []string) error {
	client, err := vaultpkg.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	checker := vaultpkg.NewPolicyChecker(client)

	result, err := checker.CheckPath(context.Background(), policyPath)
	if err != nil {
		return fmt.Errorf("policy check failed: %w", err)
	}

	if len(result.Capabilities) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "path: %s\ncapabilities: none\n", result.Path)
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "path: %s\ncapabilities: %s\n",
		result.Path,
		strings.Join(result.Capabilities, ", "),
	)

	for _, required := range []string{"read", "update"} {
		status := "no"
		if result.HasCapability(required) {
			status = "yes"
		}
		fmt.Fprintf(cmd.OutOrStdout(), "  %-8s: %s\n", required, status)
	}

	return nil
}
