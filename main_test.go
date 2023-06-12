package main

import (
	"bytes"
	_ "embed"
	"io"
	"os"
	"testing"
)

//go:embed testdata/golden.yaml
var goldenFile []byte

func Test_main(t *testing.T) {
	rescueStdout := os.Stdout
	rescueArgs := os.Args

	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{rescueArgs[0], "testdata/test-scripts"}

	main()

	w.Close()
	out, _ := io.ReadAll(r)

	os.Stdout = rescueStdout
	os.Args = rescueArgs

	if !bytes.Equal(out, goldenFile) {
		t.Logf("output does not match golden.yaml content")
		t.Fail()
	}
}
