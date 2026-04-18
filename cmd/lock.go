package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultdiff/vault"
)

var lockCmd = &cobra.Command{
	Use:   "lock [path]",
	Short: "Lock a secret path to prevent writes",
	Args:  cobra.ExactArgs(1),
	RunE:  runLock,
}

var unlockCmd = &cobra.Command{
	Use:   "unlock [path]",
	Short: "Unlock a previously locked secret path",
	Args:  cobra.ExactArgs(1),
	RunE:  runUnlock,
}

func init() {
	lockCmd.Flags().String("mount", "secret", "KV mount path")
	lockCmd.Flags().String("reason", "", "Reason for locking the secret")

	unlockCmd.Flags().String("mount", "secret", "KV mount path")

	rootCmd.AddCommand(lockCmd)
	rootCmd.AddCommand(unlockCmd)
}

func runLock(cmd *cobra.Command, args []string) error {
	path := args[0]

	mount, err := cmd.Flags().GetString("mount")
	if err != nil {
		return fmt.Errorf("failed to read mount flag: %w", err)
	}

	reason, err := cmd.Flags().GetString("reason")
	if err != nil {
		return fmt.Errorf("failed to read reason flag: %w", err)
	}

	client, err := vault.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	locker := vault.NewLocker(client, mount)
	if err := locker.Lock(cmd.Context(), path, reason); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to lock %s: %v\n", path, err)
		return err
	}

	fmt.Printf("locked: %s/%s\n", mount, path)
	return nil
}

func runUnlock(cmd *cobra.Command, args []string) error {
	path := args[0]

	mount, err := cmd.Flags().GetString("mount")
	if err != nil {
		return fmt.Errorf("failed to read mount flag: %w", err)
	}

	client, err := vault.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	locker := vault.NewLocker(client, mount)
	if err := locker.Unlock(cmd.Context(), path); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to unlock %s: %v\n", path, err)
		return err
	}

	fmt.Printf("unlocked: %s/%s\n", mount, path)
	return nil
}
