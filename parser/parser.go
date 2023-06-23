package parser

import (
	"github.com/dawkaka/go-interpreter/ast"
	"github.com/dawkaka/go-interpreter/lexer"
	"github.com/dawkaka/go-interpreter/token"
)

type Parser struct {
	l         *lexer.Lexer
	currToken token.Token
	peekToken token.Token
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
	default:
		return nil
	}
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
	for !p.currTokenIs(token.SEMICOLON) {
		p.NextToken()
	}
	return ltStm
}

func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.NextToken()
		return true
	}
	return false
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.NextToken()
	p.NextToken()
	return p
}
