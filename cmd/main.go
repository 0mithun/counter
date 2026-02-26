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
	args := display.NewOptionsArgs{}

	flag.BoolVar(&args.ShowWords, "w", false, "Used to toggle whether or not to show the word count")
	flag.BoolVar(&args.ShowBytes, "b", false, "Used to toggle whether or not to show the byte count")
	flag.BoolVar(&args.ShowLines, "l", false, "Used to toggle whether or not to show the line count")

	//flag.BoolVar(&display.ShowHeader, "header", false, "Used to toggle whether or not to show the header")

	flag.Parse()

	opts := display.NewOptions(args)

	wr := tabwriter.NewWriter(os.Stdout, 0, 8, 1, ' ', tabwriter.AlignRight)

	totals := counter.Counts{}
	filenames := flag.Args()
	didError := false
	totals.PrintHeader(wr, display.ShowHeader, opts)

	ch := counter.CountFiles(filenames)

	results := make([]counter.FileCountsResult, len(filenames))
	filenameIndex := make(map[string]int, len(filenames))

	for i, filename := range filenames {
		filenameIndex[filename] = i
	}

	for res := range ch {
		index, ok := filenameIndex[res.Filename]
		if !ok {
			continue
		}
		results[index] = res
	}

	for _, res := range results {
		if res.Err != nil {
			didError = true
			fmt.Fprintln(os.Stderr, "counter:", res.Err)
			continue
		}

		totals = totals.Add(res.Counts)
		res.Counts.Print(wr, opts, res.Filename)
	}

	if len(filenames) > 1 {
		totals.Print(wr, opts, "totals")
	}

	if len(filenames) == 0 {
		counter.GetCounts(os.Stdin).Print(wr, opts)
	}

	wr.Flush()

	if didError {
		os.Exit(1)
	}
}
