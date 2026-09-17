package content

import (
	"fmt"
	"io"
)

// =====================================================================
// Entry points: selective extraction API and flow control.
// =====================================================================

// Extract snapshots a file and directly emits selected content. The file is
// closed before return. A semantic error can accompany a partial Extraction;
// callers must check err even when the result is non-nil.
func Extract(path string, options ExtractOptions) (*Extraction, error) {
	options, err := normalizeExtractOptions(options)
	if err != nil {
		return nil, err
	}
	doc, err := ParseFile(path, options.ReadOptions)
	if err != nil {
		return nil, err
	}
	return extractDocument(doc, options)
}

// ExtractReader takes one bounded input snapshot and never closes r. Content
// kinds share one interpreter; reused Forms still execute under each caller's
// state. This API does not promise bounded total process memory or streaming I/O.
func ExtractReader(r io.ReaderAt, size int64, options ExtractOptions) (*Extraction, error) {
	options, err := normalizeExtractOptions(options)
	if err != nil {
		return nil, err
	}
	doc, err := Parse(r, size, options.ReadOptions)
	if err != nil {
		return nil, err
	}
	return extractDocument(doc, options)
}

func normalizeExtractOptions(options ExtractOptions) (ExtractOptions, error) {
	if options.Content == 0 {
		options.Content = ContentText
	}
	if options.Content & ^ContentAll != 0 {
		return options, fmt.Errorf("unknown extraction content flags: %d", options.Content)
	}
	if options.Glyphs {
		if options.Content&ContentText == 0 {
			return options, fmt.Errorf("glyph extraction requires text")
		}
		options.Positions = true
	}
	return options, nil
}

func extractDocument(doc *Document, options ExtractOptions) (*Extraction, error) {
	result := &Extraction{Content: options.Content, Pages: []ExtractedPage{}, options: options}
	if options.Provenance {
		result.Document = doc
	}
	_, err := BuildPDFEngine(doc, result)
	for _, diagnostic := range doc.Structure.Diagnostics {
		result.addDiagnostic(diagnostic, -1)
	}
	return result, err
}

func (b *semanticBuilder) wants(kind ContentKind) bool {
	return b.extract == nil || b.extract.Content&kind != 0
}

func (b *semanticBuilder) wantPositions() bool {
	return b.extract == nil || b.extract.options.Positions
}

func (b *semanticBuilder) wantStyles() bool {
	return b.extract == nil || b.extract.options.Styles
}

func (b *semanticBuilder) wantGlyphs() bool {
	return b.extract == nil || b.extract.options.Glyphs
}

func (b *semanticBuilder) wantProvenance() bool {
	return b.extract == nil || b.extract.options.Provenance
}

func (b *semanticBuilder) wantPaths() bool {
	return b.wants(ContentGraphics) || b.wantStyles()
}

func (b *semanticBuilder) wantTransforms() bool {
	return b.wantPositions() || b.wants(ContentGraphics) || b.wantStyles()
}

func (b *semanticBuilder) wantPageContent() bool {
	return b.wants(ContentText | ContentGraphics | ContentImages)
}

func (p *Extraction) addDiagnostic(diagnostic Diagnostic, page int) {
	d := ExtractionDiagnostic{
		Severity: diagnostic.Severity, Code: diagnostic.Code, Message: diagnostic.Message, Page: page,
	}
	if p.options.Provenance {
		d.Span = &diagnostic.Span
	}
	p.Diagnostics = append(p.Diagnostics, d)
}

// =====================================================================
// Request and result model: what callers select and what they receive.
// =====================================================================

// ContentKind selects independently emitted content categories.
type ContentKind uint8

const (
	ContentText ContentKind = 1 << iota
	ContentGraphics
	ContentImages
	ContentAnnotations
	ContentAll = ContentText | ContentGraphics | ContentImages | ContentAnnotations
)

// ExtractOptions controls work and retained output. Zero options extract Unicode
// text and compact font metadata. Glyphs requires text and enables Positions.
// Positions use unrotated page user space, as in the detailed API.
type ExtractOptions struct {
	Content     ContentKind
	Positions   bool
	Styles      bool
	Glyphs      bool
	Provenance  bool
	ReadOptions ReadOptions
}

// Extraction contains only requested output, in content execution order within
// each kind. It does not reconstruct reading order or render pixels. Treat
// results as read-only. Without Provenance, no Document or DetailedPDF is retained.
type Extraction struct {
	Content        ContentKind
	Pages          []ExtractedPage
	ImageResources []ExtractedImageResource `json:",omitempty"`
	Diagnostics    []ExtractionDiagnostic   `json:",omitempty"`
	// Document is present only with Provenance. Its lazy methods are not
	// concurrent-safe. It owns the snapshot; no Close call is required.
	Document *Document `json:"-"`
	options  ExtractOptions
}

// ExtractedPage groups selected elements. An omitted category was not requested
// or has no elements; Extraction.Content distinguishes those cases.
type ExtractedPage struct {
	Index             int
	MediaBox, CropBox Rect
	Rotate            int
	UserUnit          float64
	Texts             []ExtractedText       `json:",omitempty"`
	Graphics          []ExtractedGraphic    `json:",omitempty"`
	Images            []ExtractedImage      `json:",omitempty"`
	Annotations       []ExtractedAnnotation `json:",omitempty"`
	Operations        []Operation           `json:",omitempty"`
	// Complete concerns requested content, not validation of skipped resources.
	// Always check the returned error, including errors before a page is added.
	Complete bool
}

// ExtractedText is one text-show operation, with independently selected details.
type ExtractedText struct {
	Unicode        string
	Font           *FontInfo `json:",omitempty"`
	DecodeComplete bool
	Position       *TextPosition  `json:",omitempty"`
	Style          *PaintStyle    `json:",omitempty"`
	Glyphs         []Glyph        `json:",omitempty"`
	Source         *ElementSource `json:",omitempty"`
}

// TextPosition describes text before the show operation, before font scaling.
type TextPosition struct {
	FontSize      float64
	Matrix        Matrix
	RenderingMode int
	Complete      bool
}

// ExtractedGraphic is a painted path, not a complete chart or figure. Geometry
// is always in unrotated page user space, independent of Positions.
type ExtractedGraphic struct {
	Segments              []PathSegment
	Paint                 string
	Stroke, Fill, EvenOdd bool
	Style                 *PaintStyle    `json:",omitempty"`
	Source                *ElementSource `json:",omitempty"`
}

// ExtractedImage is one placement of a shared image resource.
type ExtractedImage struct {
	Resource int
	Matrix   *Matrix        `json:",omitempty"`
	Style    *PaintStyle    `json:",omitempty"`
	Source   *ElementSource `json:",omitempty"`
}

// ExtractedImageResource retains image metadata, not decoded pixels. ColorSpace
// preserves the PDF value (including complex color-space parameters). Object is
// present only with Provenance; its Stream can be used with Extraction.Document.
type ExtractedImageResource struct {
	ID               ObjectID
	Width, Height    int
	BitsPerComponent int
	ColorSpace       Object
	ImageMask        bool
	Object           *Object `json:",omitempty"`
}

// ExtractedAnnotation contains annotation metadata; appearance streams are not executed.
type ExtractedAnnotation struct {
	Subtype Name
	Rect    Rect
	Source  *Span `json:",omitempty"`
}

// ExtractionDiagnostic reports a requested interpretation limitation. Page is
// zero based, or -1 for a document-level issue. Span needs Provenance to resolve.
type ExtractionDiagnostic struct {
	Severity      Severity
	Code, Message string
	Page          int
	Span          *Span `json:",omitempty"`
}

// =====================================================================
// Emission: the only boundary that chooses a retained output model. The
// interpreter's transient state and resource caches are shared by both
// the basic/detailed and extraction APIs.
// =====================================================================

// Emission is the only boundary that chooses a retained output model. The
// interpreter's transient state and resource caches are shared by both APIs.
func (c *contentInterpreter) emitText(text DetailedText) {
	if c.b.extract == nil {
		c.item(ElementText, len(c.b.pdf.Texts))
		c.b.pdf.Texts = append(c.b.pdf.Texts, text)
		return
	}
	b := c.b
	output := ExtractedText{Unicode: text.Unicode, DecodeComplete: text.DecodeComplete}
	if text.Font >= 0 {
		if b.fontInfos == nil {
			b.fontInfos = make(map[int]*FontInfo)
		}
		info, ok := b.fontInfos[text.Font]
		if !ok {
			font := b.pdf.Fonts[text.Font]
			info = &FontInfo{BaseFont: font.BaseFont, Subtype: font.Subtype}
			b.fontInfos[text.Font] = info
		}
		output.Font = info
	}
	if b.wantPositions() {
		output.Position = &TextPosition{
			FontSize: text.FontSize, Matrix: text.Matrix,
			RenderingMode: text.RenderingMode, Complete: text.PositionComplete,
		}
	}
	if b.wantGlyphs() {
		output.Glyphs = text.Glyphs
	}
	output.Style = b.extractStyle(text.State)
	if b.wantProvenance() {
		source := text.Source
		output.Source = &source
	}
	page := &b.extract.Pages[c.page]
	page.Texts = append(page.Texts, output)
}

func (c *contentInterpreter) emitGraphic(graphic DetailedGraphic) {
	if c.b.extract == nil {
		c.item(ElementGraphic, len(c.b.pdf.Graphics))
		c.b.pdf.Graphics = append(c.b.pdf.Graphics, graphic)
		return
	}
	segments := make([]PathSegment, len(graphic.Segments))
	for i, segment := range graphic.Segments {
		// Path points are immutable after construction; the path accumulator is
		// replaced after paint, so output can own these without another copy.
		segments[i] = PathSegment{Operator: segment.Operator, Points: segment.Points}
	}
	output := ExtractedGraphic{
		Segments: segments, Paint: graphic.Paint, Stroke: graphic.Stroke,
		Fill: graphic.Fill, EvenOdd: graphic.EvenOdd, Style: c.b.extractStyle(graphic.State),
	}
	if c.b.wantProvenance() {
		source := graphic.Source
		output.Source = &source
	}
	page := &c.b.extract.Pages[c.page]
	page.Graphics = append(page.Graphics, output)
}

func (b *semanticBuilder) emitImageResource(resource ImageResource) int {
	if b.extract == nil {
		index := len(b.pdf.ImageResources)
		b.pdf.ImageResources = append(b.pdf.ImageResources, resource)
		return index
	}
	output := ExtractedImageResource{
		ID: resource.ID, Width: resource.Width, Height: resource.Height,
		BitsPerComponent: resource.BitsPerComponent, ColorSpace: resource.ColorSpace,
		ImageMask: resource.ImageMask,
	}
	if b.wantProvenance() {
		object := resource.Object
		output.Object = &object
	}
	index := len(b.extract.ImageResources)
	b.extract.ImageResources = append(b.extract.ImageResources, output)
	return index
}

func (c *contentInterpreter) emitImage(image DetailedImage) {
	if c.b.extract == nil {
		c.item(ElementImage, len(c.b.pdf.Images))
		c.b.pdf.Images = append(c.b.pdf.Images, image)
		return
	}
	output := ExtractedImage{Resource: image.Resource, Style: c.b.extractStyle(image.State)}
	if c.b.wantPositions() {
		// A pointer into image would also retain its graphics state and clips.
		matrix := image.Matrix
		output.Matrix = &matrix
	}
	if c.b.wantProvenance() {
		source := image.Source
		output.Source = &source
	}
	page := &c.b.extract.Pages[c.page]
	page.Images = append(page.Images, output)
}

func (b *semanticBuilder) extractStyle(state GraphicsState) *PaintStyle {
	if !b.wantStyles() {
		return nil
	}
	style := BasicStyle(state)
	return &style
}

func (b *semanticBuilder) emitAnnotation(annotation Annotation) {
	if b.extract == nil {
		index := len(b.pdf.Annotations)
		b.pdf.Annotations = append(b.pdf.Annotations, annotation)
		b.pdf.Pages[annotation.Page].Annotations = append(b.pdf.Pages[annotation.Page].Annotations, index)
		return
	}
	output := ExtractedAnnotation{Subtype: annotation.Subtype, Rect: annotation.Rect}
	if b.wantProvenance() {
		span := annotation.Object.Span
		output.Source = &span
	}
	page := &b.extract.Pages[annotation.Page]
	page.Annotations = append(page.Annotations, output)
}
