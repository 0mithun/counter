package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	//filename := "words.txt"

	filename := "big.txt"
	//now := time.Now()
	file, err := os.Open(filename)

	if err != nil {
		log.Fatalln("failed to read file:", err)
	}
	defer file.Close()

	//scanner := bufio.NewScanner(file)
	//for scanner.Scan() {
	//	data := []byte(scanner.Bytes())
	//	wordCount := CountWords(data)
	//
	//	fmt.Println(wordCount)
	//}

	//data, err := io.ReadAll(file)

	wordCount := CountWords(file)

	fmt.Println(wordCount)
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
