package gopd

import (
	"io"

	"github.com/MyungSub0519/gopd/internal/content"
	"github.com/MyungSub0519/gopd/internal/document"
	"github.com/MyungSub0519/gopd/internal/pdfmodel"
)

// ParsePDF snapshots a file and returns basic text and graphics grouped by
// page. Only text and graphics are interpreted; images, annotations, and
// provenance require Extract or Open. A semantic failure can return both a
// partial PDF and an error. Callers must check err before treating the result
// as successful. No Close call is required.
func ParsePDF(path string) (*PDF, error) {
	ext, err := content.Extract(path, parsePDFOptions())
	return basicPDF(ext), err
}

// Open returns detailed contents, resources, operations, and byte provenance.
// Use ParsePDF for a basic result or Extract for a selected subset. Open
// closes the input file.
func Open(path string) (*DetailedPDF, error) {
	d, err := document.ParseFile(path)
	if err != nil {
		return nil, err
	}
	return content.BuildPDFEngine(d, nil)
}

// Read parses a bounded snapshot of r and returns detailed page contents.
// The caller retains ownership of r; no Close call is required on the result.
func Read(r io.ReaderAt, size int64) (*DetailedPDF, error) {
	d, err := document.Parse(r, size)
	if err != nil {
		return nil, err
	}
	return content.BuildPDFEngine(d, nil)
}

// Int converts a direct PDF Integer into int64, rejecting other types and
// out-of-range values. Resolve indirect references before calling Int.
func Int(object Object) (int64, error) { return pdfmodel.Int(object) }

// Number converts a direct Integer or Real into a finite float64. It does not
// resolve indirect references or preserve exact decimal precision.
func Number(object Object) (float64, error) { return pdfmodel.Number(object) }

// IsStream reports whether object directly contains a Stream. It does not
// follow references; detailed and extraction results expose resolved values.
func IsStream(object Object) bool { return pdfmodel.IsStream(object) }

// IdentityMatrix returns a transformation that leaves coordinates unchanged.
func IdentityMatrix() Matrix { return pdfmodel.IdentityMatrix() }
