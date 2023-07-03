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

func TestReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if program == nil {
		t.Fatal("Test Return statements: Program is nil")
	}
	checkParsedErrors(t, p)

	if len(program.Statements) < 3 {
		t.Fatalf("Expected length to be 3 but got %d", len(program.Statements))
	}
	for _, st := range program.Statements {
		rtStm, ok := st.(*ast.ReturnStatement)
		if !ok {
			t.Fatal("Not a return statement")
		}
		if rtStm.TokenLiteral() != "return" {
			t.Errorf("ReturnStatement.TokenLiteral not 'return', got %q", rtStm.TokenLiteral())
		}

	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParsedErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expected stmt.Expression to be *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			ident.TokenLiteral())
	}
}
