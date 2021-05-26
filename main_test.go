package main

import (
	"bytes"
	_ "embed"
	"io/ioutil"
	"os"
	"testing"
)

//go:embed golden.yaml
var goldenFile []byte

func Test_main(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	if !bytes.Equal(out, goldenFile) {
		t.Logf("output does not match golden.yaml content")
		t.Fail()
	}
}
