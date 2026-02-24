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

	total := 0
	filenames := os.Args[1:]
	didError := false
	for _, filename := range filenames {
		wordCount, err := CountWordsInFile(filename)
		if err != nil {
			didError = true
			fmt.Fprintln(os.Stderr, "counter:", err)
			continue
		}
		fmt.Println(wordCount, filename)
		total += wordCount
	}

	if len(filenames) == 0 {
		total = CountWords(os.Stdin)
		fmt.Println(total)
		return
	}

	if len(filenames) > 1 {
		fmt.Println(total, "total")
	}

	fmt.Println(total)
	if didError {
		os.Exit(1)
	}
}
