package parser

import (
	"fmt"

	"github.com/dawkaka/go-interpreter/ast"
	"github.com/dawkaka/go-interpreter/lexer"
	"github.com/dawkaka/go-interpreter/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX // -X or !X
	CALL   // myFunction(X)
)

type Parser struct {
	l              *lexer.Lexer
	currToken      token.Token
	peekToken      token.Token
	errors         []string
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.NextToken()
	p.NextToken()
	return p
}

func (p *Parser) NextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.currToken.Type != token.EOF {
		stm := p.ParseStatement()
		if stm != nil {
			program.Statements = append(program.Statements, stm)
		}
		p.NextToken()
	}
	return program
}

func (p *Parser) ParseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.ParseLetStatement()
	case token.RETURN:
		return p.ParseReturnStatement()
	default:
		return p.ParseExpressionStatement()
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekTokenError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	if p.peekToken.Type == t {
		return true
	}
	p.peekTokenError(t)
	return false
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.NextToken()
		return true
	}
	return false
}

func (p *Parser) registerPrefixFn(t token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[t] = fn
}

func (p *Parser) registerInfixFn(t token.TokenType, fn infixParseFn) {
	p.infixParseFns[t] = fn
}

func (p *Parser) ParseLetStatement() ast.Statement {
	ltStm := &ast.LetStatement{Token: p.currToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	id := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	ltStm.Name = id
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	for !p.curTokenIs(token.SEMICOLON) {
		p.NextToken()
	}
	return ltStm
}

func (p *Parser) ParseReturnStatement() ast.Statement {
	rs := &ast.ReturnStatement{Token: p.currToken}
	p.NextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.NextToken()
	}
	return rs
}

func (p *Parser) ParseExpressionStatement() ast.Statement {
	stm := &ast.ExpressionStatement{Tokken: p.currToken}
	//stm.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return stm
}
