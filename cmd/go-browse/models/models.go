package models

import (
	"strings"
)

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

func (s *StyledNode) value(name string) *string {
	value := s.SpecifiedValues[name]

	if value == "" {
		return nil
	}
	return &value
}

// Display returns the value corresponding to the 'display' property on a StyledNode
func (s *StyledNode) Display() Display {
	displayValue := s.value("display")

	if displayValue == nil {
		return Inline
	}

	// should create default display types based on the element
	switch *displayValue {
	case "block":
		return Block
	case "none":
		return None
	default:
		return Inline
	}
}

// Specificity represents how specifically a selector applies to a certain element
type Specificity struct {
	IDSpecificity      int
	ClassSpecificity   int
	ElementSpecificity int
}

// MatchedRule is a rule and its specificity that was matched to a particular element
type MatchedRule struct {
	Rule        Rule
	Specificity Specificity
}

// Dimensions represents everything needed to position and layout a CSS Box
type Dimensions struct {
	Content Rectangle
	Padding EdgeSizes
	Border  EdgeSizes
	Margin  EdgeSizes
}

// Rectangle represents a 2D rectangle drawn in (x,y) space
type Rectangle struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

// EdgeSizes represents the distance drawn on each side from the previous content area layer
type EdgeSizes struct {
	Left   float32
	Right  float32
	Top    float32
	Bottom float32
}

// LayoutBox represents a node on the Layout Tree
type LayoutBox struct {
	Dimensions Dimensions
	BoxType    BoxType
	Node       *StyledNode
	Children   []LayoutBox
}

// NewLayoutBox is a constructor for a LayoutBox with a certain box type
func NewLayoutBox(boxType BoxType) LayoutBox {
	return LayoutBox{
		BoxType:  boxType,
		Children: make([]LayoutBox, 0),
	}
}

// GetInlineContainer is called when we need the proper container Box for an Inline Element
func (lb LayoutBox) GetInlineContainer() LayoutBox {
	// Switch based on the parent's BoxType
	switch lb.BoxType {
	case InlineNode: // Inline boxes can have inline children
		return lb
	case AnonymousBlock: // Anonymous boxes can have inline children
		return lb
	case BlockNode: // Blocks need to create an anonymous box to hold inline children
		if len(lb.Children) == 0 || lb.Children[len(lb.Children)-1].BoxType != AnonymousBlock { // If the latest child isn't an anonymous box...
			lb.Children = append(lb.Children, NewLayoutBox(AnonymousBlock)) // make it so
		}
		return lb.Children[len(lb.Children)-1] // Return the latest child as the thing that will contain this incoming inline elemen
	}

	return lb // shouldn't happen
}

// BoxType is an enum corresponding to the CSS Box Type of a LayoutBox
type BoxType int

const (
	// BlockNode corresponds to a Block element
	BlockNode BoxType = iota
	// InlineNode corresponds to an Inline element
	InlineNode
	// AnonymousBlock corresponds to an anonymous block
	AnonymousBlock
)

// Display is an enum containing supported values for the css display property
type Display int

const (
	// Inline corresponds to display:inline
	Inline Display = iota
	// Block corresponds to display:block
	Block
	// None corresponds to display:none
	None
)
