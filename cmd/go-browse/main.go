package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/bern/go-browse/cmd/go-browse/models"
	"github.com/bern/go-browse/cmd/go-browse/utils"
)

func main() {
	// tests.DrawSomething()
	// createBasicTree()
	createBasicParser()
}

func createBasicTree() {
	helloNode := utils.TextNode("Hello, World!")
	goodbyeNode := utils.TextNode("Goodbye, World!")
	divChildren := []models.Node{goodbyeNode}
	divNode := utils.ElementNode("div", make(map[string]string, 0), divChildren)
	rootChildren := []models.Node{helloNode, divNode}
	rootNode := utils.ElementNode("html", make(map[string]string, 0), rootChildren)
	utils.PrintNode(rootNode, 0)
}

func createBasicParser() {
	parser := utils.Parser{
		Input: "hello, world!",
	}
	fmt.Println(parser.NextChar())
	prefix := "hello,"
	fmt.Printf(
		"does the next bit of the parser start with prefix \"%s\"? %+v\n",
		prefix,
		parser.StartsWith(prefix),
	)

	parser.Input = "c"
	parser.Pos = 1
	fmt.Println(parser.EOF()) // should be true

	parser.Input = "a           b"
	parser.Pos = 0
	fmt.Println(parser.ConsumeChar())
	parser.ConsumeWhitespace()
	fmt.Println(parser.ConsumeChar())

	parser.Input = "abc123ABC/>"
	parser.Pos = 0
	fmt.Println(parser.ConsumeName())

	parser.Input = "hello</>"
	parser.Pos = 0
	helloNode := parser.ParseNode()
	utils.PrintNode(helloNode, 0)

	parser.Input = "<div height=\"100\" width=\"100\">"
	parser.Pos = 0
	divNode := parser.ParseNode()
	utils.PrintNode(divNode, 0)

	filePath := "test1.html"
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("failed to open ", filePath)
	}
	parentNode := utils.ParseHTML(string(dat))
	utils.PrintNode(parentNode, 0)
}
