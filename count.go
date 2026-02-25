package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
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

func (c Counts) Print(w io.Writer, opts DisplayOptions, suffix ...string) {
	//fmt.Fprintln(w, c.Lines, c.Words, c.Bytes, filename)
	xs := []string{}

	if opts.ShouldShowLines() {
		xs = append(xs, strconv.Itoa(c.Lines))
	}

	if opts.ShouldShowWords() {
		xs = append(xs, strconv.Itoa(c.Words))
	}

	if opts.ShouldShowBytes() {
		xs = append(xs, strconv.Itoa(c.Bytes))
	}

	xs = append(xs, suffix...)

	line := strings.Join(xs, " ")
	fmt.Fprintln(w, line)

	//fmt.Fprintf(w, "%d %d %d", c.Lines, c.Words, c.Bytes)
	//
	//for _, filename := range suffix {
	//	fmt.Fprintf(w, " %s", filename)
	//}

}

func GetCounts(f io.Reader) Counts {
	//const offsetStart = 0
	//lines := CountLines(f)
	//f.Seek(offsetStart, io.SeekStart)
	//
	//words := CountWords(f)
	//f.Seek(offsetStart, io.SeekStart)
	//
	//bytes := CountBytes(f)
	//f.Seek(offsetStart, io.SeekStart)
	//
	//return Counts{
	//	Lines: lines,
	//	Words: words,
	//	Bytes: bytes,
	//}

	res := Counts{}

	isInsideWord := false

	reader := bufio.NewReader(f)

	for {
		r, size, err := reader.ReadRune()
		if err != nil {
			break
		}
		res.Bytes += size

		if r == '\n' {
			res.Lines++
		}

		isSpace := unicode.IsSpace(r)

		if !isSpace && !isInsideWord {
			res.Words++
		}

		isInsideWord = !isSpace
	}

	return res
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
