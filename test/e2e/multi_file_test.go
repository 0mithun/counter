package e2e

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/0mithun/counter/test/assert"
)

func TestMultipleFiles(t *testing.T) {
	fileA, err := createFile("one two three four five\n")

	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(fileA.Name())

	fileB, err := createFile("foo bar baz\n\n")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(fileB.Name())

	fileC, err := createFile("")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(fileC.Name())

	cmd, err := getCommand(fileA.Name(), fileB.Name(), fileC.Name())
	if err != nil {
		t.Fatal("failed to create command", err)
	}

	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout

	err = cmd.Run()
	if err != nil {
		t.Fatal("failed to execute command", err)
	}

	//wants := map[string]string{
	//	fileA.Name(): fmt.Sprintf(" 1 5 24 %s", fileA.Name()),
	//	fileB.Name(): fmt.Sprintf(" 2 3 13 %s", fileB.Name()),
	//	fileC.Name(): fmt.Sprintf(" 0 0  0 %s", fileC.Name()),
	//	"totals":     " 3 8 37 totals",
	//}

	wants := fmt.Sprintf(` 1 5 24 %s
 2 3 13 %s
 0 0  0 %s
 3 8 37 totals
`, fileA.Name(), fileB.Name(), fileC.Name())

	assert.Equal(t, wants, stdout.String())
}
