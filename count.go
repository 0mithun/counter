package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Counts struct {
	Bytes int
	Words int
	Lines int
}

// Add will modify the values of the count by
// adding the values from the other.
func (c Counts) Add(other Counts) Counts {
	c.Bytes += other.Bytes
	c.Words += other.Words
	c.Lines += other.Lines

	return c
}

func (c Counts) Print(w io.Writer, filenames ...string) {
	//fmt.Fprintln(w, c.Lines, c.Words, c.Bytes, filename)
	fmt.Fprintf(w, "%d %d %d", c.Lines, c.Words, c.Bytes)

	for _, filename := range filenames {
		fmt.Fprintf(w, " %s", filename)
	}

	fmt.Fprintf(w, "\n")
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
