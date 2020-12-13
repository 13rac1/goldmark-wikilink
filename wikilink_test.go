package wikilink_test

import (
	"testing"

	wikilink "github.com/13rac1/goldmark-wikilink"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/testutil"
)

func TestWikilink(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			wikilink.Extension,
		),
	)

	testutil.DoTestCaseFile(markdown, "wikilinks.txt", t)

}
