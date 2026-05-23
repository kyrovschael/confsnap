package report

import (
	"fmt"
	"io"

	"github.com/user/confsnap/internal/diff"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorGray   = "\033[90m"
)

func writeText(w io.Writer, r *Report) error {
	fmt.Fprintf(w, "confsnap diff: %s → %s\n", r.LabelA, r.LabelB)
	fmt.Fprintf(w, "Generated: %s\n", r.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "%s\n", r.Summary())
	fmt.Fprintln(w, strings.Repeat("─", 60))

	for _, e := range r.Entries {
		var color, prefix string
		switch e.Status {
		case diff.StatusAdded:
			color, prefix = colorGreen, "[+]"
		case diff.StatusRemoved:
			color, prefix = colorRed, "[-]"
		case diff.StatusChanged:
			color, prefix = colorYellow, "[~]"
		default:
			color, prefix = colorGray, "[=]"
		}
		fmt.Fprintf(w, "%s%s %s%s\n", color, prefix, e.Path, colorReset)
		if e.Status == diff.StatusChanged && e.Diff != "" {
			fmt.Fprintf(w, "    %s\n", e.Diff)
		}
	}
	return nil
}
