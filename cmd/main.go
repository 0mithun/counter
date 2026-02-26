package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"text/tabwriter"

	"github.com/0mithun/counter"
	"github.com/0mithun/counter/display"
)

type FileCountsResult struct {
	counts   counter.Counts
	filename string
}

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

	wg := sync.WaitGroup{}
	ch := make(chan FileCountsResult)

	wg.Add(len(filenames))

	for _, filename := range filenames {
		go func() {
			defer wg.Done()

			counts, err := counter.CountFile(filename)
			if err != nil {
				didError = true
				fmt.Fprintln(os.Stderr, "counter:", err)
				return
			}
			ch <- FileCountsResult{
				counts:   counts,
				filename: filename,
			}c

		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for res := range ch {
		totals = totals.Add(res.counts)
		res.counts.Print(wr, opts, res.filename)
	}

	if len(filenames) > 1 {
		totals.Print(wr, opts, "total")
	}

	if len(filenames) == 0 {
		counter.GetCounts(os.Stdin).Print(wr, opts)
	}

	wr.Flush()

	if didError {
		os.Exit(1)
	}
}
