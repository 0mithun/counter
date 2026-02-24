package main

import (
	"fmt"
	"os"
)

func main() {
	//filename := "words.txt"
	//if len(os.Args) < 2 {
	//	log.Fatalln("error: no filename provided")
	//}
	//filename := os.Args[1]

	totals := Counts{}
	filenames := os.Args[1:]
	didError := false
	for _, filename := range filenames {
		counts, err := CountFile(filename)
		if err != nil {
			didError = true
			fmt.Fprintln(os.Stderr, "counter:", err)
			continue
		}
		totals = Counts{
			Lines: totals.Lines + counts.Lines,
			Words: totals.Words + counts.Words,
			Bytes: totals.Bytes + counts.Bytes,
		}

		counts.Print(os.Stdout, filename)
	}

	if len(filenames) == 0 {
		GetCounts(os.Stdin).Print(os.Stdout, "")
	}

	if len(filenames) > 1 {
		totals.Print(os.Stdout, "totals")
	}

	if didError {
		os.Exit(1)
	}
}
