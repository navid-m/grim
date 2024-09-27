package rpp

import (
	"fmt"
	"strings"
)

const (
	OPEN    = "<"
	CLOSE   = ">"
	NEWLINE = "\n"
)

// Start parsing
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

// Returns the next token in the lexer
func (l *Lexer) NextToken() Token {
	if l.pos >= len(l.tokens) {
		return Token{Type: "", Value: ""}
	}
	token := l.tokens[l.pos]
	l.pos++
	return token
}

// Breaks the input into tokens while handling nested tags and attributes
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

// Parse some element, including attributes and nested children
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
