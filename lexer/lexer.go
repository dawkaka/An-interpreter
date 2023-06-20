package lexer

import (
	"github.com/dawkaka/go-interpreter/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	c := l.ch
	var tok token.Token
	switch c {
	case '=':
		tok = AssignToken(token.ASSIGN, c)
	case '+':
		tok = AssignToken(token.PLUS, c)
	case '(':
		tok = AssignToken(token.LBRACE, c)
	case ')':
		tok = AssignToken(token.RBRACE, c)
	case '{':
		tok = AssignToken(token.LPAREN, c)
	case '}':
		tok = AssignToken(token.RPAREN, c)
	case ',':
		tok = AssignToken(token.COMMA, c)
	case ';':
		tok = AssignToken(token.SEMICOLON, c)
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	}
	l.readChar()
	return tok
}

func AssignToken(tokenType token.TokenType, literal byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(literal)}
}

func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.readChar()
	return l
}
