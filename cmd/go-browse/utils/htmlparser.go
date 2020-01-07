package utils

import (
	"strings"

	"github.com/bern/go-browse/cmd/go-browse/models"
)

// HTMLParser represents a parser tied to a specific HTML file
type HTMLParser struct {
	FilePath string
	Parser   *Parser
}

// assertStringParsed throws an error if a parsed string does not match its expected value
func (p *HTMLParser) assertStringParsed(actual string, expected ...string) {
	for _, test := range expected {
		if actual == test {
			return
		}
	}

	p.expectedStringError(expected...)
}

// ParseHTML parses an HTML source file and returns the root Node
func ParseHTML(filepath, src string) models.Node {
	parser := &HTMLParser{
		FilePath: filepath,
		Parser: &Parser{
			Input: src,
			Pos:   0,
		},
	}

	nodes := parser.ParseNodes()

	// We need to make sure the Node we return is the root
	if len(nodes) > 1 {
		return ElementNode("html", make(map[string]string, 0), nodes)
	}

	return nodes[0]
}

// ParseNodes recursively looks at all nodes until
// we hit either a closing tag or the end of our file
func (p *HTMLParser) ParseNodes() []models.Node {
	nodes := make([]models.Node, 0)

	for {
		p.Parser.ConsumeWhitespace()

		if p.Parser.EOF() || p.Parser.StartsWith("</") {
			break
		}

		if p.Parser.StartsWith("<!--") {
			p.ParseComment()
			continue
		}

		nodes = append(nodes, p.ParseNode())
	}

	return nodes
}

// ParseNode creates a Node based off of the parser's current position
func (p *HTMLParser) ParseNode() models.Node {
	switch p.Parser.NextChar() {
	case "<":
		return p.ParseElement()
	default:
		return p.ParseText()
	}
}

// ParseComment skims over an HTML comment and does not create a node for it
func (p *HTMLParser) ParseComment() {
	for {
		if p.Parser.EOF() || p.Parser.StartsWith("-->") {
			p.assertStringParsed(p.Parser.ConsumeChar(), "-")
			p.assertStringParsed(p.Parser.ConsumeChar(), "-")
			p.assertStringParsed(p.Parser.ConsumeChar(), ">")
			break
		}
		p.Parser.ConsumeChar()
	}
}

// ParseText creates a TextNode out of an uninterrupted text block
func (p *HTMLParser) ParseText() models.Node {
	text := p.Parser.ConsumeWhile(func(s string) bool {
		return s != "<"
	})

	return TextNode(strings.Trim(text, " \n\t"))
}

// ParseElement creates an ElementNode out of a set of open/close tags
func (p *HTMLParser) ParseElement() models.Node {
	name := ""
	attrs := make(map[string]string, 0)
	children := make([]models.Node, 0)

	// Opening tag
	p.assertStringParsed(p.Parser.ConsumeChar(), "<")
	name = p.Parser.ConsumeName()
	attrs = p.ParseAttributes()

	// Check for self-closing tag
	if p.Parser.StartsWith("/>") {
		p.assertStringParsed(p.Parser.ConsumeChar(), "/")
		p.assertStringParsed(p.Parser.ConsumeChar(), ">")

		return ElementNode(name, attrs, children)
	}

	p.assertStringParsed(p.Parser.ConsumeChar(), ">")

	// Contents
	children = p.ParseNodes()

	// Closing tag
	p.assertStringParsed(p.Parser.ConsumeChar(), "<")
	p.assertStringParsed(p.Parser.ConsumeChar(), "/")
	p.assertStringParsed(p.Parser.ConsumeName(), name)
	p.assertStringParsed(p.Parser.ConsumeChar(), ">")

	return ElementNode(name, attrs, children)
}

// ParseAttributes retrieves and maps all key="value" pairs in an element tag
func (p *HTMLParser) ParseAttributes() map[string]string {
	attrs := make(map[string]string, 0)

	for {
		p.Parser.ConsumeWhitespace()
		if p.Parser.NextChar() == ">" || p.Parser.StartsWith("/>") {
			break
		}

		name, value := p.ParseAttribute()
		attrs[name] = value
	}

	return attrs
}

// ParseAttribute returns a string pair corresponding to Name,Value of an attribute
func (p *HTMLParser) ParseAttribute() (string, string) {
	name := p.Parser.ConsumeName()
	p.assertStringParsed(p.Parser.ConsumeChar(), "=")
	value := p.ConsumeAttributeValue()
	return name, value
}

// ConsumeAttributeValue expects a quoted string and returns it
func (p *HTMLParser) ConsumeAttributeValue() string {
	openQuote := p.Parser.ConsumeChar()
	p.assertStringParsed(openQuote, "'", `"`)
	value := p.Parser.ConsumeWhile(func(s string) bool {
		return s != openQuote
	})
	p.assertStringParsed(p.Parser.ConsumeChar(), openQuote)
	return value
}
