package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"vaultd/vault"
)
var interval

	watchCmd := &cobra.Command{
		Use:   "watch [path]",
		Sh a Vault secret path for version changes",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWatch(args[0], interval)
		},
	}

	watchCmd.Flags().IntVarP(&interval, "interval", "i", 30, "Poll interval in seconds")
	rootCmd.AddCommand(watchCmd)
}

func runWatch(path string, intervalSec int) error {
	client, err := vault.NewClient()
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	duration := time.Duration(intervalSec) * time.Second
	watcher := vault.NewWatcher(client, path, duration)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		cancel()
	}()

	fmt.Fprintf(os.Stdout, "Watching %s every %v...\n", path, duration)
	for event := range watcher.Watch(ctx) {
		if event.Error != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", event.Error)
			continue
		}
		fmt.Fprintf(os.Stdout, "[%s] version changed: %d -> %d\n",
			event.Path, event.OldVersion, event.NewVersion)
	}
	return nil
}
