package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	data, _ := os.ReadFile("./words.txt")

	wordCount := CountWords(data)

	fmt.Println(wordCount)
}

func CountWords(data []byte) int {
	//wordCount := 0
	//wasSpace := true
	//for _, x := range data {
	//	isSpace := x == ' ' || x == '\n'
	//	if wasSpace && !isSpace {
	//		wordCount++
	//	}
	//
	//	wasSpace = isSpace
	//}
	//
	//return wordCount

	//words := strings.Fields(string(data))

	words := bytes.Fields(data)

	return len(words)
}
