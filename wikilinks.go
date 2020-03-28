package goldmark_wikilinks

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type wikilinksParser struct {

}

var defaultWikilinksParser = &wikilinksParser{}

func NewWikilinksParser() parser.InlineParser {
	return defaultWikilinksParser
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
	destination = append(destination, '.')
	destination = append(destination, 'h')
	destination = append(destination, 't')
	destination = append(destination, 'm')
	destination = append(destination, 'l')

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