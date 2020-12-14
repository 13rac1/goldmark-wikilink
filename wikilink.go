package wikilink

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// FilenameNormalizer is a plugin which takes link text and converts the text given to
// a filename which can be linked to in the final format of your file.
type FilenameNormalizer interface {
	Normalize(linkText string) string
}

// wikilinksParser keeps track of the plugins used for processing wikilinks.
type wikilinksParser struct {
	normalizer FilenameNormalizer
}

type linkNormalizer struct{}

func (t *linkNormalizer) Normalize(linkText string) string {
	return url.PathEscape(linkText) + ".html"
}

var defaultWikilinksParser = &wikilinksParser{
	normalizer: &linkNormalizer{},
}

// NewParser gives you back a parser that you can use to process wikilinks.
func NewParser() *wikilinksParser {
	return defaultWikilinksParser
}

// WithNormalizer is the fluent interface for replacing the default normalizer plugin.
func (p *wikilinksParser) WithNormalizer(fn FilenameNormalizer) *wikilinksParser {
	p.normalizer = fn
	return p
}

// Trigger looks for the [[ beginning of wikilinks.
func (p *wikilinksParser) Trigger() []byte {
	return []byte{'[', '['}
}

func (p *wikilinksParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, segment := block.PeekLine()
	// Must specifically confirm the second '['.
	if line[1] != '[' {
		return nil
	}

	foundEnd := false
	// Skip to 3rd position since first two must be `[[`
	pos := 2
	// Look for two ']]' to end the wikilink.
	for ; pos < len(line)-1; pos++ {
		if line[pos] != ']' {
			continue
		}
		// Can always add one, because of the -1
		pos++
		if line[pos] != ']' {
			continue
		}
		// pos == the position of the second ']'
		foundEnd = true
		break
	}

	if !foundEnd {
		return nil
	}

	// extract the text of the wikilink and normalize it to figure out the destination
	destination := block.Value(text.NewSegment(segment.Start+2, segment.Start+pos-1))
	destString := strings.TrimSpace(string(destination))
	destFilename := p.normalizer.Normalize(destString)

	block.Advance(pos + 1)

	link := ast.NewLink()
	link.Title = []byte(destString)
	link.Destination = []byte(destFilename)
	wl := NewWikilink(link)
	return wl
}

type wlExtension struct {
}

// Extension is the default extension instance.
var Extension = &wlExtension{}

// Extend adds a wikilink parser to a Goldmark parser
func (wl *wlExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(util.Prioritized(defaultWikilinksParser, 102)),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(NewHTMLRenderer(), 500),
		),
	)
}

// Wikilink struct represents a Wikilink of the Markdown text.
type Wikilink struct {
	ast.Link
}

// KindWikilink is a NodeKind of the Wikilink node.
var KindWikilink = ast.NewNodeKind("Wikilink")

// Kind implements Node.Kind.
func (n *Wikilink) Kind() ast.NodeKind {
	return KindWikilink
}

// NewWikilink returns a new Wikilink node.
func NewWikilink(l *ast.Link) *Wikilink {
	c := &Wikilink{
		Link: *l,
	}
	c.Destination = l.Destination // AKA Target
	c.Title = l.Title

	return c
}

// HTMLRenderer struct is a renderer.NodeRenderer implementation for the extension.
type HTMLRenderer struct{}

// NewHTMLRenderer builds a new HTMLRenderer with given options and returns it.
func NewHTMLRenderer() renderer.NodeRenderer {
	return &HTMLRenderer{}
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs.
func (r *HTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindWikilink, r.render)
}

func (r *HTMLRenderer) render(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		return ast.WalkContinue, nil
	}

	wl := node.(*Wikilink)
	out := fmt.Sprintf(`<a href="%s">%s</a>`, wl.Destination, wl.Title)
	w.Write([]byte(out))
	return ast.WalkContinue, nil
}
