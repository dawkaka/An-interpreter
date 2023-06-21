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

func (l *Lexer) readIdentifier() string {
	pos := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == '\r' || l.ch == ' ' || l.ch == '\n' || l.ch == '\t' {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	pos := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhiteSpace()
	c := l.ch
	var tok token.Token
	switch c {
	case '=':
		tok = AssignToken(token.ASSIGN, c)
	case '+':
		tok = AssignToken(token.PLUS, c)
	case '(':
		tok = AssignToken(token.LPAREN, c)
	case ')':
		tok = AssignToken(token.RPAREN, c)
	case '{':
		tok = AssignToken(token.LBRACE, c)
	case '}':
		tok = AssignToken(token.RBRACE, c)
	case ',':
		tok = AssignToken(token.COMMA, c)
	case ';':
		tok = AssignToken(token.SEMICOLON, c)
	case '!':
		tok = AssignToken(token.BANG, c)
	case '-':
		tok = AssignToken(token.MINUS, c)
	case '/':
		tok = AssignToken(token.SLASH, c)
	case '<':
		tok = AssignToken(token.LT, c)
	case '>':
		tok = AssignToken(token.GT, c)
	case '*':
		tok = AssignToken(token.ASTERISK, c)

	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = AssignToken(token.ILLEGAL, l.ch)
		}
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

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
