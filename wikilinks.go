package goldmark_wikilinks

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"log"
)

type FilenameNormalizer interface {
	Normalize(linkText string) string
}

type WikilinkTracker interface {
	LinkWithContext(destination string, context string)
}

type wikilinksParser struct {
	normalizer FilenameNormalizer
	tracker WikilinkTracker
}

var defaultWikilinksParser = &wikilinksParser{}

func NewWikilinksParser() *wikilinksParser {
	return defaultWikilinksParser
}

func (wl *wikilinksParser) WithNormalizer(fn FilenameNormalizer) *wikilinksParser {
	log.Println("Normalizer Set")
	wl.normalizer = fn
	return wl
}

func (wl *wikilinksParser) WithTracker(wlt WikilinkTracker) *wikilinksParser {
	wl.tracker = wlt
	return wl
}

func (wl *wikilinksParser) Trigger() []byte {
	return []byte{'[', '['}
}

func (wl *wikilinksParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, segment := block.PeekLine()
	if line[1] != '[' {
		return nil
	}
	gotFirst := false
	pos := 2
	for ; pos < len(line); pos++ {
		b := line[pos]
		if b == ']' {
			if gotFirst {
				break
			} else {
				gotFirst = true
			}
		} else if gotFirst {
			gotFirst = false
		}
	}

	if !gotFirst && pos >= len(line) {
		return nil
	}

	destination := block.Value(text.NewSegment(segment.Start+2, segment.Start+pos-1))
	destText := string(destination)
	if wl.normalizer != nil {
		log.Println("Calling normalizer for " + destText)
		destText = wl.normalizer.Normalize(destText)
	} else {
		destText += ".html"
	}
	destination = []byte(destText)

	if wl.tracker != nil {
		context := ""
		lines := parent.Lines()
		for i := 0; i < lines.Len(); i++ {
			seg := lines.At(i)
			context += string(block.Value(seg))
		}
		wl.tracker.LinkWithContext(destText, context)
	}

	block.Advance(pos+1)

	link := ast.NewLink()
	link.Destination = destination
	newText := ast.NewText()
	newText.Segment = text.NewSegment(segment.Start+2, segment.Start+pos-1)
	link.AppendChild(link, newText)
	return link
}

type wikilinks struct {

}

var Wikilinks = &wikilinks{}

func (wl *wikilinks) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(util.Prioritized(defaultWikilinksParser, 102)),
	)
}