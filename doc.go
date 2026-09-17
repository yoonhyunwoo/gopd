// Package gopd parses PDF page content and preserves its underlying structure.
//
// # Selective content
//
// Extract and ExtractReader return an Extraction grouped by page, generating
// only selected content. Zero ExtractOptions select Unicode text. Content flags
// combine text, graphics, images, and annotations; Positions, Styles, Glyphs,
// and Provenance select details. Glyphs requires text and enables Positions.
// Graphics always include path geometry. Without Provenance, the returned
// extraction retains neither a Document nor a DetailedPDF. Provenance retains
// source access and page operations; its Document is excluded from JSON.
//
// # Basic content
//
// ParsePDF interprets only text and graphics and returns them grouped by page,
// retaining neither the input snapshot nor a detailed result. Use it when a
// caller only needs everyday text and path content.
//
// # Detailed content
//
// Open and Read return a DetailedPDF containing page geometry, text glyphs,
// graphics, images, fonts, annotations, content operations, and source spans.
//
// # Selective content
//
// Extract and ExtractReader return an Extraction grouped by page, generating
// only selected content. Content flags combine text, graphics, images, and
// annotations; Positions, Styles, Glyphs, and Provenance select details.
// Without Provenance, the returned extraction retains neither a Document nor
// a DetailedPDF.
//
// # Ownership and interpretation limits
//
// Parsing snapshots input and caches decoded sources in memory during a call.
// Detailed and basic results retain the snapshot; extraction retains it only
// with Provenance. This is not streaming I/O or an exact process-memory bound.
// No Close call is needed. Treat results as read-only. Lazy Document methods
// are not safe for concurrent calls.
// Basic parsing interprets only the requested text and graphics, but it is
// not a low-memory mode for detailed work. Always check returned errors,
// including when a partial result is non-nil.
// Unsupported effects may instead be reported in result diagnostics. Selective
// extraction does not validate skipped resources or unrequested interpretation.
// Coordinates use unrotated PDF user space, and content order is drawing order,
// not reconstructed reading order. Rendering, OCR, and PDF editing are outside
// the current implementation.
package gopd
