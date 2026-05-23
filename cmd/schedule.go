package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/user/confsnap/internal/schedule"
)

func init() {
	var interval string
	var files []string

	addCmd := &cobra.Command{
		Use:   "add <label>",
		Short: "Add or update a scheduled snapshot entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := time.ParseDuration(interval)
			if err != nil {
				return fmt.Errorf("invalid interval %q: %w", interval, err)
			}
			s, err := schedule.Load()
			if err != nil {
				return err
			}
			s.Add(schedule.Entry{Label: args[0], Files: files, Interval: d})
			if err := schedule.Save(s); err != nil {
				return err
			}
			fmt.Printf("scheduled %s every %s\n", args[0], interval)
			return nil
		},
	}
	addCmd.Flags().StringVarP(&interval, "interval", "i", "1h", "snapshot interval (e.g. 30m, 6h)")
	addCmd.Flags().StringSliceVarP(&files, "files", "f", nil, "comma-separated list of config files")
	_ = addCmd.MarkFlagRequired("files")

	removeCmd := &cobra.Command{
		Use:   "remove <label>",
		Short: "Remove a scheduled snapshot entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := schedule.Load()
			if err != nil {
				return err
			}
			if !s.Remove(args[0]) {
				return fmt.Errorf("label %q not found", args[0])
			}
			return schedule.Save(s)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List scheduled snapshot entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := schedule.Load()
			if err != nil {
				return err
			}
			if len(s.Entries) == 0 {
				fmt.Println("no scheduled entries")
				return nil
			}
			for _, e := range s.Entries {
				fmt.Printf("%-20s %s  [%s]\n", e.Label, e.Interval, strings.Join(e.Files, ", "))
			}
			return nil
		},
	}

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run due scheduled snapshots now",
		RunE: func(cmd *cobra.Command, args []string) error {
			return schedule.RunOnce(os.Stdout, time.Now())
		},
	}

	scheduleCmd := &cobra.Command{
		Use:   "schedule",
		Short: "Manage scheduled snapshot jobs",
	}
	scheduleCmd.AddCommand(addCmd, removeCmd, listCmd, runCmd)
	rootCmd.AddCommand(scheduleCmd)
}
