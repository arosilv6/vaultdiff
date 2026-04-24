package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/vault"
)

var (
	archiveVersion  int
	archiveDestPath string
)

var archiveCmd = &cobra.Command{
	Use:   "archive <path>",
	Short: "Archive a specific version of a secret to a designated path",
	Args:  cobra.ExactArgs(1),
	RunE:  runArchive,
}

func init() {
	archiveCmd.Flags().IntVar(&archiveVersion, "version", 0, "Version of the secret to archive (required)")
	archiveCmd.Flags().StringVar(&archiveDestPath, "dest", "", "Destination archive path (required)")
	_ = archiveCmd.MarkFlagRequired("version")
	_ = archiveCmd.MarkFlagRequired("dest")
	archiveCmd.Flags().StringVar(&mountFlag, "mount", "secret", "KV mount path")
	rootCmd.AddCommand(archiveCmd)
}

func runArchive(cmd *cobra.Command, args []string) error {
	srcPath := args[0]

	client, err := vault.NewClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating vault client: %v\n", err)
		return err
	}

	archiver := vault.NewArchiver(client, mountFlag)
	entry, err := archiver.Archive(srcPath, archiveVersion, archiveDestPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "archive failed: %v\n", err)
		return err
	}

	fmt.Printf("Archived %s (v%d) → %s at %s\n",
		srcPath,
		entry.Version,
		entry.Path,
		entry.ArchivedAt.Format("2006-01-02T15:04:05Z"),
	)
	return nil
}
