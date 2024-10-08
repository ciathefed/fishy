package ast

import (
	"fishy/pkg/datatype"
	"fishy/pkg/log"
	"fishy/pkg/token"
)

type Statement interface {
	String() string
}

type Label struct {
	Name string
}

type Instruction struct {
	Name     string
	DataType datatype.DataType
	Args     []Value
}

type Sequence struct {
	Name   string
	Values []Value
}

type RegisterOffset struct {
	Left     Register
	Operator Operator
	Right    NumberLiteral
}

type LabelOffset struct {
	Left     Value
	Operator Operator
	Right    NumberLiteral
}

func (l *Label) String() string          { return "LABEL" }
func (i *Instruction) String() string    { return "INSTRUCTION" }
func (s *Sequence) String() string       { return "SEQUENCE" }
func (b *RegisterOffset) String() string { return "REGISTER_OFFSET" }
func (b *LabelOffset) String() string    { return "LABEL_OFFSET" }

type Operator int

const (
	ADD Operator = iota
	SUBTRACT
	MULTIPLY
	DIVIDE
)

func (o Operator) String() string {
	switch o {
	case ADD:
		return "+"
	case SUBTRACT:
		return "-"
	case MULTIPLY:
		return "*"
	case DIVIDE:
		return "/"
	default:
		log.Fatal("unknown operator", "kind", int(o))
		return ""
	}
}

func OperatorFromTokenKind(kind token.TokenKind) Operator {
	switch kind {
	case token.PLUS:
		return ADD
	case token.MINUS:
		return SUBTRACT
	case token.STAR:
		return MULTIPLY
	case token.SLASH:
		return DIVIDE
	default:
		log.Fatal("unknown operator", "kind", kind)
		return -1
	}
}
