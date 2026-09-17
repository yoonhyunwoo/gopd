package content

import (
	"github.com/MyungSub0519/gopd/internal/document"
	"github.com/MyungSub0519/gopd/internal/pdfmodel"
	"github.com/MyungSub0519/gopd/internal/syntax"
)

// Bridge to the shared data model and lower layers; keeps engine code free
// of qualifiers.
type (
	Value          = pdfmodel.Value
	Object         = pdfmodel.Object
	Null           = pdfmodel.Null
	Boolean        = pdfmodel.Boolean
	Integer        = pdfmodel.Integer
	Real           = pdfmodel.Real
	Name           = pdfmodel.Name
	PDFString      = pdfmodel.PDFString
	Array          = pdfmodel.Array
	Dictionary     = pdfmodel.Dictionary
	ObjectID       = pdfmodel.ObjectID
	Reference      = pdfmodel.Reference
	Stream         = pdfmodel.Stream
	StreamBoundary = pdfmodel.StreamBoundary
	SourceID       = pdfmodel.SourceID
	Position       = pdfmodel.Position
	Span           = pdfmodel.Span
	Source         = pdfmodel.Source
	Derivation     = pdfmodel.Derivation
	Transform      = pdfmodel.Transform
	Diagnostic     = pdfmodel.Diagnostic
	Severity       = pdfmodel.Severity
	Limits         = pdfmodel.Limits
	Structure      = pdfmodel.Structure
	Point          = pdfmodel.Point
	Rect           = pdfmodel.Rect
	Matrix         = pdfmodel.Matrix
	Token          = pdfmodel.Token
	TokenKind      = pdfmodel.TokenKind
	Document       = document.Document
	ReadOptions    = document.ReadOptions
)

const (
	TokenWhitespace = pdfmodel.TokenWhitespace
	TokenComment    = pdfmodel.TokenComment
	TokenEOF        = pdfmodel.TokenEOF
	TokenInteger    = pdfmodel.TokenInteger
	TokenReal       = pdfmodel.TokenReal
	TokenName       = pdfmodel.TokenName
	TokenKeyword    = pdfmodel.TokenKeyword
	SeverityInfo    = pdfmodel.SeverityInfo
	SeverityWarning = pdfmodel.SeverityWarning
	SeverityError   = pdfmodel.SeverityError
)

var (
	ErrLimit       = pdfmodel.ErrLimit
	ErrMissingKey  = pdfmodel.ErrMissingKey
	Lex            = syntax.Lex
	Parse          = document.Parse
	ParseFile      = document.ParseFile
	Int            = pdfmodel.Int
	Number         = pdfmodel.Number
	IsStream       = pdfmodel.IsStream
	IdentityMatrix = pdfmodel.IdentityMatrix
)
