package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/nicholasgasior/confsnap/internal/baseline"
	"github.com/nicholasgasior/confsnap/internal/export"
	"github.com/spf13/cobra"
)

func init() {
	var baselineLabel string
	var format string
	var output string

	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "Export drift results to a file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExport(baselineLabel, format, output)
		},
	}

	exportCmd.Flags().StringVarP(&baselineLabel, "baseline", "b", "", "baseline label to compare against (required)")
	exportCmd.Flags().StringVarP(&format, "format", "f", "csv", "output format: csv, markdown, html, pdf")
	exportCmd.Flags().StringVarP(&output, "output", "o", "", "output file path (stdout if omitted)")
	_ = exportCmd.MarkFlagRequired("baseline")

	rootCmd.AddCommand(exportCmd)
}

func runExport(baselineLabel, format, outputPath string) error {
	bl, err := baseline.Load(baselineLabel)
	if err != nil {
		return fmt.Errorf("load baseline: %w", err)
	}

	results, err := baseline.CheckDrift(bl)
	if err != nil {
		return fmt.Errorf("check drift: %w", err)
	}

	w := os.Stdout
	if outputPath != "" {
		f, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("create output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	now := time.Now()
	switch format {
	case "csv":
		return export.DriftToCSV(w, results, now, baselineLabel)
	case "markdown":
		return export.DriftToMarkdown(w, results, now, baselineLabel)
	case "html":
		return export.DriftToHTML(w, results, now, baselineLabel)
	case "pdf":
		return export.DriftToPDF(w, results, now, baselineLabel)
	default:
		return fmt.Errorf("unknown format %q: choose csv, markdown, html, or pdf", format)
	}
}
