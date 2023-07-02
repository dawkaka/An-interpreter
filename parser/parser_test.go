package parser

import (
	"fmt"
	"testing"

	"github.com/dawkaka/go-interpreter/ast"
	"github.com/dawkaka/go-interpreter/lexer"
)

func TestLetStatement(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatal("Program is nil")
	}
	checkParsedErrors(t, p)
	if len(program.Statements) != 3 {
		t.Fatalf("Expected length to be 3 but got %d", len(program.Statements))
	}

	tests := []struct{ expectedIdentifier string }{{"x"}, {"y"}, {"foobar"}}

	for i, tt := range tests {
		st := program.Statements[i]
		if !testLetStatement(t, st, tt.expectedIdentifier) {
			return
		}

	}
}

func checkParsedErrors(t *testing.T, p *Parser) {
	errors := p.errors
	if len(errors) > 0 {
		t.Errorf("parser has %d errors", len(errors))
		for _, err := range errors {
			fmt.Println(err)
		}
		t.FailNow()
	}
}

func testLetStatement(t *testing.T, st ast.Statement, name string) bool {
	if st.TokenLiteral() != "let" {
		t.Errorf("Expected 'let' but got '%s'", st.TokenLiteral())
		return false
	}

	letStatement, ok := st.(*ast.LetStatement)
	if !ok {
		t.Error("Not a let statement")
		return false
	}
	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("letStatement.Name.TokenLiteral(); expected: %s, got: %s", letStatement.Name.TokenLiteral(), name)
		return false
	}
	if letStatement.Name.Value != name {
		t.Errorf("letStatement.Name.Value not '%s'. got=%s", name, letStatement.Name.Value)
		return false
	}

	return true
}