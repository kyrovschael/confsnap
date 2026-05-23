package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"confsnap/internal/snapshot"
	"confsnap/internal/snapshot/store"
	"confsnap/internal/watch"
)

var (
	watchInterval int
	watchLabel    string
)

func init() {
	watchCmd := &cobra.Command{
		Use:   "watch [files...]",
		Short: "Watch config files and auto-snapshot on change",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runWatch,
	}
	watchCmd.Flags().IntVarP(&watchInterval, "interval", "i", 30, "poll interval in seconds")
	watchCmd.Flags().StringVarP(&watchLabel, "label", "l", "watch", "label prefix for auto snapshots")
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	interval := time.Duration(watchInterval) * time.Second
	w, err := watch.New(args, interval)
	if err != nil {
		return fmt.Errorf("initialising watcher: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Watching %d file(s) every %s. Press Ctrl+C to stop.\n", len(args), interval)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-sig:
			fmt.Fprintln(os.Stdout, "\nStopping watcher.")
			return nil
		case <-ticker.C:
			events, err := w.Poll()
			if err != nil {
				fmt.Fprintf(os.Stderr, "poll error: %v\n", err)
				continue
			}
			for _, ev := range events {
				label := fmt.Sprintf("%s-%d", watchLabel, ev.ChangedAt.Unix())
				snap, serr := snapshot.Capture(ev.Path)
				if serr != nil {
					fmt.Fprintf(os.Stderr, "snapshot error for %s: %v\n", ev.Path, serr)
					continue
				}
				if serr = store.Save(label, []snapshot.Entry{snap}); serr != nil {
					fmt.Fprintf(os.Stderr, "save error: %v\n", serr)
					continue
				}
				fmt.Fprintf(os.Stdout, "[%s] change detected in %s → snapshot %q saved\n",
					ev.ChangedAt.Format(time.RFC3339), ev.Path, label)
			}
		}
	}
}
