package lexer_test

import (
	"fishy/internal/lexer"
	"fishy/pkg/token"
	"testing"
)

type InputExpectedTest struct {
	input    string
	expected []token.Token
}

func TestImmediate(t *testing.T) {
	tests := []InputExpectedTest{
		{
			input: "10",
			expected: []token.Token{
				{Kind: token.IMMEDIATE, Value: "10"},
			},
		},
		{
			input: "100",
			expected: []token.Token{
				{Kind: token.IMMEDIATE, Value: "100"},
			},
		},
		{
			input: "0x10",
			expected: []token.Token{
				{Kind: token.IMMEDIATE, Value: "16"},
			},
		},
		{
			input: "0x100",
			expected: []token.Token{
				{Kind: token.IMMEDIATE, Value: "256"},
			},
		},
		{
			input: "0o10",
			expected: []token.Token{
				{Kind: token.IMMEDIATE, Value: "8"},
			},
		},
		{
			input: "0o100",
			expected: []token.Token{
				{Kind: token.IMMEDIATE, Value: "64"},
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		for _, expectedToken := range tt.expected {
			tok := l.NextToken()
			if tok.Kind != expectedToken.Kind || tok.Value != expectedToken.Value {
				t.Fatalf("expected token %v, got %v", expectedToken, tok)
			}
		}
	}
}

func TestInstruction(t *testing.T) {
	tests := []InputExpectedTest{
		{
			input: "mov x0, 10",
			expected: []token.Token{
				{Kind: token.INSTRUCTION, Value: "mov"},
				{Kind: token.REGISTER, Value: "x0"},
				{Kind: token.COMMA, Value: ","},
				{Kind: token.IMMEDIATE, Value: "10"},
			},
		},
		{
			input: "mov x1, x2",
			expected: []token.Token{
				{Kind: token.INSTRUCTION, Value: "mov"},
				{Kind: token.REGISTER, Value: "x1"},
				{Kind: token.COMMA, Value: ","},
				{Kind: token.REGISTER, Value: "x2"},
			},
		},
		{
			input: "mov x3, [x4]",
			expected: []token.Token{
				{Kind: token.INSTRUCTION, Value: "mov"},
				{Kind: token.REGISTER, Value: "x3"},
				{Kind: token.COMMA, Value: ","},
				{Kind: token.LEFT_BRACKET, Value: "["},
				{Kind: token.REGISTER, Value: "x4"},
				{Kind: token.RIGHT_BRACKET, Value: "]"},
			},
		},
		{
			input: "mov x5, [x6 + 0x10]",
			expected: []token.Token{
				{Kind: token.INSTRUCTION, Value: "mov"},
				{Kind: token.REGISTER, Value: "x5"},
				{Kind: token.COMMA, Value: ","},
				{Kind: token.LEFT_BRACKET, Value: "["},
				{Kind: token.REGISTER, Value: "x6"},
				{Kind: token.PLUS, Value: "+"},
				{Kind: token.IMMEDIATE, Value: "16"},
				{Kind: token.RIGHT_BRACKET, Value: "]"},
			},
		},
		{
			input: "syscall",
			expected: []token.Token{
				{Kind: token.INSTRUCTION, Value: "syscall"},
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		for _, expectedToken := range tt.expected {
			tok := l.NextToken()
			if tok.Kind != expectedToken.Kind || tok.Value != expectedToken.Value {
				t.Fatalf("expected token %v, got %v", expectedToken, tok)
			}
		}
	}
}

func TestSequence(t *testing.T) {
	tests := []InputExpectedTest{
		{
			input: "db \"Hello, World!\n\", 0",
			expected: []token.Token{
				{Kind: token.SEQUENCE, Value: "db"},
				{Kind: token.STRING, Value: "Hello, World!\n"},
				{Kind: token.COMMA, Value: ","},
				{Kind: token.IMMEDIATE, Value: "0"},
			},
		},
		{
			input: "db 0x10, 0x1c",
			expected: []token.Token{
				{Kind: token.SEQUENCE, Value: "db"},
				{Kind: token.IMMEDIATE, Value: "16"},
				{Kind: token.COMMA, Value: ","},
				{Kind: token.IMMEDIATE, Value: "28"},
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		for _, expectedToken := range tt.expected {
			tok := l.NextToken()
			if tok.Kind != expectedToken.Kind || tok.Value != expectedToken.Value {
				t.Fatalf("expected token %v, got %v", expectedToken, tok)
			}
		}
	}
}
