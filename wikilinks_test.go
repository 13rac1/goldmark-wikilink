package goldmark_wikilinks

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/testutil"
	"github.com/yuin/goldmark/util"
	"testing"
)

type Backlink struct {
	destination string
	context string
}

type Tracker struct {
	backlinks []Backlink
}

func (t *Tracker) LinkWithContext(dest string, context string) {
	bl := Backlink{
		destination: dest,
		context:     context,
	}
	t.backlinks = append(t.backlinks, bl)
}

func (t *Tracker) Normalize(linkText string) string {
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
	tracker := &Tracker{}

	markdown.Parser().AddOptions(
		parser.WithInlineParsers(util.Prioritized(NewWikilinksParser().
			WithNormalizer(tracker).WithTracker(tracker), 102)),
	)
	testutil.DoTestCaseFile(markdown, "wikilinks.txt", t)
	if len(tracker.backlinks) != 7 {
		t.Errorf("Expected 7 backlinks but saw %d", len(tracker.backlinks))
	}
	listlink := tracker.backlinks[5]
	if listlink.context != "That has a [[Wiki Link]] in the second bullet" {
		t.Errorf("Did not get expected context for bullet: %s", listlink.context)
	}
	paralink := tracker.backlinks[6]
	if paralink.context != `Here is a multi-line paragraph which is full of text and
also has a [[Wiki Link]] in the middle of it, but I should
get the _full set_ of text in my tracker.` {
		t.Errorf("Did not get expected paragraph context: %s", paralink.context,)
	}
}