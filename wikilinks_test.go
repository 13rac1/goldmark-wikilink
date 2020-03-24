package goldmark_wikilinks

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/testutil"
	"testing"
)

func TestWikilinks(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
				html.WithUnsafe(),
			),
		goldmark.WithExtensions(Wikilinks),
	)
	testutil.DoTestCaseFile(markdown, "wikilinks.txt", t)
}