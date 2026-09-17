package content

import "github.com/MyungSub0519/gopd/internal/pdftest"

func semanticStream(dict, content string) string {
	return pdftest.Stream(dict, content)
}
