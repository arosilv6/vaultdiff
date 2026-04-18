package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	vaultpkg "vaultdiff/vault"
)

func init() {
	var mount string
	var version int

	cmd := &cobra.Command{
		Use:   "expire <path>",
		Short: "Check expiry/TTL information for a secret version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExpire(args[0], mount, version)
		},
	}

	cmd.Flags().StringVar(&mount, "mount", "secret", "KV mount path")
	cmd.Flags().IntVar(&version, "version", 1, "Secret version to check")

	rootCmd.AddCommand(cmd)
}

func runExpire(path, mount string, version int) error {
	client, err := vaultpkg.NewClient()
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	expirer := vaultpkg.NewExpirer(client, mount)
	info, err := expirer.CheckExpiry(context.Background(), path, version)
	if err != nil {
		return fmt.Errorf("checking expiry: %w", err)
	}

	if info.ExpiresAt.IsZero() {
		fmt.Fprintf(os.Stdout, "path=%s version=%d no expiry set\n", info.Path, info.Version)
		return nil
	}

	status := "active"
	if info.Expired {
		status = "EXPIRED"
	}

	fmt.Fprintf(os.Stdout, "path=%s version=%d expires_at=%s ttl=%s status=%s\n",
		info.Path, info.Version, info.ExpiresAt.Format("2006-01-02T15:04:05Z"), info.TTL.Round(1e9), status)
	return nil
}
