package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultdiff/vault"
)

var (
	encryptMount   string
	encryptKey     string
	encryptDecrypt bool
)

func init() {
	encryptCmd := &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt or decrypt a value using the Vault Transit engine",
		Args:  cobra.ExactArgs(1),
		RunE:  runEncrypt,
	}
	encryptCmd.Flags().StringVar(&encryptMount, "mount", "transit", "Transit secrets engine mount path")
	encryptCmd.Flags().StringVar(&encryptKey, "key", "", "Transit key name (required)")
	encryptCmd.Flags().BoolVar(&encryptDecrypt, "decrypt", false, "Decrypt instead of encrypt")
	_ = encryptCmd.MarkFlagRequired("key")
	rootCmd.AddCommand(encryptCmd)
}

func runEncrypt(cmd *cobra.Command, args []string) error {
	client, err := vault.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	enc := vault.NewEncrypter(client, encryptMount, encryptKey)
	ctx := context.Background()

	if encryptDecrypt {
		res, err := enc.Decrypt(ctx, args[0])
		if err != nil {
			return fmt.Errorf("decrypt failed: %w", err)
		}
		decoded, err := base64.StdEncoding.DecodeString(res.Plaintext)
		if err != nil {
			// Print raw base64 if decoding fails
			fmt.Fprintln(os.Stdout, res.Plaintext)
			return nil
		}
		fmt.Fprintln(os.Stdout, string(decoded))
		return nil
	}

	b64 := base64.StdEncoding.EncodeToString([]byte(args[0]))
	res, err := enc.Encrypt(ctx, b64)
	if err != nil {
		return fmt.Errorf("encrypt failed: %w", err)
	}
	fmt.Fprintln(os.Stdout, res.Ciphertext)
	return nil
}
