package main_test

import (
	"bytes"
	"strings"
	"testing"

	counter "github.com/0mithun/counter"
)

func TestCountWords(t *testing.T) {
	type testCase struct {
		name  string
		input string
		wants int
	}

	testCases := []testCase{
		{name: "count five", input: "one two three four five", wants: 5},
		{name: "empty string", input: "", wants: 0},
		{name: "empty space", input: " ", wants: 0},
		{name: "prefix space", input: " one", wants: 1},
		{name: "suffix space", input: "one ", wants: 1},
		{name: "both side space", input: " one ", wants: 1},
		{name: "multi word left space", input: " one two", wants: 2},
		{name: "new lines", input: "one two three\nfour five", wants: 5},
		{name: "multiple spaces", input: "one two  three four five", wants: 5},
		{name: "tab character", input: "one two\tthree four five", wants: 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.input)

			results := counter.CountWords(r)
			if results != tc.wants {
				t.Errorf("countWords(%q) = %d, want %d", tc.input, results, tc.wants)
			}
		})
	}
}

func TestCountLines(t *testing.T) {
	type testCase struct {
		name  string
		input string
		wants int
	}

	testCases := []testCase{
		{
			name:  "Simple five words, 1 new lines",
			input: "one two three four five\n",
			wants: 1,
		},
		{
			name:  "empty file",
			input: "",
			wants: 0,
		},
		{
			name:  "no new lines",
			input: "one two three four five",
			wants: 0,
		},
		{
			name:  "no new lines at end",
			input: "one two three four five\nsix",
			wants: 1,
		},
		{
			name:  "multi newline string",
			input: "\n\n\n\n",
			wants: 4,
		},
		{
			name:  "multi word line string",
			input: "one\ntwo\nthree\nfour\nfive\n",
			wants: 5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.input)
			results := counter.CountLines(r)
			if results != tc.wants {
				t.Errorf("countLines(%q) = %d, want %d", tc.input, results, tc.wants)
			}
		})
	}
}

func TestCountBytes(t *testing.T) {
	type testCase struct {
		name  string
		input string
		wants int
	}

	testCases := []testCase{
		{name: "empty string", input: "", wants: 0},
		{name: "five words", input: "one two three four five", wants: 23},
		{name: "all spaces", input: "     ", wants: 5},
		{name: "newlines and words", input: "one\ntwo\nthree\nfour\t\n", wants: 20},
		{name: "unicode characters", input: "θβ", wants: 4},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.input)
			results := counter.CountBytes(r)
			if results != tc.wants {
				t.Errorf("countBytes(%q) = %d, want %d", tc.input, results, tc.wants)
			}
		})
	}
}

func TestGetCount(t *testing.T) {
	type testCase struct {
		name  string
		input string
		wants counter.Counts
	}

	testCases := []testCase{
		{
			name:  "simple five words",
			input: "one two three four five\n",
			wants: counter.Counts{Lines: 1, Words: 5, Bytes: 24},
		},
		{
			name:  "five words no new line",
			input: "one two three four five",
			wants: counter.Counts{Lines: 0, Words: 5, Bytes: 23},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.input)
			results := counter.GetCounts(r)

			if results != tc.wants {
				t.Errorf("GetCounts(%q) = %d, want %d", tc.input, results, tc.wants)
			}
		})
	}
}

func TestPrintCounts(t *testing.T) {
	type inputs struct {
		counts   counter.Counts
		opts     counter.DisplayOptions
		filename []string
	}

	type testCase struct {
		name  string
		input inputs
		wants string
	}

	testCases := []testCase{
		{
			name: "simple five words.txt",
			input: inputs{
				counts:   counter.Counts{Lines: 1, Words: 5, Bytes: 24},
				filename: []string{"words.txt"},
				opts: counter.DisplayOptions{
					ShowLines: true,
					ShowWords: true,
					ShowBytes: true,
				},
			},
			wants: "1 5 24 words.txt\n",
		},
		{
			name: "simple five words.txt show lines",
			input: inputs{
				counts:   counter.Counts{Lines: 1, Words: 5, Bytes: 24},
				filename: []string{"words.txt"},
				opts: counter.DisplayOptions{
					ShowLines: true,
					ShowWords: false,
					ShowBytes: false,
				},
			},
			wants: "1 words.txt\n",
		},
		{
			name: "simple five words.txt show words",
			input: inputs{
				counts:   counter.Counts{Lines: 1, Words: 5, Bytes: 24},
				filename: []string{"words.txt"},
				opts: counter.DisplayOptions{
					ShowLines: false,
					ShowWords: true,
					ShowBytes: false,
				},
			},
			wants: "5 words.txt\n",
		},
		{
			name: "simple five words.txt show bytes",
			input: inputs{
				counts:   counter.Counts{Lines: 1, Words: 5, Bytes: 24},
				filename: []string{"words.txt"},
				opts: counter.DisplayOptions{
					ShowLines: false,
					ShowWords: false,
					ShowBytes: true,
				},
			},
			wants: "24 words.txt\n",
		},
		{
			name: "simple five words.txt show bytes and lines",
			input: inputs{
				counts:   counter.Counts{Lines: 1, Words: 5, Bytes: 24},
				filename: []string{"words.txt"},
				opts: counter.DisplayOptions{
					ShowLines: true,
					ShowWords: false,
					ShowBytes: true,
				},
			},
			wants: "1 24 words.txt\n",
		},
		{
			name: "simple five words.txt no options",
			input: inputs{
				counts:   counter.Counts{Lines: 1, Words: 5, Bytes: 24},
				filename: []string{"words.txt"},
			},
			wants: "1 5 24 words.txt\n",
		},
		{
			name: "empty filename",
			input: inputs{
				counts: counter.Counts{Lines: 1, Words: 4, Bytes: 20},
				opts: counter.DisplayOptions{
					ShowLines: true,
					ShowWords: true,
					ShowBytes: true,
				},
			},
			wants: "1 4 20\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buffer := bytes.Buffer{}
			tc.input.counts.Print(&buffer, tc.input.opts, tc.input.filename...)
			if got := buffer.String(); got != tc.wants {
				t.Errorf("PrintCounts(%q) = %q, want %q", tc.input.filename, got, tc.wants)
			}
		})
	}
}

func TestAddCounts(t *testing.T) {
	type inputs struct {
		counts counter.Counts
		other  counter.Counts
	}
	type testCase struct {
		name  string
		input inputs
		wants counter.Counts
	}

	testCases := []testCase{
		{
			name: "simple add by one",
			input: inputs{
				counts: counter.Counts{Lines: 1, Words: 5, Bytes: 24},
				other:  counter.Counts{Lines: 1, Words: 1, Bytes: 1},
			},
			wants: counter.Counts{Lines: 2, Words: 6, Bytes: 25},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			totals := tc.input.counts
			results := totals.Add(tc.input.other)
			if results != tc.wants {
				t.Errorf("AddCounts(%q) = %d, want %d", tc.input.other, results, tc.wants)
			}
		})
	}
}
