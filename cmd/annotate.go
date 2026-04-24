package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	vaultpkg "github.com/yourusername/vaultdiff/vault"
)

var (
	annotateMount   string
	annotateAuthor  string
	annotateNote    string
	annotateList    bool
)

func init() {
	annotateCmd := &cobra.Command{
		Use:   "annotate <path> [version]",
		Short: "Annotate a secret version with a note",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  runAnnotate,
	}

	annotateCmd.Flags().StringVar(&annotateMount, "mount", "secret", "KV mount path")
	annotateCmd.Flags().StringVar(&annotateAuthor, "author", "", "Author of the annotation")
	annotateCmd.Flags().StringVar(&annotateNote, "note", "", "Annotation note text")
	annotateCmd.Flags().BoolVar(&annotateList, "list", false, "List all annotations for the path")

	rootCmd.AddCommand(annotateCmd)
}

func runAnnotate(cmd *cobra.Command, args []string) error {
	client, err := vaultpkg.NewClient()
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	path := args[0]
	annotator := vaultpkg.NewAnnotator(client, annotateMount, path)

	if annotateList {
		annotations, err := annotator.GetAnnotations()
		if err != nil {
			return fmt.Errorf("get annotations: %w", err)
		}
		if len(annotations) == 0 {
			fmt.Fprintln(os.Stdout, "No annotations found.")
			return nil
		}
		for _, a := range annotations {
			fmt.Fprintf(os.Stdout, "v%d  [%s]  %s\n", a.Version, a.Author, a.Note)
		}
		return nil
	}

	if len(args) < 2 {
		return fmt.Errorf("version argument required when not using --list")
	}
	version, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid version %q: %w", args[1], err)
	}

	if err := annotator.Annotate(version, annotateAuthor, annotateNote); err != nil {
		return fmt.Errorf("annotate: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Annotation saved for %s version %d\n", path, version)
	return nil
}
