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

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`
	p := New(lexer.New(input))
	checkParsedErrors(t, p)
	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	integer, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected stmt.Expression to be *ast.Interger. got=%T", stmt.Expression)
	}

	if integer.Value != 5 {
		t.Errorf("ident.Value not %s. got=%d", "foobar", integer.Value)
	}
	if integer.TokenLiteral() != "5" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", integer.TokenLiteral())
	}
}

func TestParsingPrefixExpression(t *testing.T) {

	prefixExpressions := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixExpressions {
		p := New(lexer.New(tt.input))
		checkParsedErrors(t, p)
		program := p.ParseProgram()
		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}

}

func testIntegerLiteral(t *testing.T, il ast.Expression, intVal int64) bool {
	in, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if in.Value != intVal {
		t.Errorf("integ.Value not %d. got=%d", intVal, in.Value)
		return false
	}
	if in.TokenLiteral() != fmt.Sprint(intVal) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", intVal,
			in.TokenLiteral())
		return false
	}

	return true
}
