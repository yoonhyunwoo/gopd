package gopd

import (
	"io"

	"github.com/MyungSub0519/gopd/internal/content"
)

// Aliases re-export the detailed model and interpretation results. The
// analysis engine lives under internal/; these names are the supported
// public surface.
type (
	DetailedPDF            = content.DetailedPDF
	DetailedPage           = content.DetailedPage
	ElementKind            = content.ElementKind
	ElementRef             = content.ElementRef
	ElementSource          = content.ElementSource
	FormCall               = content.FormCall
	Operation              = content.Operation
	DetailedPathSegment    = content.DetailedPathSegment
	DetailedGraphic        = content.DetailedGraphic
	Glyph                  = content.Glyph
	DetailedText           = content.DetailedText
	DetailedImage          = content.DetailedImage
	ImageResource          = content.ImageResource
	Annotation             = content.Annotation
	Color                  = content.Color
	PaintStyle             = content.PaintStyle
	GraphicsState          = content.GraphicsState
	ClipPath               = content.ClipPath
	PathSegment            = content.PathSegment
	FontInfo               = content.FontInfo
	ContentKind            = content.ContentKind
	ExtractOptions         = content.ExtractOptions
	Extraction             = content.Extraction
	ExtractedPage          = content.ExtractedPage
	ExtractedText          = content.ExtractedText
	TextPosition           = content.TextPosition
	ExtractedGraphic       = content.ExtractedGraphic
	ExtractedImage         = content.ExtractedImage
	ExtractedImageResource = content.ExtractedImageResource
	ExtractedAnnotation    = content.ExtractedAnnotation
	ExtractionDiagnostic   = content.ExtractionDiagnostic
)

const (
	ContentText        = content.ContentText
	ContentGraphics    = content.ContentGraphics
	ContentImages      = content.ContentImages
	ContentAnnotations = content.ContentAnnotations
	ContentAll         = content.ContentAll
	ElementText        = content.ElementText
	ElementGraphic     = content.ElementGraphic
	ElementImage       = content.ElementImage
)

// Extract snapshots a file and directly emits selected content.
func Extract(path string, options ExtractOptions) (*Extraction, error) {
	return content.Extract(path, options)
}

// ExtractReader is the ReaderAt form of Extract.
func ExtractReader(r io.ReaderAt, size int64, options ExtractOptions) (*Extraction, error) {
	return content.ExtractReader(r, size, options)
}
