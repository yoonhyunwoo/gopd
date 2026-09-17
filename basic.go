package gopd

// PDF is the basic result returned by ParsePDF. The outer slice index is the
// zero-based page index; empty pages contain non-nil empty slices. Each inner
// slice follows content execution order for its kind, not reading order.
// Treat results and shared resources as read-only. ParsePDF computes only text
// and graphics; use Extract or Open for images, annotations, or provenance.
type PDF struct {
	Texts    [][]Text
	Graphics [][]Graphic
	// Diagnostics reports skipped or degraded constructs encountered while
	// producing the basic result.
	Diagnostics []ExtractionDiagnostic `json:",omitempty"`
}

// Text is one text-show operation. Matrix is in unrotated page user space,
// before font scaling. Glyph positions and source bytes are available through
// Extract with Glyphs or Provenance, and through Open.
type Text struct {
	Page             int
	Unicode          string
	Font             *FontInfo
	FontSize         float64
	Matrix           Matrix
	RenderingMode    int
	Style            PaintStyle
	DecodeComplete   bool
	PositionComplete bool
}

// Graphic is one path painting operation; it is not a whole chart or figure.
type Graphic struct {
	Page                  int
	Segments              []PathSegment
	Paint                 string
	Stroke, Fill, EvenOdd bool
	Style                 PaintStyle
}

// PathSegment points are already transformed into unrotated page user space.
// FontInfo contains display metadata, shared by original font resource identity.
// FontSize belongs to Text because the same font can be used at different sizes.

// basicPDF projects a text-and-graphics extraction into the basic result.
func basicPDF(ext *Extraction) *PDF {
	if ext == nil {
		return nil
	}
	p := &PDF{
		Texts:       make([][]Text, len(ext.Pages)),
		Graphics:    make([][]Graphic, len(ext.Pages)),
		Diagnostics: ext.Diagnostics,
	}
	for page := range ext.Pages {
		p.Texts[page] = []Text{}
		p.Graphics[page] = []Graphic{}
		pageIndex := ext.Pages[page].Index
		for _, text := range ext.Pages[page].Texts {
			t := Text{
				Page:           pageIndex,
				Unicode:        text.Unicode,
				Font:           text.Font,
				DecodeComplete: text.DecodeComplete,
			}
			if text.Position != nil {
				t.FontSize = text.Position.FontSize
				t.Matrix = text.Position.Matrix
				t.RenderingMode = text.Position.RenderingMode
				t.PositionComplete = text.Position.Complete
			}
			if text.Style != nil {
				t.Style = *text.Style
			}
			p.Texts[page] = append(p.Texts[page], t)
		}
		for _, graphic := range ext.Pages[page].Graphics {
			segments := make([]PathSegment, len(graphic.Segments))
			for j, segment := range graphic.Segments {
				segments[j] = PathSegment{Operator: segment.Operator, Points: append([]Point{}, segment.Points...)}
			}
			g := Graphic{
				Page:     pageIndex,
				Segments: segments,
				Paint:    graphic.Paint,
				Stroke:   graphic.Stroke,
				Fill:     graphic.Fill,
				EvenOdd:  graphic.EvenOdd,
			}
			if graphic.Style != nil {
				g.Style = *graphic.Style
			}
			p.Graphics[page] = append(p.Graphics[page], g)
		}
	}
	return p
}

// parsePDFOptions is the basic work contract: text and graphics with the
// position and style detail the basic result exposes, and nothing retained.
func parsePDFOptions() ExtractOptions {
	return ExtractOptions{Content: ContentText | ContentGraphics, Positions: true, Styles: true}
}
