package main

import (
	"bufio"
	"io"
	"os"
)

type Counts struct {
	Bytes int
	Words int
	Lines int
}

func GetCounts(f io.ReadSeeker) Counts {
	const offsetStart = 0
	lines := CountLines(f)
	f.Seek(offsetStart, io.SeekStart)

	words := CountWords(f)
	f.Seek(offsetStart, io.SeekStart)

	bytes := CountBytes(f)
	f.Seek(offsetStart, io.SeekStart)

	return Counts{
		Lines: lines,
		Words: words,
		Bytes: bytes,
	}
}

func CountFile(filename string) (Counts, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Counts{}, err
	}
	defer file.Close()

	counts := GetCounts(file)

	return counts, err
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

func CountBytes(r io.Reader) int {
	bytesCount, _ := io.Copy(io.Discard, r)

	return int(bytesCount)
}
