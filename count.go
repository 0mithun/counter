package main

import (
	"bufio"
	"io"
	"os"
)

func CountWordsInFile(filename string) (int, error) {
	file, err := os.Open(filename)

	if err != nil {
		return 0, err
	}
	defer file.Close()

	return CountWords(file), nil
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

func CountLines(r io.Reader) int {
	//lineCount := 0
	//scanner := bufio.NewScanner(r)
	//scanner.Split(bufio.ScanLines)
	//for scanner.Scan() {
	//	lineCount++
	//}
	//return lineCount

	linesCount := 0
	reader := bufio.NewReader(r)
	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if r == '\n' {
			linesCount++
		}
	}

	return linesCount
}
