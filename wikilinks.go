package goldmark_wikilinks

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// FilenameNormalizer is a plugin which takes link text and converts the text given to
// a filename which can be linked to in the final format of your file.
type FilenameNormalizer interface {
	Normalize(linkText string) string
}

// WikilinkTracker is a plugin that can get called for each discovered link and gather up
// information about the links. This is useful for creating backlinks, for example (the
// purpose for which I created this plugin).
type WikilinkTracker interface {
	LinkWithContext(destText string, destFilename string, context string)
}

// wikilinksParser keeps track of the plugins used for processing wikilinks.
type wikilinksParser struct {
	normalizer FilenameNormalizer
	tracker WikilinkTracker
}

var defaultWikilinksParser = &wikilinksParser{}

// NewWikilinksParser gives you back a parser that you can use to process wikilinks.
func NewWikilinksParser() *wikilinksParser {
	return defaultWikilinksParser
}

// WithNormalizer is the fluent interface for adding a normalizer plugin
func (wl *wikilinksParser) WithNormalizer(fn FilenameNormalizer) *wikilinksParser {
	wl.normalizer = fn
	return wl
}

// WithTracker is the fluent interface for adding a wikilink tracker plugin
func (wl *wikilinksParser) WithTracker(wlt WikilinkTracker) *wikilinksParser {
	wl.tracker = wlt
	return wl
}

// Trigger looks for the [[ beginning of wikilinks.
func (wl *wikilinksParser) Trigger() []byte {
	return []byte{'[', '['}
}

func (wl *wikilinksParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, segment := block.PeekLine()
	// Did we not actually find a wikilink?
	if line[1] != '[' {
		return nil
	}
	gotFirst := false
	pos := 2
	for ; pos < len(line); pos++ {
		b := line[pos]
		// look for two ]] to close out the wikilink
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

	// extract the text of the wikilink and normalize it to figure out the destination
	destination := block.Value(text.NewSegment(segment.Start+2, segment.Start+pos-1))
	destText := string(destination)
	destFilename := destText
	if wl.normalizer != nil {
		destFilename = wl.normalizer.Normalize(destFilename)
	} else {
		destFilename += ".html"
	}
	destination = []byte(destFilename)

	if wl.tracker != nil {
		context := ""
		lines := parent.Lines()
		for i := 0; i < lines.Len(); i++ {
			seg := lines.At(i)
			context += string(block.Value(seg))
		}
		wl.tracker.LinkWithContext(destText, destFilename, context)
	}

	block.Advance(pos+1)

	// This replaces the wikilink in the AST with a normal markdown link
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

// Extend adds a wikilink parser to a Goldmark parser
func (wl *wikilinks) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(util.Prioritized(defaultWikilinksParser, 102)),
	)
}