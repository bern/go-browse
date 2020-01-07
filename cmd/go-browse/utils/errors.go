package utils

import (
	"log"
	"strconv"
)

func (p *HTMLParser) expectedStringError(expected ...string) {
	if len(expected) < 1 {
		return
	}

	expectedStrings := expected[0]
	for i, expectedString := range expected {
		if i == 0 {
			continue
		}

		expectedStrings += " or " + expectedString
	}

	log.Fatal(
		"error in parsing file: ", p.FilePath,
		", expected string ", expectedStrings,
		" at position ", strconv.Itoa(p.Parser.Pos),
	)
}
