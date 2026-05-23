package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"confsnap/internal/baseline"
	"confsnap/internal/snapshot"
)

func init() {
	var name, label string

	createCmd := &cobra.Command{
		Use:   "baseline create [files...]",
		Short: "Create a named baseline from current file state",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hashes := make(map[string]string)
			for _, path := range args {
				snap, err := snapshot.Capture(path)
				if err != nil {
					return fmt.Errorf("capture %q: %w", path, err)
				}
				hashes[path] = snap.Hash
			}
			b := &baseline.Baseline{
				Name:      name,
				CreatedAt: time.Now(),
				Label:     label,
				Files:     args,
				Hashes:    hashes,
			}
			if err := baseline.Save(b); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Baseline %q created with %d file(s)\n", name, len(args))
			return nil
		},
	}
	createCmd.Flags().StringVarP(&name, "name", "n", "", "Baseline name (required)")
	createCmd.Flags().StringVarP(&label, "label", "l", "", "Optional label or tag")
	_ = createCmd.MarkFlagRequired("name")

	driftCmd := &cobra.Command{
		Use:   "baseline drift <name>",
		Short: "Check current files against a saved baseline",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := baseline.Load(args[0])
			if err != nil {
				return err
			}
			drifts, err := baseline.CheckDrift(b)
			if err != nil {
				return err
			}
			for _, d := range drifts {
				fmt.Fprintf(cmd.OutOrStdout(), "%-12s %s\n", d.Status, d.Path)
			}
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "baseline list",
		Short: "List all saved baselines",
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := baseline.List()
			if err != nil {
				return err
			}
			if len(names) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No baselines found.")
				return nil
			}
			for _, n := range names {
				fmt.Fprintln(cmd.OutOrStdout(), n)
			}
			return nil
		},
	}

	baselineCmd := &cobra.Command{Use: "baseline", Short: "Manage configuration baselines"}
	baselineCmd.AddCommand(createCmd, driftCmd, listCmd)
	rootCmd.AddCommand(baselineCmd)
	_ = os.MkdirAll(".confsnap/baselines", 0o755)
}
