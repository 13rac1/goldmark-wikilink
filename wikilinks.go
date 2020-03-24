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
	panic("implement me")
}

type wikilinks struct {

}

var Wikilinks = &wikilinks{}

func (wl *wikilinks) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(util.Prioritized(defaultWikilinksParser, 102)),
	)
}