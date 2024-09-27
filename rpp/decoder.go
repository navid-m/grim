package rpp

import (
	"fmt"
	"strconv"
	"strings"
)

// Token types
const (
	OPEN    = "<"
	CLOSE   = ">"
	NEWLINE = "\n"
)

// Element represents an element in the RPP file
type Element struct {
	RootFileName string
	Tag          string
	Attrib       []string
	Children     []*Element
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

// Lexer is the structure to tokenize RPP content
type Lexer struct {
	input  string
	tokens []Token
	pos    int
}

// Token represents a lexeme
type Token struct {
	Type  string
	Value string
}

// NewLexer initializes a new lexer for the input
func NewLexer(input string) *Lexer {
	return &Lexer{input: input, tokens: tokenize(input)}
}

// NextToken returns the next token in the lexer
func (l *Lexer) NextToken() Token {
	if l.pos >= len(l.tokens) {
		return Token{Type: "", Value: ""}
	}
	token := l.tokens[l.pos]
	l.pos++
	return token
}

// Tokenize breaks the input into tokens, handling nested tags and attributes
func tokenize(input string) []Token {
	var tokens []Token
	lines := strings.Split(input, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check if line starts with '<' (Opening tag)
		if strings.HasPrefix(line, OPEN) {
			tokens = append(tokens, Token{Type: "OPEN", Value: OPEN})
			line = strings.TrimPrefix(line, OPEN)
		}

		// Check if line ends with '>' (Closing tag)
		if strings.HasSuffix(line, CLOSE) {
			line = strings.TrimSuffix(line, CLOSE)
			tokens = append(tokens, Token{Type: "CLOSE", Value: CLOSE})
		}

		// Handle standalone lines like "TEMPO 70 4 4"
		if line != "" {
			tokens = append(tokens, Token{Type: "STRING", Value: line})
		}

		tokens = append(tokens, Token{Type: "NEWLINE", Value: NEWLINE})
	}
	return tokens
}

// Parser parses tokens into an Element tree
type Parser struct {
	lexer *Lexer
}

// NewParser initializes a new parser
func NewParser(input string) *Parser {
	return &Parser{lexer: NewLexer(input)}
}

// Parse starts the parsing process
func (p *Parser) Parse() (*Element, error) {
	token := p.lexer.NextToken()
	if token.Type != "OPEN" {
		return nil, fmt.Errorf("expected opening token, got %s", token.Type)
	}
	element, err := p.parseElement()
	if err != nil {
		return nil, fmt.Errorf("error parsing element: %v", err)
	}
	return element, nil
}

// Parse some element, supporting attributes and nested children
func (p *Parser) parseElement() (*Element, error) {
	token := p.lexer.NextToken()
	if token.Type != "STRING" {
		return nil, fmt.Errorf("expected STRING token, got %s", token.Type)
	}

	root := &Element{Tag: token.Value, Attrib: []string{}, Children: []*Element{}}

	for {
		token := p.lexer.NextToken()

		switch token.Type {
		case "OPEN":
			// Handle nested child elements
			child, err := p.parseElement()
			if err != nil {
				return nil, err
			}
			root.Children = append(root.Children, child)

		case "CLOSE":
			// Return when encountering a closing tag
			return root, nil

		case "STRING":
			// Handle both attributes and standalone tags like TEMPO
			// We treat it as an attribute if no OPEN/CLOSE follows
			root.Attrib = append(root.Attrib, token.Value)

		case "NEWLINE":
			// Ignore newlines

		default:
			return nil, fmt.Errorf("unexpected token type: %s", token.Type)
		}
	}
}