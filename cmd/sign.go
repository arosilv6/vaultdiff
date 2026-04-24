package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/vault"
)

var (
	signMount   string
	signVersion int
	signHMACKey string
)

func init() {
	signCmd := &cobra.Command{
		Use:   "sign <path>",
		Short: "Compute an HMAC-SHA256 signature for a secret version",
		Long: `Reads a KV v2 secret and computes a deterministic HMAC-SHA256
signature over its key/value pairs. Use this to detect out-of-band
tampering by comparing signatures across time or environments.`,
		Args: cobra.ExactArgs(1),
		RunE: runSign,
	}

	signCmd.Flags().StringVar(&signMount, "mount", "secret", "KV v2 mount path")
	signCmd.Flags().IntVar(&signVersion, "version", 0, "Secret version (0 = latest)")
	signCmd.Flags().StringVar(&signHMACKey, "key", "", "HMAC key used to sign the secret data (required)")
	_ = signCmd.MarkFlagRequired("key")

	rootCmd.AddCommand(signCmd)
}

func runSign(cmd *cobra.Command, args []string) error {
	path := args[0]

	client, err := vault.NewClient()
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	signer := vault.NewSigner(client)
	result, err := signer.Sign(signMount, path, signVersion, signHMACKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return err
	}

	fmt.Printf("path:      %s\n", result.Path)
	fmt.Printf("version:   %d\n", result.Version)
	fmt.Printf("signature: %s\n", result.Signature)
	return nil
}
