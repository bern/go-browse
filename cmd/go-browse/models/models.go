package models

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
