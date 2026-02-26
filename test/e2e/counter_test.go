package e2e

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/0mithun/counter/test/assert"
)

func TestStdin(t *testing.T) {
	cmd, err := getCommand()
	if err != nil {
		t.Fatal("could not get current working directory")
	}

	output := &bytes.Buffer{}

	cmd.Stdin = strings.NewReader("one two three\n")
	cmd.Stdout = output

	if err := cmd.Run(); err != nil {
		t.Fatal("failed to run command:", err)
	}

	wants := " 1 3 14\n"

	assert.Equal(t, wants, output.String())
}

func TestSingleFile(t *testing.T) {
	file, err := os.CreateTemp("", "counter-test-*")
	if err != nil {
		t.Fatal("could not create temporary file:", err)
	}
	defer os.Remove(file.Name())

	_, err = file.WriteString("foo bar baz\nbaz bar foo\none two three\n")
	if err != nil {
		t.Fatal("could not write to temporary file:", err)
	}
	defer file.Close()

	cmd, err := getCommand(file.Name())
	if err != nil {
		t.Fatal("could not get current working directory")
	}

	output := &bytes.Buffer{}
	cmd.Stdout = output

	err = cmd.Run()
	if err != nil {
		t.Fatal("failed to run command:", err)
	}

	//wants := fmt.Sprintln(" 3 9 38", file.Name())
	wants := fmt.Sprintf(" 3 9 38 %s\n", file.Name())

	assert.Equal(t, wants, output.String())
}

func TestNotExitsFile(t *testing.T) {
	filename := "noexist.txt"
	cmd, err := getCommand(filename)
	if err != nil {
		t.Fatal("could not get current working directory")
	}

	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err = cmd.Run()
	if err == nil {
		t.Log("command succeeded when should have failed")
		t.Fail()
	}

	if err.Error() != "exit status 1" {
		t.Logf("got %q, want \"exit status 1\"", err)
		t.Fail()
	}

	wantsStderr := fmt.Sprintf("counter: open %s: no such file or directory\n", filename)
	wantStdOutput := ""

	assert.Equal(t, wantsStderr, stderr.String())
	assert.Equal(t, wantStdOutput, stdout.String())
}

func TestFlags(t *testing.T) {
	file, err := createFile("one two three four five\none two three\n")
	if err != nil {
		t.Fatal("could not create temporary file:", err)
	}
	defer os.Remove(file.Name())

	testCases := []struct {
		name  string
		want  string
		flags []string
	}{
		{
			name:  "line flag",
			want:  fmt.Sprintf(" 2 %s\n", file.Name()),
			flags: []string{"-l"},
		},
		{
			name:  "word flag",
			want:  fmt.Sprintf(" 8 %s\n", file.Name()),
			flags: []string{"-w"},
		},
		{
			name:  "byte flag",
			want:  fmt.Sprintf(" 38 %s\n", file.Name()),
			flags: []string{"-b"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputs := append(tc.flags, file.Name())
			cmd, err := getCommand(inputs...)
			if err != nil {
				t.Log("failed to get command:", err)
				t.Fail()
			}

			stdout := &bytes.Buffer{}
			cmd.Stdout = stdout

			err = cmd.Run()
			if err != nil {
				t.Log("failed to run command:", err)
				t.Fail()
			}

			wants := tc.want
			if stdout.String() != wants {
				t.Errorf("got %q, want %q", stdout.String(), wants)
			}

			assert.Equal(t, wants, stdout.String())
		})
	}

}
