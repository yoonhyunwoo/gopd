package document

import (
	"bytes"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/MyungSub0519/gopd/internal/pdfmodel"
	"github.com/MyungSub0519/gopd/internal/pdftest"
)

const contractContent = "0 0 m 10 10 l S BT /F 12 Tf 1 0 0 1 20 30 Tm (A) Tj ET"

// The former public low-level contract, now enforced at the package boundary.
func TestDocumentObjectsAndStreamsContract(t *testing.T) {
	encoded := hex.EncodeToString([]byte(contractContent)) + ">"
	objects := []string{
		`<< /Type /Catalog /Pages 2 0 R >>`,
		`<< /Type /Pages /Kids [3 0 R] /Count 1 /MediaBox [0 0 100 100] >>`,
		`<< /Type /Page /Parent 2 0 R /Contents 4 0 R /Resources << /Font << /F 5 0 R >> >> >>`,
		pdftest.Stream("/Filter /ASCIIHexDecode", encoded),
		`<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding /WinAnsiEncoding /FirstChar 65 /Widths [600] >>`,
	}
	data := pdftest.File(objects, "")
	path := filepath.Join(t.TempDir(), "contract.pdf")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	options := ReadOptions{MaxFileBytes: int64(len(data))}
	for _, parse := range []func() (*Document, error){
		func() (*Document, error) { return ParseFile(path, options) },
		func() (*Document, error) { return Parse(bytes.NewReader(data), int64(len(data)), options) },
	} {
		doc, err := parse()
		if err != nil {
			t.Fatal(err)
		}
		catalog, err := doc.Catalog()
		if err != nil {
			t.Fatal(err)
		}
		typeObject, err := catalog.Value.(pdfmodel.Dictionary).Get("Type")
		if err != nil || typeObject.Value != pdfmodel.Name("Catalog") {
			t.Fatalf("catalog lookup = %+v, %v", typeObject, err)
		}
		id := pdfmodel.ObjectID{Number: 4, Generation: 0}
		object, err := doc.Load(id)
		if err != nil || !pdfmodel.IsStream(object.Body) {
			t.Fatalf("stream lookup = %+v, %v", object, err)
		}
		ref := pdfmodel.Reference{ID: id}
		resolved, err := doc.Resolve(ref)
		if err != nil || resolved.Span != object.Body.Span {
			t.Fatalf("reference resolution = %+v, %v", resolved, err)
		}
		resolved, err = doc.ResolveObject(pdfmodel.Object{Value: ref})
		if err != nil || resolved.Span != object.Body.Span {
			t.Fatalf("object resolution = %+v, %v", resolved, err)
		}
		raw, err := doc.RawObject(resolved)
		if err != nil || !bytes.Equal(raw, data[resolved.Span.Start:resolved.Span.End]) {
			t.Fatal("raw source bytes changed")
		}
		source, err := doc.DecodeStream(resolved.Value.(pdfmodel.Stream))
		if err != nil {
			t.Fatal(err)
		}
		decoded, err := doc.Bytes(pdfmodel.Span{Source: source.ID, Start: 0, End: source.Size})
		if err != nil || string(decoded) != contractContent {
			t.Fatalf("decoded source = %q, %v", decoded, err)
		}
	}
}
