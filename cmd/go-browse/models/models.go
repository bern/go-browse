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

type Specificity struct {
	IDSpecificity      int
	ClassSpecificity   int
	ElementSpecificity int
}

// BySpecificity sorts a slice of selectors with the highest specificty first
type BySpecificity []Selector

func (a BySpecificity) Len() int {
	return len(a)
}

func (a BySpecificity) Less(i, j int) bool {
	specI := CalculateSpecificity(a[i])
	specJ := CalculateSpecificity(a[j])

	if specI.IDSpecificity > specJ.IDSpecificity {
		return true
	} else if specI.IDSpecificity < specJ.IDSpecificity {
		return false
	}

	if specI.ClassSpecificity > specJ.ClassSpecificity {
		return true
	} else if specI.ClassSpecificity < specJ.ClassSpecificity {
		return false
	}

	if specI.ElementSpecificity > specJ.ElementSpecificity {
		return true
	} else if specI.ElementSpecificity < specJ.ElementSpecificity {
		return false
	}

	return false
}

func (a BySpecificity) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func CalculateSpecificity(s Selector) Specificity {
	idSpecificity := 0
	classSpecificity := 0
	elementSpecificity := 0

	if s.ID != nil {
		idSpecificity = 1
	}

	if s.Classes != nil {
		classSpecificity = len(*s.Classes)
	}

	if s.TagName != nil {
		elementSpecificity = 1
	}

	return Specificity{
		IDSpecificity:      idSpecificity,
		ClassSpecificity:   classSpecificity,
		ElementSpecificity: elementSpecificity,
	}
}
