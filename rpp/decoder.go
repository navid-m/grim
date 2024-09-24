package rpp

import (
	"fmt"
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
	Tag      string
	Attrib   []string
	Children []*Element
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

// Tokenize breaks the input into tokens
func tokenize(input string) []Token {
	var tokens []Token
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, OPEN) {
			tokens = append(tokens, Token{Type: "OPEN", Value: OPEN})
			line = strings.TrimPrefix(line, OPEN)
		}
		if strings.HasSuffix(line, CLOSE) {
			line = strings.TrimSuffix(line, CLOSE)
			tokens = append(tokens, Token{Type: "CLOSE", Value: CLOSE})
		}
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
		return nil, err
	}
	return element, nil
}

// Parse some element
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
			child, err := p.parseElement()
			if err != nil {
				return nil, err
			}
			root.Children = append(root.Children, child)
		case "CLOSE":
			return root, nil
		case "STRING":
			root.Attrib = append(root.Attrib, token.Value)
		case "NEWLINE":
		default:
			return nil, fmt.Errorf("unexpected token type: %s", token.Type)
		}
	}
}
