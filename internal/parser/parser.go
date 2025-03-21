package parser

import (
	"fishy/internal/lexer"
	"fishy/pkg/ast"
	"fishy/pkg/datatype"
	"fishy/pkg/token"
	"fishy/pkg/utils"
	"fmt"
	"log"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) Parse() ([]ast.Statement, error) {
	var statements []ast.Statement

	for p.currentToken.Kind != token.EOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}

	return statements, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.currentToken.Kind {
	case token.LABEL:
		return p.parseLabel()
	case token.IDENTIFIER:
		return p.parseIdentifier()
	case token.INSTRUCTION:
		return p.parseInstruction()
	case token.SEQUENCE:
		return p.parseSequence()
	case token.COMMENT:
		p.nextToken()
		return p.parseStatement()
	default:
		return nil, fmt.Errorf("unexpected token: %v", p.currentToken.Value)
	}
}

func (p *Parser) parseLabel() (ast.Statement, error) {
	name := p.currentToken.Value
	p.nextToken()
	return &ast.Label{Name: name}, nil
}

func (p *Parser) parseIdentifier() (ast.Statement, error) {
	id := p.currentToken.Value
	p.nextToken()
	return &ast.Identifier{Value: id}, nil
}

func (p *Parser) parseInstruction() (ast.Statement, error) {
	instructionName := p.currentToken.Value
	instruction := &ast.Instruction{Name: instructionName}
	p.nextToken()

	if p.currentToken.Kind == token.DATA_TYPE {
		instruction.DataType = datatype.FromString(p.currentToken.Value)
		p.nextToken()
	} else {
		instruction.DataType = datatype.UNSET
	}

	var args []ast.Value
	for p.currentToken.Kind != token.EOF && p.currentToken.Kind != token.COMMA && p.currentToken.Kind != token.INSTRUCTION && p.currentToken.Kind != token.SEQUENCE && p.currentToken.Kind != token.LABEL {
		arg, err := p.parseArgument()
		if err != nil {
			return nil, err
		}
		if arg == nil {
			break
		}
		args = append(args, arg)

		if p.currentToken.Kind != token.COMMA {
			break
		}
		p.nextToken()
	}

	instruction.Args = args

	return instruction, nil
}

func (p *Parser) parseSequence() (ast.Statement, error) {
	sequenceName := p.currentToken.Value
	p.nextToken()

	var values []ast.Value
	for p.currentToken.Kind != token.EOF && p.currentToken.Kind != token.COMMA && p.currentToken.Kind != token.INSTRUCTION && p.currentToken.Kind != token.SEQUENCE && p.currentToken.Kind != token.LABEL {
		arg, err := p.parseArgument()
		if err != nil {
			return nil, err
		}
		values = append(values, arg)

		if p.currentToken.Kind != token.COMMA {
			break
		}
		p.nextToken()
	}

	return &ast.Sequence{Name: sequenceName, Values: values}, nil
}

func (p *Parser) parseArgument() (ast.Value, error) {
	switch p.currentToken.Kind {
	case token.IDENTIFIER:
		id := p.currentToken.Value
		p.nextToken()
		return &ast.Identifier{Value: id}, nil
	case token.REGISTER:
		name := p.currentToken.Value
		p.nextToken()
		return &ast.Register{Value: utils.RegisterToIndex(name)}, nil
	case token.IMMEDIATE:
		lit := p.currentToken.Value
		p.nextToken()
		return &ast.NumberLiteral{Value: lit}, nil
	case token.STRING:
		lit := p.currentToken.Value
		p.nextToken()
		return &ast.StringLiteral{Value: lit}, nil
	case token.LEFT_BRACKET:
		p.nextToken()
		address, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if p.currentToken.Kind != token.RIGHT_BRACKET {
			return nil, fmt.Errorf("expected ']', got: %v", p.currentToken)
		}
		p.nextToken()
		return &ast.AddressOf{Value: address}, nil
	case token.COMMENT:
		p.nextToken()
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected token in argument: %v", p.currentToken)
	}
}

func (p *Parser) parseExpression() (ast.Value, error) {
	leftExpr, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	var operator ast.Operator

	if p.currentToken.Kind == token.PLUS || p.currentToken.Kind == token.MINUS || p.currentToken.Kind == token.STAR || p.currentToken.Kind == token.SLASH {
		operator = ast.OperatorFromTokenKind(p.currentToken.Kind)
		p.nextToken()
	} else {
		return leftExpr, nil
	}

	rightExpr, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	switch left := leftExpr.(type) {
	case *ast.Identifier:
		switch right := rightExpr.(type) {
		case *ast.NumberLiteral:
			return &ast.LabelOffsetNumber{
				Left:     left,
				Operator: operator,
				Right:    *right,
			}, nil
		case *ast.Register:
			return &ast.LabelOffsetRegister{
				Left:     left,
				Operator: operator,
				Right:    *right,
			}, nil
		default:
			log.Fatal("unknown expression", "left", left.String(), "op", operator.String(), "right", right.String())
		}
	case *ast.Register:
		switch right := rightExpr.(type) {
		case *ast.NumberLiteral:
			return &ast.RegisterOffsetNumber{
				Left:     *left,
				Operator: operator,
				Right:    *right,
			}, nil
		case *ast.Register:
			return &ast.RegisterOffsetRegister{
				Left:     *left,
				Operator: operator,
				Right:    *right,
			}, nil
		default:
			log.Fatal("unknown expression", "left", left.String(), "op", operator.String(), "right", right.String())
		}
	default:
		log.Fatal("unknown expression", "left", leftExpr.String(), "op", operator.String(), "right", rightExpr.String())
	}

	return leftExpr, nil
}

func (p *Parser) parseTerm() (ast.Value, error) {
	switch p.currentToken.Kind {
	case token.REGISTER:
		name := p.currentToken.Value
		p.nextToken()
		return &ast.Register{Value: utils.RegisterToIndex(name)}, nil
	case token.IMMEDIATE:
		lit := p.currentToken.Value
		p.nextToken()
		return &ast.NumberLiteral{Value: lit}, nil
	case token.IDENTIFIER:
		ident := p.currentToken.Value
		p.nextToken()
		return &ast.Identifier{Value: ident}, nil
	default:
		return nil, fmt.Errorf("unexpected token in term: %v", p.currentToken)
	}
}
