package utils

import (
	"fmt"
	"log"
	"strconv"

	"github.com/bern/go-browse/cmd/go-browse/models"
)

// BuildLayoutTree recurses down a LayoutBox node and builds a layout tree
func BuildLayoutTree(styleNode models.StyledNode) models.LayoutBox {
	root := models.LayoutBox{}
	switch styleNode.Display() {
	case models.Block:
		root = models.NewLayoutBox(models.BlockNode, &styleNode)
		break
	case models.Inline:
		root = models.NewLayoutBox(models.InlineNode, &styleNode)
		break
	case models.None:
		log.Fatal("the root node is set to display: none !!!")
	default:
		log.Fatal("couldn't find a box type on the root node")
	}

	for _, child := range styleNode.Children {
		switch child.Display() {
		case models.Block:
			root.Children = append(root.Children, BuildLayoutTree(child))
		case models.Inline:
			children := root.GetInlineContainer().Children
			children = append(children, BuildLayoutTree(child))
			root.Children = children
		case models.None:
			continue
		}
	}

	return root
}

// PrintLayoutBox recurses down a LayoutBox, printing all box types
func PrintLayoutBox(root models.LayoutBox, level int) {
	printedValue := ""
	for i := 0; i < level; i++ {
		printedValue += "  "
	}
	printedValue += "| -- "

	boxType := root.BoxType
	switch boxType {
	case models.BlockNode:
		printedValue += "block"
	case models.InlineNode:
		printedValue += "inline"
	case models.AnonymousBlock:
		printedValue += "anonymous"
	}

	printedValue += " (" +
		"x:" + strconv.Itoa(root.Dimensions.Content.X) +
		" y:" + strconv.Itoa(root.Dimensions.Content.Y) +
		" width:" + strconv.Itoa(root.Dimensions.Content.Width) +
		" height:" + strconv.Itoa(root.Dimensions.Content.Height) +
		")"

	fmt.Println(printedValue)

	for _, child := range root.Children {
		PrintLayoutBox(child, level+1)
	}
}
