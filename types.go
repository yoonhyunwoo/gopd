package gopd

import (
	"github.com/MyungSub0519/gopd/internal/document"
	"github.com/MyungSub0519/gopd/internal/pdfmodel"
)

// Document owns an immutable input snapshot and lazily decoded sources.
// Treat returned objects/slices as read-only. Lazy methods are not concurrent-safe.
// ReadOptions bounds input, recursive parsing, xrefs and decoded stream data.
// Zero fields use defaults. Negative values are invalid.
type ReadOptions = document.ReadOptions

type Point = pdfmodel.Point

type Rect = pdfmodel.Rect

// Matrix is [a b c d e f], mapping (x,y) to (a*x+c*y+e,b*x+d*y+f).
// Coordinates use PDF user space; page rotation is retained separately.
type Matrix = pdfmodel.Matrix

type SourceID = pdfmodel.SourceID

// Source == 0은 위치 없음. 오프셋은 항상 바이트 단위다.
type Position = pdfmodel.Position

// [Start, End): Start를 포함하고 End를 제외한다.
type Span = pdfmodel.Span

const (
	TransformFilter  = pdfmodel.TransformFilter
	TransformDecrypt = pdfmodel.TransformDecrypt
)

type Value = pdfmodel.Value

type Object = pdfmodel.Object

type Null = pdfmodel.Null

type Boolean = pdfmodel.Boolean

// PDF 숫자 문법을 검증한 십진 표기. 자동으로 float64로 바꾸지 않는다.
type Integer = pdfmodel.Integer

type Real = pdfmodel.Real

// '/'를 제외하고 #xx escape를 해제한 바이트열. UTF-8을 보장하지 않는다.
type Name = pdfmodel.Name

type StringForm = pdfmodel.StringForm

type PDFString = pdfmodel.PDFString

type Array = pdfmodel.Array

type DictionaryEntry = pdfmodel.DictionaryEntry

type Dictionary = pdfmodel.Dictionary

type ObjectID = pdfmodel.ObjectID

type Reference = pdfmodel.Reference

// 손상된 구문을 표현하는 분석기 전용 값. PDF의 null과 구별한다.
const (
	StringLiteral = pdfmodel.StringLiteral
	StringHex     = pdfmodel.StringHex
)

type StreamBoundary = pdfmodel.StreamBoundary

type Stream = pdfmodel.Stream

const (
	StreamUnresolved = pdfmodel.StreamUnresolved
	StreamFromLength = pdfmodel.StreamFromLength
	StreamRecovered  = pdfmodel.StreamRecovered
)

const (
	XRefTable  = pdfmodel.XRefTable
	XRefStream = pdfmodel.XRefStream
)

const (
	TokenInvalid       = pdfmodel.TokenInvalid
	TokenEOF           = pdfmodel.TokenEOF
	TokenWhitespace    = pdfmodel.TokenWhitespace
	TokenComment       = pdfmodel.TokenComment
	TokenInteger       = pdfmodel.TokenInteger
	TokenReal          = pdfmodel.TokenReal
	TokenName          = pdfmodel.TokenName
	TokenLiteralString = pdfmodel.TokenLiteralString
	TokenHexString     = pdfmodel.TokenHexString
	TokenArrayOpen     = pdfmodel.TokenArrayOpen
	TokenArrayClose    = pdfmodel.TokenArrayClose
	TokenDictOpen      = pdfmodel.TokenDictOpen
	TokenDictClose     = pdfmodel.TokenDictClose
	TokenKeyword       = pdfmodel.TokenKeyword
)

type Severity = pdfmodel.Severity

type Diagnostic = pdfmodel.Diagnostic

type Limits = pdfmodel.Limits

const (
	RegionUnknown        = pdfmodel.RegionUnknown
	RegionHeader         = pdfmodel.RegionHeader
	RegionTrivia         = pdfmodel.RegionTrivia
	RegionIndirectObject = pdfmodel.RegionIndirectObject
	RegionXRef           = pdfmodel.RegionXRef
	RegionTrailer        = pdfmodel.RegionTrailer
	RegionFileTail       = pdfmodel.RegionFileTail
)

const (
	SeverityInfo    = pdfmodel.SeverityInfo
	SeverityWarning = pdfmodel.SeverityWarning
	SeverityError   = pdfmodel.SeverityError
)

// ErrMissingKey identifies absent dictionary keys. Use errors.Is to distinguish
// a missing key from a duplicate key or another lookup failure.
var ErrMissingKey = pdfmodel.ErrMissingKey

// ErrLimit identifies exhausted resource budgets. Use errors.Is even when
// parsing also returns a partial result.
var ErrLimit = pdfmodel.ErrLimit
