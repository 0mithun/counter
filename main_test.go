package main_test

import (
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
			results := counter.CountWords([]byte(tc.input))
			if results != tc.wants {
				t.Errorf("countWords(%q) = %d, want %d", tc.input, results, tc.wants)
			}
		})
	}

}
