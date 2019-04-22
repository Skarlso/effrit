package pkg

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
)

var yellow = color.New(color.FgYellow)
var red = color.New(color.FgRed)
var green = color.New(color.FgGreen, color.Bold)

// Display displays the analyies information in a pretty way...
// TODO: Add multiple display options and Graph generation.
func Display(packages map[string]Package) {
	const padding = 3
	table := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	writeColumns(table, []string{"NAME", "STABILITY"})
	for _, p := range packages {
		c := &color.Color{}
		if p.Stability < 0.5 {
			c = green
		} else if p.Stability >= 0.5 && p.Stability < 1 {
			c = yellow
		} else if p.Stability == 1 {
			c = red
		}
		stability := fmt.Sprintf("%.1f", p.Stability)
		writeColumns(table, []string{p.FullName, c.Sprint(stability)})
	}
	table.Flush()
}

func writeColumns(w io.Writer, cols []string) {
	fmt.Fprintln(w, strings.Join(cols, "\t"))
}
