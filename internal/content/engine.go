package content

import "fmt"

func BuildPDFEngine(d *Document, extraction *Extraction) (*DetailedPDF, error) {
	if d == nil {
		return nil, fmt.Errorf("nil Document")
	}
	p := &DetailedPDF{document: d, structure: &d.Structure}
	p.Diagnostics = append(p.Diagnostics, d.Structure.Diagnostics...)
	if d.Encrypted {
		return p, fmt.Errorf("semantic decoding of encrypted PDF is unsupported")
	}
	b := &semanticBuilder{
		extract:     extraction,
		doc:         d,
		pdf:         p,
		fonts:       make(map[Span]int),
		images:      make(map[Span]int),
		activeForms: make(map[Span]bool),
		activePages: make(map[Span]bool),
		page:        -1,
		maxDepth:    128,
		maxObjects:  1000000,
	}
	b.cmaps = make(map[Span]*CMap)
	b.maxUnicodeBytes = 256 << 20
	if d.Options.Limits.MaxDecodedBytes > 0 && d.Options.Limits.MaxDecodedBytes < b.maxUnicodeBytes {
		b.maxUnicodeBytes = d.Options.Limits.MaxDecodedBytes
	}
	if d.Options.Limits.MaxDepth > 0 && d.Options.Limits.MaxDepth < b.maxDepth {
		b.maxDepth = d.Options.Limits.MaxDepth
	}
	if d.Options.Limits.MaxObjects > 0 && d.Options.Limits.MaxObjects < b.maxObjects {
		b.maxObjects = d.Options.Limits.MaxObjects
	}
	b.maxSemanticObjects = d.Options.Limits.MaxSemanticObjects
	if b.maxSemanticObjects == 0 {
		b.maxSemanticObjects = b.maxObjects
	}
	b.maxContentBytes = d.Options.Limits.MaxContentBytes
	if b.maxContentBytes == 0 {
		b.maxContentBytes = 256 << 20
	}
	b.maxValues = d.Options.Limits.MaxValues
	if b.maxValues == 0 {
		b.maxValues = 1 << 20
	}
	catalog, err := d.Catalog()
	if err != nil {
		return p, err
	}
	dict, err := semDictionary(catalog)
	if err != nil {
		return p, err
	}
	pages, err := dict.Get("Pages")
	if err != nil {
		return p, err
	}
	err = b.walkPages(pages, map[Name]Object{}, 0)
	return p, err
}

// ParseFile snapshots a file and prepares low-level object access. It closes
// the file before returning. At most one ReadOptions value may be supplied;
// a structural error can return a partial Document together with the error.
