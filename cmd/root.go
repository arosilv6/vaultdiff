package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vaultdiff",
	Short: "Diff and between HashiCorp Vault secret versions",
}//rintln(os.Stderr, err)
		os.Exit(1)
	}
}
