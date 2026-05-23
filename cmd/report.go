package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/confsnap/internal/report"
)

var (
	reportFormat string
)

func init() {
	reportCmd := &cobra.Command{
		Use:   "report <label-a> <label-b>",
		Short: "Generate a diff report between two snapshots",
		Args:  cobra.ExactArgs(2),
		RunE:  runReport,
	}
	reportCmd.Flags().StringVarP(&reportFormat, "format", "f", "text", "Output format: text or json")
	rootCmd.AddCommand(reportCmd)
}

func runReport(cmd *cobra.Command, args []string) error {
	labelA, labelB := args[0], args[1]

	r, err := report.New(labelA, labelB)
	if err != nil {
		return fmt.Errorf("building report: %w", err)
	}

	var fmt report.Format
	switch reportFormat {
	case "json":
		fmt = report.FormatJSON
	default:
		fmt = report.FormatText
	}

	if err := r.Write(os.Stdout, fmt); err != nil {
		return fmt.Errorf("writing report: %w", err)
	}
	return nil
}
