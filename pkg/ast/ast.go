package ast

import "fishy/pkg/token"

type Statement interface {
	String() string
}

type Program struct {
	Statements []Statement
}

type Label struct {
	Name string
}

type Instruction struct {
	Name string
	Args []Value
}

type Sequence struct {
	Name   string
	Values []Value
}

type BinaryExpression struct {
	Left     Value
	Operator token.TokenKind
	Right    Value
}

func (p *Program) String() string          { return "PROGRAM" }
func (l *Label) String() string            { return "LABEL" }
func (i *Instruction) String() string      { return "INSTRUCTION" }
func (s *Sequence) String() string         { return "SEQUENCE" }
func (b *BinaryExpression) String() string { return "BINARY_EXPRESSION" }
