package lexer

import (
	"testing"

	"github.com/dawkaka/go-interpreter/token"
)

func TextNextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{expectedType: token.ASSIGN, expectedLiteral: "="},
		{expectedType: token.PLUS, expectedLiteral: "+"},
		{expectedType: token.LBRACE, expectedLiteral: "("},
		{expectedType: token.RBRACE, expectedLiteral: ")"},
		{expectedType: token.LPAREN, expectedLiteral: "{"},
		{expectedType: token.RPAREN, expectedLiteral: "}"},
		{expectedType: token.COMMA, expectedLiteral: ","},
		{expectedType: token.SEMICOLON, expectedLiteral: ";"},
		{expectedType: token.EOF, expectedLiteral: ""},
	}

	l := New(input)

	for i, test := range tests {
		tok := l.NextToken()
		if tok.Type != test.expectedType {
			t.Fatalf("tests:[%d] wrong type; expected:[%q] but got: [%q]", i, tok.Type, test.expectedType)
		}
		if tok.Literal != test.expectedLiteral {
			t.Fatalf("tests:[%d] wrong literal; expected:[%q] but got: [%q]", i, tok.Literal, test.expectedLiteral)
		}
	}
}
