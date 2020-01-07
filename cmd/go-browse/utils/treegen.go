package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bern/go-browse/cmd/go-browse/models"
)

// TextNode func creates a Node out of a given string
func TextNode(data string) models.Node {
	return models.Node{
		Children: make([]models.Node, 0),
		NodeType: models.Text,
		Text:     &data,
	}
}

// ElementNode func creates a Node out of a given element
func ElementNode(name string, attrs map[string]string, children []models.Node) models.Node {
	return models.Node{
		Children: children,
		NodeType: models.Element,
		Element: &models.ElementData{
			TagName:    name,
			Attributes: attrs,
		},
	}
}

// PrintNode will recurse down a node and print all of its descendants
func PrintNode(root models.Node, level int) {
	printedValue := ""
	for i := 0; i < level; i++ {
		printedValue += "\t"
	}
	printedValue += "| -- "

	nodeType := root.NodeType
	switch nodeType {
	case models.Text:
		printedValue += fmt.Sprintf("TextNode(\"%s\")", *root.Text)
		break
	case models.Element:
		printedValue += fmt.Sprintf("ElementNode(\"%s\"", root.Element.TagName)
		for name, val := range root.Element.Attributes {
			printedValue += fmt.Sprintf(", %s=\"%s\"", name, val)
		}
		printedValue += ")"
		break
	default:
		printedValue += "I'm not sure how to print this node..."
	}

	fmt.Println(printedValue)

	for _, child := range root.Children {
		PrintNode(child, level+1)
	}
}

// PrintStylesheet will crawl down a stylesheet and print every rule
func PrintStylesheet(sheet models.Stylesheet) {
	for i, rule := range sheet.Rules {
		PrintRule(rule, i)
	}
}

// PrintRule takes a rule and prints outs its selectors and declarations
func PrintRule(rule models.Rule, index int) {
	fmt.Printf("**Rule #%s**\n", strconv.Itoa(index))
	fmt.Println("Selectors")
	for _, selector := range rule.Selectors {
		if selector.TagName != nil {
			fmt.Println("Type:", *selector.TagName)
		} else {
			fmt.Println("Type: universal")
		}

		if selector.ID != nil {
			fmt.Println("ID:", *selector.ID)
		}

		if selector.Classes != nil {
			fmt.Println("Classes:", strings.Join(*selector.Classes, ", "))
		}

		fmt.Println()
	}
	fmt.Println("Declarations")
	for _, declaration := range rule.Declarations {
		fmt.Println("Name:", declaration.Name)
		fmt.Println("Value:", declaration.Value)
		fmt.Println()
	}
}
