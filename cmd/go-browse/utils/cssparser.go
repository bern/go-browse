package utils

import (
	"github.com/bern/go-browse/cmd/go-browse/models"
)

// CSSParser represents a parser tied to a specific CSS file
type CSSParser struct {
	FilePath string
	Parser   *Parser
}

// ParseCSS parses a CSS source file and returns a Stylesheet
func ParseCSS(filepath, src string) models.Stylesheet {
	parser := &CSSParser{
		FilePath: filepath,
		Parser: &Parser{
			Input: src,
			Pos:   0,
		},
	}

	stylesheet := parser.ParseRules()

	return stylesheet
}

// ParseRules crawls down a CSS file parsing each rule as it is encountered
func (p *CSSParser) ParseRules() models.Stylesheet {
	stylesheet := models.Stylesheet{
		Rules: make([]models.Rule, 0),
	}

	for {
		p.Parser.ConsumeWhitespace()

		if p.Parser.EOF() {
			break
		}

		/*
			div#about.main__info.main__info--blue {
				color: blue;
			}
		*/

		// parse selector
		// TODO: Support multiple comma-separated selectors
		selector := p.ParseSelector()

		// parse declarations
		declarations := p.ParseDeclarations()

		// create a rule
		stylesheet.Rules = append(stylesheet.Rules, models.Rule{
			Selectors:    []models.Selector{selector},
			Declarations: declarations,
		})
	}

	return stylesheet
}

func (p *CSSParser) ParseSelector() models.Selector {
	/*
		A simple selector is either a type selector or universal selector
		followed immediately by zero or more attribute selectors, ID selectors,
		or pseudo-classes, in any order.
	*/

	selector := models.Selector{}

	// if the first character is neither a . nor a # nor a * then it is the type selector
	if p.Parser.NextChar() != "." && p.Parser.NextChar() != "#" && p.Parser.NextChar() != "*" {
		tagName := p.Parser.ConsumeWhile(func(s string) bool {
			return s != "." && s != "#" && s != "*"
		})
		selector.TagName = &tagName
	}

	classes := make([]string, 0)

	for {
		p.Parser.ConsumeWhitespace()
		if p.Parser.NextChar() == "{" {
			break
		}

		switch p.Parser.NextChar() {
		case ".":
			p.Parser.ConsumeChar() // .
			className := p.ConsumeSelectorComponent()

			classes = append(classes, className)
			break
		case "#":
			p.Parser.ConsumeChar() // #
			idName := p.ConsumeSelectorComponent()

			selector.ID = &idName
			break
		case "*":
			p.Parser.ConsumeChar() // *
			break
		default: // unrecognized
		}
	}

	if len(classes) > 0 {
		selector.Classes = &classes
	}

	return selector
}

func (p *CSSParser) ConsumeSelectorComponent() string {
	return p.Parser.ConsumeWhile(func(s string) bool {
		return s != "." && s != "#" && s != "{" && s != " " && s != "\n" && s != "\t"
	})
}

func (p *CSSParser) ParseDeclarations() []models.Declaration {
	declarations := make([]models.Declaration, 0)

	p.Parser.ConsumeChar() // {
	for {
		p.Parser.ConsumeWhitespace()

		if p.Parser.EOF() { // TODO: check if this should throw an error??
			break
		} else if p.Parser.NextChar() == "}" {
			p.Parser.ConsumeChar() // }
			break
		}

		declaration := p.ParseDeclaration()
		declarations = append(declarations, declaration)
	}

	return declarations
}

func (p *CSSParser) ParseDeclaration() models.Declaration {
	name := p.Parser.ConsumeWhile(func(s string) bool {
		return s != ":"
	})

	p.Parser.ConsumeChar() // :
	p.Parser.ConsumeWhitespace()

	// TODO: support a number, a string, a percentage, or a hex color code
	value := p.Parser.ConsumeWhile(func(s string) bool {
		return s != ";"
	})

	p.Parser.ConsumeChar() // ;

	return models.Declaration{
		Name:  name,
		Value: value,
	}
}
