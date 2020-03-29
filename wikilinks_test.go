package goldmark_wikilinks

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/testutil"
	"github.com/yuin/goldmark/util"
	"testing"
)

func normalizer(linkText string) string {
	if linkText == "change me" {
		return "ChangeMe.html"
	}
	return linkText + ".html"
}



func TestWikilinks(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
				html.WithUnsafe(),
			),
	)
	markdown.Parser().AddOptions(
		parser.WithInlineParsers(util.Prioritized(NewWikilinksParser().WithNormalizer(normalizer), 102)),
	)
	testutil.DoTestCaseFile(markdown, "wikilinks.txt", t)
}