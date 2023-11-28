package cmd

import (
	"bytes"
	"os"
	"testing"
)

// TestLSSuccess calls cmd.ls with an array of string arguments (path to file or directory)
// should print out a list of files
func TestLSSuccess(t *testing.T) {
	filepath := "../assets"

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	ls([]string{filepath})

	w.Close()
	os.Stdout = old

	// Read the captured output
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	// Check if the output contains the expected message
	expected := "dew.txt"
	if !bytes.Contains(buf.Bytes(), []byte(expected)) {
		t.Errorf("Expected message not found in stdout. Got: %s", buf.String())
	}
}

// TestLSFail calls cmd.ls with an array of string arguments (path to file or directory)
// should print out failed to read directory to the stderr
func TestLSFail(t *testing.T) {
	filepath := "assets"

	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	ls([]string{filepath})

	w.Close()
	os.Stderr = old

	// Read the captured output
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	// Check if the output contains the expected message
	expected := "failed to read directory"
	if !bytes.Contains(buf.Bytes(), []byte(expected)) {
		t.Errorf("Expected message not found in stderr. Got: %s", buf.String())
	}
}

// TestCATSuccess calls cmd.ls with an array of string arguments (path to file or directory)
// should print out a list of files
func TestCATSuccess(t *testing.T) {
	filepath := "../assets/for_you.txt"

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cat([]string{filepath})

	w.Close()
	os.Stdout = old

	// Read the captured output
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	// Check if the output contains the expected message
	expected := "Sonia Sanchez"
	if !bytes.Contains(buf.Bytes(), []byte(expected)) {
		t.Errorf("Expected message not found in stdout. Got: %s", buf.String())
	}
}

// TestCATFail calls cmd.cat with an array of string arguments (path to file)
// should print out failed to read file to the stderr
func TestCATFail(t *testing.T) {
	filepath := "assets/for_you.txt"

	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	ls([]string{filepath})

	w.Close()
	os.Stderr = old

	// Read the captured output
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	// Check if the output contains the expected message
	expected := "no such file"
	if !bytes.Contains(buf.Bytes(), []byte(expected)) {
		t.Errorf("Expected message not found in stderr. Got: %s", buf.String())
	}
}
