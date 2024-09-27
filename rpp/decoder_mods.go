package rpp

import (
	"fmt"
	"strconv"
)

// An element in the RPP file
type Element struct {
	RootFileName string
	Tag          string
	Attrib       []string
	Children     []*Element
}

// Lexer to tokenize RPP content
type Lexer struct {
	input  string
	tokens []Token
	pos    int
}

// Some lexeme
type Token struct {
	Type  string
	Value string
}

// Currently a lexer container
type Parser struct {
	lexer *Lexer
}

// Initialize a parser
func NewParser(input string) *Parser {
	return &Parser{lexer: NewLexer(input)}
}

// Create new lexer for the input
func NewLexer(input string) *Lexer {
	return &Lexer{input: input, tokens: tokenize(input)}
}

func (e Element) String() string {
	toret := ""

	if e.RootFileName != "" {
		toret += fmt.Sprintln("Root File Name: ", e.RootFileName)
	}

	toret += fmt.Sprintln("Tag: ", e.Tag)

	for i, attrib := range e.Attrib {
		toret += fmt.Sprintln("\t - Attrib #" + strconv.Itoa(i) + ": " + attrib)
	}

	for i, child := range e.Children {
		toret += fmt.Sprintln("\t - Child #"+strconv.Itoa(i)+": ", *child)
	}

	return toret
}
