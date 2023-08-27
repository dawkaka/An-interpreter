package parser

import (
	"fmt"
	"testing"

	"github.com/dawkaka/go-interpreter/ast"
	"github.com/dawkaka/go-interpreter/lexer"
	"github.com/dawkaka/go-interpreter/token"
)

func TestLetStatement(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
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

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {

		p := New(lexer.New(tt.input))
		program := p.ParseProgram()
		checkParsedErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}
		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParsedErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	expr, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, expr.Left, left) {
		return false
	}
	if expr.Operator != operator {
		return false
	}
	if !testLiteralExpression(t, expr.Right, right) {
		return false
	}
	return true
}

func TestSomeShit(t *testing.T) {
	stmt := &ast.InfixExpression{
		Left:     &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "5"}, Value: 5},
		Operator: "+",
		Right:    &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "10"}, Value: 10},
	}
	stmt2 := &ast.InfixExpression{
		Left:     &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "alice"}, Value: "alice"},
		Operator: "*",
		Right:    &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "bob"}, Value: "bob"},
	}
	if !testInfixExpression(t, stmt, 5, "+", 10) {
		t.Error("Failed stmt test")
	}
	if !testInfixExpression(t, stmt2, "alice", "*", "bob") {
		t.Error("testInfixExpression failed for stmt2")
	}
}

func TestBooleanExpression(t *testing.T) {
	input := `false;`
	p := New(lexer.New(input))
	program := p.ParseProgram()
	checkParsedErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program has not enough statements. got=%d", len(program.Statements))
		return
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("statement is not ast.ExpressionStatement. got=%T", program.Statements[0])
		return
	}
	boo, ok := stmt.Expression.(*ast.Boolean)
	if !ok {
		t.Errorf("stmt.Expression is not ast.Boolean. got=%T", stmt.Expression)
		return
	}
	if boo.Value != false {
		t.Errorf("boo.Value not set to 'false'. got=%t", boo.Value)
	}
	if boo.TokenLiteral() != "false" {
		t.Errorf("boo.TokenLiteral not set to 'False'. got=%s", boo.TokenLiteral())
	}
}
