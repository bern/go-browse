package utils

import (
	"regexp"
	"strings"
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
		return s == " "
	})
}

// ConsumeName consumes and returns non-symbolic characters
func (p *Parser) ConsumeName() string {
	isValidChar := regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

	return p.ConsumeWhile(func(s string) bool {
		return isValidChar(s)
	})
}
