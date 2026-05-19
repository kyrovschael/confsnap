package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/confsnap/confsnap/internal/snapshot"
	"github.com/spf13/cobra"
)

var (
	snapLabel      string
	snapOutputDir  string
	snapWithContent bool
)

var snapCmd = &cobra.Command{
	Use:   "snap [files...]",
	Short: "Capture a snapshot of one or more config files",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		set := &snapshot.SnapshotSet{
			Label:     snapLabel,
			CreatedAt: time.Now().UTC(),
		}

		for _, path := range args {
			snap, err := snapshot.Capture(path, snapWithContent)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: skipping %q: %v\n", path, err)
				continue
			}
			set.Snapshots = append(set.Snapshots, snap)
			fmt.Printf("  captured %-40s  sha256:%s\n", snap.Path, snap.Hash[:12])
		}

		if len(set.Snapshots) == 0 {
			return fmt.Errorf("no files were captured")
		}

		out, err := snapshot.Save(snapOutputDir, set)
		if err != nil {
			return fmt.Errorf("saving snapshot: %w", err)
		}
		fmt.Printf("snapshot saved → %s\n", out)
		return nil
	},
}

func init() {
	snapCmd.Flags().StringVarP(&snapLabel, "label", "l", "snapshot", "label for this snapshot set")
	snapCmd.Flags().StringVarP(&snapOutputDir, "output", "o", ".confsnap", "directory to store snapshot files")
	snapCmd.Flags().BoolVar(&snapWithContent, "content", false, "embed file contents in the snapshot")
	RootCmd.AddCommand(snapCmd)
}
