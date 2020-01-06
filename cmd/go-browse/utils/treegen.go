package utils

import (
	"fmt"

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
		printedValue += fmt.Sprintf("ElementNode(\"%s\")", root.Element.TagName)
		break
	default:
		printedValue += "I'm not sure how to print this node..."
	}

	fmt.Println(printedValue)

	for _, child := range root.Children {
		PrintNode(child, level+1)
	}
}
