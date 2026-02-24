package main_test

import (
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
