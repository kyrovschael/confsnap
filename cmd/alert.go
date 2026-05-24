package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/nicholasgasior/confsnap/internal/alert"
	"github.com/spf13/cobra"
)

func init() {
	var alertCmd = &cobra.Command{
		Use:   "alert",
		Short: "Manage alert configuration and dispatch test alerts",
	}

	var configureCmd = &cobra.Command{
		Use:   "configure",
		Short: "Configure alert settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			level, _ := cmd.Flags().GetString("level")
			logFile, _ := cmd.Flags().GetString("log-file")
			stdout, _ := cmd.Flags().GetBool("stdout")

			cfg := alert.Config{
				Enabled: true,
				Level:   level,
				LogFile: logFile,
				Stdout:  stdout,
			}
			if err := alert.SaveConfig(cfg); err != nil {
				return fmt.Errorf("saving alert config: %w", err)
			}
			fmt.Println("Alert configuration saved.")
			return nil
		},
	}
	configureCmd.Flags().String("level", "WARNING", "Alert level (INFO|WARNING|CRITICAL)")
	configureCmd.Flags().String("log-file", "", "Path to alert log file (optional)")
	configureCmd.Flags().Bool("stdout", true, "Print alerts to stdout")

	var testCmd = &cobra.Command{
		Use:   "test",
		Short: "Send a test alert using current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := alert.LoadConfig()
			if err != nil {
				return fmt.Errorf("loading alert config: %w", err)
			}
			if !cfg.Enabled {
				fmt.Println("Alerts are disabled.")
				return nil
			}
			d := alert.BuildDispatcher(cfg)
			a := alert.Alert{
				Timestamp: time.Now(),
				Level:     alert.LevelFromString(cfg.Level),
				File:      "test",
				Message:   "confsnap alert test",
			}
			errs := d.Dispatch(a)
			if len(errs) > 0 {
				for _, e := range errs {
					fmt.Fprintf(os.Stderr, "alert error: %v\n", e)
				}
				return fmt.Errorf("%d handler(s) failed", len(errs))
			}
			return nil
		},
	}

	alertCmd.AddCommand(configureCmd, testCmd)
	rootCmd.AddCommand(alertCmd)
}
