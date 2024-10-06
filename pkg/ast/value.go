package ast

type Value interface {
	String() string
	Index() int
}

type NumberLiteral struct {
	Value string
}

type StringLiteral struct {
	Value string
}

type Register struct {
	Value int
}

type AddressOf struct {
	Value Value
}

type Identifier struct {
	Value string
}

func (n *NumberLiteral) String() string { return "NUMBER_LITERAL" }
func (s *StringLiteral) String() string { return "STRING_LITERAL" }
func (r *Register) String() string      { return "REGISTER" }
func (a *AddressOf) String() string     { return "ADDRESS_OF" }
func (i *Identifier) String() string    { return "IDENTIFIER" }

func (n *NumberLiteral) Index() int    { return 0 }
func (s *StringLiteral) Index() int    { return 1 }
func (r *Register) Index() int         { return 2 }
func (a *AddressOf) Index() int        { return 3 }
func (i *Identifier) Index() int       { return 4 }
func (b *BinaryExpression) Index() int { return 5 }
