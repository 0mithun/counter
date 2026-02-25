package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/0mithun/counter"
	"github.com/0mithun/counter/display"
)

func main() {
	opts := display.Options{}

	flag.BoolVar(&opts.ShowWords, "w", false, "Used to toggle whether or not to show the word count")
	flag.BoolVar(&opts.ShowBytes, "b", false, "Used to toggle whether or not to show the byte count")
	flag.BoolVar(&opts.ShowLines, "l", false, "Used to toggle whether or not to show the line count")

	flag.BoolVar(&display.ShowHeader, "header", false, "Used to toggle whether or not to show the header")

	flag.Parse()

	wr := tabwriter.NewWriter(os.Stdout, 0, 8, 1, ' ', tabwriter.AlignRight)

	totals := counter.Counts{}
	filenames := flag.Args()
	didError := false
	totals.PrintHeader(wr, display.ShowHeader, opts)
	for _, filename := range filenames {
		counts, err := counter.CountFile(filename)
		if err != nil {
			didError = true
			fmt.Fprintln(os.Stderr, "counter:", err)
			continue
		}
		totals = totals.Add(counts)

		counts.Print(wr, opts, filename)
	}

	if len(filenames) == 0 {
		counter.GetCounts(os.Stdin).Print(wr, opts)
	}

	if len(filenames) > 1 {
		totals.Print(wr, opts, "total")
	}

	wr.Flush()

	if didError {
		os.Exit(1)
	}
}
