package utils

import (
	"log"
	"strconv"
	"strings"
)

func (p *HTMLParser) expectedStringError(expected ...string) {
	if len(expected) < 1 {
		return
	}

	expectedStrings := strings.Join(expected, " or ")

	log.Fatal(
		"error in parsing file: ", p.FilePath,
		", expected string ", expectedStrings,
		" at position ", strconv.Itoa(p.Parser.Pos),
	)
}
