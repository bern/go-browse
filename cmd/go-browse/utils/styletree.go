package utils

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bern/go-browse/cmd/go-browse/models"
)

// Matches returns true if a given ElementData matches a Selector
func Matches(element models.ElementData, selector models.Selector) bool {
	if selector.TagName != nil {
		if element.TagName != *selector.TagName {
			return false
		}
	}

	if selector.ID != nil {
		if element.ID() == nil || element.ID() != selector.ID {
			return false
		}
	}

	if selector.Classes != nil {
		if element.Classes() == nil {
			return false
		}

		for _, class := range *selector.Classes {
			elementContainsClass := false
			for _, elementClass := range *element.Classes() {
				if class == elementClass {
					elementContainsClass = true
					break
				}
			}
			if !elementContainsClass {
				return false
			}
		}
	}

	return true
}

// MatchRule determines the highest specificity selector through which a rule applies to an element
func MatchRule(element models.ElementData, rule models.Rule) *models.MatchedRule {
	for _, selector := range rule.Selectors {
		if Matches(element, selector) {
			return &models.MatchedRule{
				Rule:        rule,
				Specificity: CalculateSpecificity(selector),
			}
		}
	}
	return nil
}

// MatchingRules determines all rules in a stylesheet that apply to an element
func MatchingRules(element models.ElementData, stylesheet models.Stylesheet) []models.MatchedRule {
	matchingRules := make([]models.MatchedRule, 0)

	for _, rule := range stylesheet.Rules {
		matchedRule := MatchRule(element, rule)

		if matchedRule != nil {
			matchingRules = append(matchingRules, *matchedRule)
		}
	}

	return matchingRules
}

// SpecifiedValues sorts a set of rules that apply to an element by specificity and map them to property values
func SpecifiedValues(element models.ElementData, stylesheet models.Stylesheet) models.PropertyMap {
	rules := MatchingRules(element, stylesheet)
	values := make(map[string]string)

	sort.Sort(ByMatchedRuleSpecificityAscending(rules))
	for _, rule := range rules {
		for _, declaration := range rule.Rule.Declarations {
			values[declaration.Name] = declaration.Value
		}
	}

	return values
}

// StyleTree takes a root node of the DOM and recursively applies a stylesheet to it
func StyleTree(root models.Node, stylesheet models.Stylesheet) models.StyledNode {
	specifiedValues := make(map[string]string, 0)
	if root.NodeType == models.Element && root.Element != nil {
		specifiedValues = SpecifiedValues(*root.Element, stylesheet)
	}

	children := make([]models.StyledNode, 0)
	for _, child := range root.Children {
		children = append(children, StyleTree(child, stylesheet))
	}

	return models.StyledNode{
		Node:            root,
		SpecifiedValues: specifiedValues,
		Children:        children,
	}
}

// PrintStyledNode recurses down a StyledNode, printing all elements and their associated styles
func PrintStyledNode(root models.StyledNode, level int) {
	printedValue := ""
	for i := 0; i < level; i++ {
		printedValue += "  "
	}
	printedValue += "| -- "

	nodeType := root.Node.NodeType
	switch nodeType {
	case models.Text:
		printedValue += fmt.Sprintf("TextNode(\"%s\")", *root.Node.Text)
		break
	case models.Element:
		printedValue += fmt.Sprintf("ElementNode(\"%s\"", root.Node.Element.TagName)
		for name, val := range root.Node.Element.Attributes {
			printedValue += fmt.Sprintf(", %s=\"%s\"", name, val)
		}
		printedValue += ")"
		break
	default:
		printedValue += "I'm not sure how to print this node..."
	}

	printedValue += " ("

	for name, value := range root.SpecifiedValues {
		printedValue += name + ":" + value + " "
	}

	printedValue = strings.TrimRight(printedValue, " ")
	printedValue += ")"

	fmt.Println(printedValue)

	for _, child := range root.Children {
		PrintStyledNode(child, level+1)
	}
}
