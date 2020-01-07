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
	// createBasicParser()
	parseHTMLFile("test_files/index1.html")
	parseCSSFile("test_files/styles1.css")
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
}

func parseHTMLFile(path string) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("failed to open ", path)
	}
	parentNode := utils.ParseHTML(path, string(dat))
	utils.PrintNode(parentNode, 0)
}

func parseCSSFile(path string) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("failed to open ", path)
	}
	stylesheet := utils.ParseCSS(path, string(dat))
	utils.PrintStylesheet(stylesheet)
}
