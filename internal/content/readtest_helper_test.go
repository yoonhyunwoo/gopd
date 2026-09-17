package content

import (
	"bytes"
	"io"
	"testing"
)

// readAll is the detailed two-step entry used by tests: parse then interpret.
func readAll(t *testing.T, data []byte) (*DetailedPDF, error) {
	t.Helper()
	d, err := Parse(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}
	return BuildPDFEngine(d, nil)
}

func readAllErr(r io.ReaderAt, size int64) (*DetailedPDF, error) {
	d, err := Parse(r, size)
	if err != nil {
		return nil, err
	}
	return BuildPDFEngine(d, nil)
}
