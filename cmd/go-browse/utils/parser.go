package utils

import (
	"regexp"
	"strings"

	"github.com/bern/go-browse/cmd/go-browse/models"
)

// Parser represents a generic parser with an input string and a tracked position
type Parser struct {
	Pos   int
	Input string
}

// NextChar reveals the next character that will be visited by the parser
func (p *Parser) NextChar() string {
	return string(p.Input[p.Pos])
}

// StartsWith checks whether the input at its current position has a given prefix
func (p *Parser) StartsWith(prefix string) bool {
	return strings.HasPrefix(p.Input[p.Pos:], prefix)
}

// EOF returns true if all input is consumed
func (p *Parser) EOF() bool {
	return p.Pos >= len(p.Input)
}

// ConsumeChar advances the parser and returns the consumed character
func (p *Parser) ConsumeChar() string {
	if p.EOF() {
		return ""
	}

	nextChar := p.NextChar()
	p.Pos = p.Pos + 1
	return nextChar
}

// ConsumeWhile advances the parser and returns a sequence of
// consumed characters until a given test fails
func (p *Parser) ConsumeWhile(test func(string) bool) string {
	result := ""
	for !p.EOF() && test(p.NextChar()) {
		result += p.ConsumeChar()
	}
	return result
}

// ConsumeWhitespace consumes zero or more whitespace characters
func (p *Parser) ConsumeWhitespace() {
	p.ConsumeWhile(func(s string) bool {
		return s == " " || s == "\n" || s == "\t"
	})
}

// ConsumeName consumes and returns non-symbolic characters
func (p *Parser) ConsumeName() string {
	isValidChar := regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

	return p.ConsumeWhile(func(s string) bool {
		return isValidChar(s)
	})
}

// ParseHTML parses an HTML source file and returns the root Node
func ParseHTML(src string) models.Node {
	parser := &Parser{
		Input: src,
		Pos:   0,
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
func (p *Parser) ParseNodes() []models.Node {
	nodes := make([]models.Node, 0)

	for {
		p.ConsumeWhitespace()

		if p.EOF() || p.StartsWith("</") {
			break
		}

		nodes = append(nodes, p.ParseNode())
	}

	return nodes
}

// ParseNode creates a Node based off of the parser's current position
func (p *Parser) ParseNode() models.Node {
	switch p.NextChar() {
	case "<":
		return p.ParseElement()
	default:
		return p.ParseText()
	}
}

// ParseText creates a TextNode out of an uninterrupted text block
func (p *Parser) ParseText() models.Node {
	text := p.ConsumeWhile(func(s string) bool {
		return s != "<"
	})

	return TextNode(strings.Trim(text, " \n\t"))
}

// ParseElement creates an ElementNode out of a set of open/close tags
func (p *Parser) ParseElement() models.Node {
	name := ""
	attrs := make(map[string]string, 0)
	children := make([]models.Node, 0)

	// Opening tag
	p.ConsumeChar() // <
	name = p.ConsumeName()
	attrs = p.ParseAttributes()
	p.ConsumeChar() // >

	// Contents
	children = p.ParseNodes()

	// Closing tag
	p.ConsumeChar() // <
	p.ConsumeChar() // /
	p.ConsumeName() // same val as name
	p.ConsumeChar() // >

	return ElementNode(name, attrs, children)
}

// ParseAttributes retrieves and maps all key="value" pairs in an element tag
func (p *Parser) ParseAttributes() map[string]string {
	attrs := make(map[string]string, 0)

	for {
		p.ConsumeWhitespace()
		if p.NextChar() == ">" {
			break
		}

		name, value := p.ParseAttribute()
		attrs[name] = value
	}

	return attrs
}

// ParseAttribute returns a string pair corresponding to Name,Value of an attribute
func (p *Parser) ParseAttribute() (string, string) {
	name := p.ConsumeName()
	p.ConsumeChar() // =
	value := p.ConsumeAttributeValue()
	return name, value
}

// ConsumeAttributeValue expects a quoted string and returns it
func (p *Parser) ConsumeAttributeValue() string {
	openQuote := p.ConsumeChar() // ' or "
	value := p.ConsumeWhile(func(s string) bool {
		return s != openQuote
	})
	p.ConsumeChar() // same val as openQuote
	return value
}
