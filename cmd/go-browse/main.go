package main

import (
	"github.com/bern/go-browse/cmd/go-browse/models"
	"github.com/bern/go-browse/cmd/go-browse/utils"
)

func main() {
	// fmt.Println("Hello, World!")
	// tests.DrawSomething()
	createBasicTree()
}

func createBasicTree() {
	helloNode := utils.TextNode("Hello, World!")
	goodbyeNode := utils.TextNode("Goodbye, World!")
	divChildren := []models.Node{goodbyeNode}
	divNode := utils.ElementNode("div", make(map[string]string, 0), divChildren)
	rootChildren := []models.Node{helloNode, divNode}
	rootNode := utils.ElementNode("html", make(map[string]string, 0), rootChildren)
	// fmt.Printf("%+v\n", parentNode)
	utils.PrintNode(rootNode, 0)
}
