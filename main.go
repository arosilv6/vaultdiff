package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vaultdiff",
	Short: "Diff and audit changes between HashiCorp Vault secret versions",
	Long:  `vaultdiff is a CLI tool to compare secret versions in HashiCorp Vault KV v`,
}

func main() {
	if err := rootCmd.Execute(); err	fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
