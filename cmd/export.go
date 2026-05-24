package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/confsnap/internal/baseline"
	"github.com/yourusername/confsnap/internal/export"
)

func init() {
	var format string
	var baselineLabel string
	var files []string
	var outputFile string

	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "Export a drift report in various formats (csv, markdown, html)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if baselineLabel == "" {
				return fmt.Errorf("--baseline is required")
			}
			if len(files) == 0 {
				return fmt.Errorf("at least one --file is required")
			}

			bl, err := baseline.Load(baselineLabel)
			if err != nil {
				return fmt.Errorf("loading baseline %q: %w", baselineLabel, err)
			}

			results, err := baseline.CheckDrift(files, bl)
			if err != nil {
				return fmt.Errorf("checking drift: %w", err)
			}

			out := os.Stdout
			if outputFile != "" {
				f, err := os.Create(outputFile)
				if err != nil {
					return fmt.Errorf("creating output file: %w", err)
				}
				defer f.Close()
				out = f
			}

			now := time.Now()
			switch format {
			case "csv":
				return export.DriftToCSV(out, results, baselineLabel, now)
			case "markdown":
				return export.DriftToMarkdown(out, results, baselineLabel, now)
			case "html":
				return export.DriftToHTML(out, results, baselineLabel, now)
			default:
				return fmt.Errorf("unknown format %q: choose csv, markdown, or html", format)
			}
		},
	}

	exportCmd.Flags().StringVarP(&format, "format", "f", "csv", "Output format: csv, markdown, html")
	exportCmd.Flags().StringVarP(&baselineLabel, "baseline", "b", "", "Baseline label to compare against")
	exportCmd.Flags().StringSliceVarP(&files, "file", "p", nil, "Files to check (repeatable)")
	exportCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write output to file instead of stdout")

	rootCmd.AddCommand(exportCmd)
}
