package counter

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/0mithun/counter/display"
)

type Counts struct {
	bytes int
	words int
	lines int
}

type FileCountsResult struct {
	Counts   Counts
	Filename string
	Err      error
}

// Add will modify the values of the count by
// adding the values from the other.
func (c Counts) Add(other Counts) Counts {
	c.bytes += other.bytes
	c.words += other.words
	c.lines += other.lines

	return c
}

func (c Counts) Print(w io.Writer, opts display.Options, suffix ...string) {
	//fmt.Fprintln(w, c.lines, c.words, c.bytes, filename)
	stats := []string{}

	if opts.ShouldShowLines() {
		stats = append(stats, strconv.Itoa(c.lines))
	}

	if opts.ShouldShowWords() {
		stats = append(stats, strconv.Itoa(c.words))
	}

	if opts.ShouldShowBytes() {
		stats = append(stats, strconv.Itoa(c.bytes))
	}

	line := strings.Join(stats, "\t") + "\t"
	fmt.Fprint(w, line)

	suffixStr := strings.Join(suffix, " ")
	if suffixStr != "" {
		fmt.Fprintf(w, " %s", suffixStr)
	}
	fmt.Fprint(w, "\n")
}

func (c Counts) PrintHeader(w io.Writer, showHeader bool, opts display.Options) {
	if !showHeader {
		return
	}
	stats := []string{}

	if opts.ShouldShowLines() {
		stats = append(stats, "lines")
	}

	if opts.ShouldShowWords() {
		stats = append(stats, "words")
	}

	if opts.ShouldShowBytes() {
		stats = append(stats, "bytes")
	}

	line := strings.Join(stats, "\t") + "\t"
	fmt.Fprintln(w, line)
}

func CountFiles(filenames []string) <-chan FileCountsResult {
	ch := make(chan FileCountsResult)

	wg := sync.WaitGroup{}
	wg.Add(len(filenames))

	for _, filename := range filenames {
		go func() {
			defer wg.Done()
			res, err := CountFile(filename)

			ch <- FileCountsResult{
				Counts:   res,
				Filename: filename,
				Err:      err,
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

func GetCounts(r io.Reader) Counts {
	bytesReader, bytesWriter := io.Pipe()
	wordsReader, wordsWriter := io.Pipe()
	linesReader, linesWriter := io.Pipe()

	w := io.MultiWriter(bytesWriter, wordsWriter, linesWriter)

	chBytes := make(chan int)
	chWords := make(chan int)
	chLines := make(chan int)

	go func() {
		defer close(chBytes)
		chBytes <- CountBytes(bytesReader)
	}()

	go func() {
		defer close(chWords)
		chWords <- CountWords(wordsReader)
	}()

	go func() {
		defer close(chLines)
		chLines <- CountLines(linesReader)
	}()

	io.Copy(w, r)
	bytesWriter.Close()
	wordsWriter.Close()
	linesWriter.Close()

	bytesCount := <-chBytes
	wordsCount := <-chWords
	linesCount := <-chLines

	return Counts{
		bytes: bytesCount,
		words: wordsCount,
		lines: linesCount,
	}
}

func GetCountsPipe(r io.Reader) Counts {

	p1r, p1w := io.Pipe()
	p2r, p2w := io.Pipe()

	bytesReader := io.TeeReader(r, p1w)
	wordsReader := io.TeeReader(p1r, p2w)
	linesReader := p2r

	chBytes := make(chan int)
	chWords := make(chan int)
	chLines := make(chan int)

	go func() {
		defer p1w.Close()
		defer close(chBytes)
		chBytes <- CountBytes(bytesReader)

	}()

	go func() {
		defer p2w.Close()
		defer close(chWords)
		chWords <- CountWords(wordsReader)
	}()

	go func() {
		defer close(chLines)
		chLines <- CountLines(linesReader)
	}()

	bytesCount := <-chBytes
	wordsCount := <-chWords
	linesCount := <-chLines

	return Counts{
		bytes: bytesCount,
		words: wordsCount,
		lines: linesCount,
	}
}

func GetCountsTeeReader(r io.Reader) Counts {
	byteBuf := &bytes.Buffer{}
	wordBuf := &bytes.Buffer{}

	bytesReader := io.TeeReader(r, byteBuf)
	wordsReader := io.TeeReader(byteBuf, wordBuf)
	lineReader := wordBuf

	byteCount := CountBytes(bytesReader)
	wordCount := CountWords(wordsReader)
	lineCount := CountLines(lineReader)

	return Counts{
		bytes: byteCount,
		words: wordCount,
		lines: lineCount,
	}
}

func GetCountsSinglePass(f io.Reader) Counts {
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
	//	lines: lines,
	//	words: words,
	//	bytes: bytes,
	//}

	res := Counts{}

	isInsideWord := false

	reader := bufio.NewReader(f)

	for {
		r, size, err := reader.ReadRune()
		if err != nil {
			break
		}
		res.bytes += size

		if r == '\n' {
			res.lines++
		}

		isSpace := unicode.IsSpace(r)

		if !isSpace && !isInsideWord {
			res.words++
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
