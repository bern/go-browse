package models

import (
	"fmt"
	"strconv"
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
	X      int // should be float32
	Y      int
	Width  int
	Height int
}

// EdgeSizes represents the distance drawn on each side from the previous content area layer
type EdgeSizes struct {
	Left   int
	Right  int
	Top    int
	Bottom int
}

// LayoutBox represents a node on the Layout Tree
type LayoutBox struct {
	Dimensions Dimensions
	BoxType    BoxType
	Node       *StyledNode
	Children   []LayoutBox
}

// NewLayoutBox is a constructor for a LayoutBox with a certain box type
func NewLayoutBox(boxType BoxType, node *StyledNode) LayoutBox {
	return LayoutBox{
		BoxType:  boxType,
		Node:     node,
		Children: make([]LayoutBox, 0),
	}
}

// GetStyledNode returns the StyledNode associated with a LayoutBox
func (lb LayoutBox) GetStyledNode() StyledNode {
	node := lb.Node
	if node == nil {
		return StyledNode{}
	}

	return *node
}

// Layout determines the dimensions for a LayoutBox based on its container
func (lb *LayoutBox) Layout(container Dimensions) {
	if lb.BoxType == BlockNode {
		lb.LayoutBlock(container)
	}

	// TODO: layout InlineNode
}

// LayoutBlock determines the dimensions of a LayoutBox of BoxType block
func (lb *LayoutBox) LayoutBlock(container Dimensions) {
	fmt.Println("a")
	// Calculate the box's width based on the parent
	lb.CalculateBlockWidth(container)

	// Calculate the box's position based on the parent
	lb.CalculateBlockPosition(container)

	// Perform all position and width calculations on box's children
	lb.LayoutBlockChildren()

	// THEN calculate the height of the box
	lb.CalculateBlockHeight()
}

// CalculateBlockWidth calculates the box's width based on the dimensions of its container
func (lb *LayoutBox) CalculateBlockWidth(container Dimensions) {
	styledNode := lb.GetStyledNode()

	fmt.Printf("%+v\n", styledNode)

	width := "auto"
	styledWidth := styledNode.value("width")
	if styledWidth != nil {
		width = *styledWidth
	}

	zero := "0"

	marginLeft := styledNode.Lookup([]string{"margin-left", "margin"}, zero)
	marginRight := styledNode.Lookup([]string{"margin-right", "margin"}, zero)

	fmt.Println("marginLeft", marginLeft)

	borderLeft := styledNode.Lookup([]string{"border-left-width", "border-left"}, zero)
	borderRight := styledNode.Lookup([]string{"border-right-width", "border-right"}, zero)

	paddingLeft := styledNode.Lookup([]string{"padding-left", "padding"}, zero)
	paddingRight := styledNode.Lookup([]string{"padding-right", "padding"}, zero)

	totalWidth := 0
	for _, partialWidth := range []string{marginLeft, marginRight, borderLeft, borderRight, paddingLeft, paddingRight, width} {
		totalWidth += convertToPixels(partialWidth)
	}

	if width != "auto" && totalWidth > container.Content.Width {
		// overflow! should do something here
	}

	underflow := container.Content.Width - totalWidth
	if width != "auto" && marginLeft != "auto" && marginRight != "auto" {
		marginRight = string(convertToPixels(marginRight) + underflow)
	}

	if width != "auto" && marginLeft != "auto" && marginRight == "auto" {
		marginRight = string(underflow)
	}

	if width != "auto" && marginLeft == "auto" && marginRight != "auto" {
		marginLeft = string(underflow)
	}

	if width == "auto" {
		if marginLeft == "auto" {
			marginLeft = zero
		}
		if marginRight == "auto" {
			marginRight = zero
		}

		if underflow >= 0 {
			width = string(underflow)
		} else {
			width = zero
			marginRight = string(convertToPixels(marginRight) + underflow)
		}
	}

	if width != "auto" && marginLeft == "auto" && marginRight == "auto" {
		marginLeft = string(underflow / 2)
		marginRight = string(underflow / 2)
	}

	lb.Dimensions = Dimensions{
		Content: Rectangle{
			Width: convertToPixels(width),
		},
		Padding: EdgeSizes{
			Left:  convertToPixels(paddingLeft),
			Right: convertToPixels(paddingRight),
		},
		Border: EdgeSizes{
			Left:  convertToPixels(borderLeft),
			Right: convertToPixels(borderRight),
		},
		Margin: EdgeSizes{
			Left:  convertToPixels(marginLeft),
			Right: convertToPixels(marginRight),
		},
	}

	fmt.Printf("%+v\n", lb.Dimensions)
}

func convertToPixels(s string) int {
	if s == "auto" {
		return 0
	}

	// support %

	if strings.Contains(s, "px") {
		s = strings.Trim(s, "px")
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

// Lookup sees if any element in a slice of fields has a corresponding value in a StyledNode
// if not, it returns a default string
func (s StyledNode) Lookup(fields []string, defaultVal string) string {
	for _, field := range fields {
		if s.value(field) != nil {
			return *s.value(field)
		}
	}
	return defaultVal
}

// CalculateBlockPosition calculates the box's position based on the parent
func (lb *LayoutBox) CalculateBlockPosition(container Dimensions) {
	styledNode := lb.GetStyledNode()
	d := lb.Dimensions

	zero := "0"

	d.Margin.Top = convertToPixels(styledNode.Lookup([]string{"margin-top", "margin"}, zero))
	d.Margin.Bottom = convertToPixels(styledNode.Lookup([]string{"margin-bottom", "margin"}, zero))

	d.Border.Top = convertToPixels(styledNode.Lookup([]string{"border-top-width", "border-top"}, zero))
	d.Border.Bottom = convertToPixels(styledNode.Lookup([]string{"border-bottom-width", "border-bottom"}, zero))

	d.Padding.Top = convertToPixels(styledNode.Lookup([]string{"padding-top", "padding"}, zero))
	d.Padding.Bottom = convertToPixels(styledNode.Lookup([]string{"padding-bottom", "padding"}, zero))

	d.Content.X = container.Content.X + d.Margin.Left + d.Border.Left + d.Padding.Left
	d.Content.Y = container.Content.Height + container.Content.Y + d.Margin.Top + d.Border.Top + d.Padding.Top
}

// LayoutBlockChildren performs all position and width calculations on box's childrenr
func (lb *LayoutBox) LayoutBlockChildren() {
	d := lb.Dimensions

	for _, child := range lb.Children {
		child.Layout(d)
		d.Content.Height = d.Content.Height + child.Dimensions.MarginBox().Height
	}
}

// PaddingBox calculates the area covered by the content area plus its padding
func (d Dimensions) PaddingBox() Rectangle {
	return d.Content.ExpandedBy(d.Padding)
}

// BorderBox calculates the area covered by the PaddingBox plus its border
func (d Dimensions) BorderBox() Rectangle {
	return d.PaddingBox().ExpandedBy(d.Border)
}

// MarginBox calculates the area covered by the MarginBox plus its margin
func (d Dimensions) MarginBox() Rectangle {
	return d.BorderBox().ExpandedBy(d.Margin)
}

// ExpandedBy expands a Rectangle by a set of EdgeSizes
func (r Rectangle) ExpandedBy(edge EdgeSizes) Rectangle {
	return Rectangle{
		X:      r.X - edge.Left,
		Y:      r.Y - edge.Right,
		Width:  r.Width + edge.Left + edge.Right,
		Height: r.Height + edge.Top + edge.Bottom,
	}
}

// CalculateBlockHeight calculates the height of a box
func (lb *LayoutBox) CalculateBlockHeight() {
	styledNode := lb.GetStyledNode()
	if styledNode.value("height") != nil {
		lb.Dimensions.Content.Height = convertToPixels(*styledNode.value("height"))
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
			lb.Children = append(lb.Children, NewLayoutBox(AnonymousBlock, nil)) // make it so
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
