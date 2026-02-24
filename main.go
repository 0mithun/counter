package main

import (
	"bufio"
	"fmt"
	"io"
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

func CountWordsInFile(filename string) (int, error) {
	file, err := os.Open(filename)

	if err != nil {
		return 0, err
	}
	defer file.Close()

	return CountWords(file), nil
}

func CountWordsInReader(file io.Reader) int {
	wordCount := 0

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		wordCount++
		//wordCount += CountWords(scanner.Bytes())
	}

	return wordCount
}

func CountWords(r io.Reader) int {
	wordCount := 0

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		wordCount++
	}

	return wordCount
}
