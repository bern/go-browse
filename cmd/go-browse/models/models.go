package models

import "strings"

// NodeType informs us of what kind of data we can expect the node to hold
type NodeType int

const (
	// Text NodeType for strings
	Text NodeType = iota
	// Element NodeType for everything else (lol)
	Element
)

// Node represents a node on a tree structure
type Node struct {
	Children []Node
	NodeType NodeType
	Text     *string
	Element  *ElementData
}

// ElementData represents an HTML element
type ElementData struct {
	TagName    string
	Attributes map[string]string
}

// ID returns the ID for a given element
func (e ElementData) ID() *string {
	if e.Attributes["id"] != "" {
		id := e.Attributes["id"]
		return &id
	}
	return nil
}

// Classes returns the class list for a given element
func (e ElementData) Classes() *[]string {
	if e.Attributes["class"] != "" {
		classStr := e.Attributes["class"]
		classes := strings.Split(classStr, " ")
		return &classes
	}
	return nil
}

// Stylesheet represents a set of CSS Rules
type Stylesheet struct {
	Rules []Rule
}

// Rule represents the selectors and declaractions that make a CSS Rule
type Rule struct {
	Selectors    []Selector
	Declarations []Declaration
}

// Selector represents the mapping for which elements certain CSS Rules should be applied
type Selector struct {
	TagName *string
	ID      *string
	Classes *[]string
}

// Declaration represents a <name, value> pair for CSS Properties
type Declaration struct {
	Name  string
	Value string // TODO: Extend into multiple types, should also be a slice
}

// PropertyMap represents a collection of CSS properties and their corresponding values
type PropertyMap map[string]string

// StyledNode is a DOM node that is associated with a map of CSS properties and values
type StyledNode struct {
	Node            Node
	SpecifiedValues PropertyMap
	Children        []StyledNode
}

// Specificity represents how specifically a selector applies to a certain element
type Specificity struct {
	IDSpecificity      int
	ClassSpecificity   int
	ElementSpecificity int
}

type MatchedRule struct {
	Rule        Rule
	Specificity Specificity
}
