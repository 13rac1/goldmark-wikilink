# goldmark-wikilink

goldmark-wikilink is an extension for the [goldmark][goldmark] library that extends
Markdown to support `[[title]]` [Wikilink][help-link] style links with a new AST
type and HTML Renderer.

[goldmark]: http://github.com/yuin/goldmark
[help-link]: https://en.wikipedia.org/wiki/Help:Link
[goldmark-wikilinks]: https://github.com/dangoor/goldmark-wikilinks

## Demo

This markdown:

```md
# Hello goldmark-wikilink

[[Example Link]]
```

Becomes this HTML, with the default Normalizer:

```html
<h1>Hello goldmark-wikilink</h1>
<p><a href="Example%20Link.html">Example Link</a></p>
```

### Installation

```bash
go get github.com/13rac1/goldmark-wikilink
```

## Usage

```go
  markdown := goldmark.New(
    goldmark.WithExtensions(
      wikilink.Extension,
    ),
  )
  var buf bytes.Buffer
  if err := markdown.Convert([]byte(source), &buf); err != nil {
    panic(err)
  }
  fmt.Print(buf)
}
```

## TODO

* Support [Piped Links][piped-link] in the form `[target|displayed text]`.
* Support [Section linking][section-linking] in the forms:
  * External link: `[[Page name#Section name|displayed text]]`
  * Internal link: `[[#Section name|displayed text]]`

[piped-link]: https://en.wikipedia.org/wiki/Help:Link#Piped_link
[section-linking]: https://en.wikipedia.org/wiki/Help:Link#Section_linking_(anchors)

## License

MIT

## Author

Brad Erickson & Kevin Dangoor

Fork of [dangoor/goldmark-wikilinks][goldmark-wikilinks]. Adds
a Wikilink AST type to ease extending and removes the _WikilinkTracker_ which
can implemented in a separate AST walker.
