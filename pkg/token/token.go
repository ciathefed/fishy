package token

type TokenKind string

const (
	DIRECTIVE     TokenKind = "DIRECTIVE"
	INSTRUCTION   TokenKind = "INSTRUCTION"
	SEQUENCE      TokenKind = "SEQUENCE"
	LABEL         TokenKind = "LABEL"
	IDENTIFIER    TokenKind = "IDENTIFIER"
	IMMEDIATE     TokenKind = "IMMEDIATE"
	REGISTER      TokenKind = "REGISTER"
	STRING        TokenKind = "STRING"
	COMMENT       TokenKind = "COMMENT"
	COMMA         TokenKind = "COMMA"
	COLON         TokenKind = "COLON"
	LEFT_BRACKET  TokenKind = "LEFT_BRACKET"
	RIGHT_BRACKET TokenKind = "RIGHT_BRACKET"
	PLUS          TokenKind = "PLUS"
	MINUS         TokenKind = "MINUS"
	STAR          TokenKind = "STAR"
	SLASH         TokenKind = "SLASH"
	EOF           TokenKind = "EOF"
)

type Token struct {
	Kind  TokenKind
	Value string
	Start int
	End   int
}
