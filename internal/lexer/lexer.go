package lexer

import (
	"fishy/pkg/token"
	"fishy/pkg/utils"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	start := l.position

	switch l.ch {
	case ',':
		l.readChar()
		return token.Token{Kind: token.COMMA, Value: ",", Start: start, End: l.position}
	case ':':
		l.readChar()
		return token.Token{Kind: token.COLON, Value: ":", Start: start, End: l.position}
	case '[':
		l.readChar()
		return token.Token{Kind: token.LEFT_BRACKET, Value: "[", Start: start, End: l.position}
	case ']':
		l.readChar()
		return token.Token{Kind: token.RIGHT_BRACKET, Value: "]", Start: start, End: l.position}
	case '+':
		l.readChar()
		return token.Token{Kind: token.PLUS, Value: "+", Start: start, End: l.position}
	case '-':
		l.readChar()
		return token.Token{Kind: token.MINUS, Value: "-", Start: start, End: l.position}
	case '*':
		l.readChar()
		return token.Token{Kind: token.STAR, Value: "*", Start: start, End: l.position}
	case '/':
		l.readChar()
		return token.Token{Kind: token.SLASH, Value: "/", Start: start, End: l.position}
	case '"':
		return l.readString()
	case '$':
		return l.readImmediate()
	case '#':
		return l.readDirective()
	case ';':
		l.skipComment()
		return token.Token{Kind: token.COMMENT, Value: "", Start: start, End: l.position}
	case 0:
		return token.Token{Kind: token.EOF, Value: "", Start: start, End: start}
	}

	if isLetter(l.ch) || l.ch == '_' || l.ch == '.' {
		ident := l.readIdentifier()
		if l.ch == ':' {
			l.readChar()
			return token.Token{Kind: token.LABEL, Value: ident, Start: start, End: l.position}
		}
		end := l.position
		if isInstruction(ident) {
			return token.Token{Kind: token.INSTRUCTION, Value: ident, Start: start, End: end}
		}
		if isSequence(ident) {
			return token.Token{Kind: token.SEQUENCE, Value: ident, Start: start, End: end}
		}
		if isRegister(ident) {
			return token.Token{Kind: token.REGISTER, Value: ident, Start: start, End: end}
		}
		return token.Token{Kind: token.IDENTIFIER, Value: ident, Start: start, End: end}
	}

	if isDigit(l.ch) {
		return l.readImmediate()
	}

	l.readChar()
	return token.Token{Kind: token.EOF, Value: "", Start: start, End: start}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' || l.ch == '.' {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readString() token.Token {
	var str strings.Builder
	start := l.position
	l.readChar()

	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				str.WriteByte('\n')
			case 't':
				str.WriteByte('\t')
			case 'x':
				hexValue := l.readHexSequence(2)
				if hexValue != -1 {
					str.WriteByte(byte(hexValue))
				}
			default:
				str.WriteByte(l.ch)
			}
		} else {
			str.WriteByte(l.ch)
		}
		l.readChar()
	}

	l.readChar()
	return token.Token{Kind: token.STRING, Value: str.String(), Start: start, End: l.position}
}

func (l *Lexer) readImmediate() token.Token {
	start := l.position
	isHex := false
	isOct := false

	if l.ch == '$' {
		l.readChar()
	}

	if l.ch == '0' && (l.peekChar() == 'x' || l.peekChar() == 'X') {
		isHex = true
		l.readChar()
		l.readChar()
		for isHexDigit(l.ch) {
			l.readChar()
		}
	} else if l.ch == '0' && (l.peekChar() == 'o' || l.peekChar() == 'O') {
		isOct = true
		l.readChar()
		l.readChar()
		for isOctDigit(l.ch) {
			l.readChar()
		}
	} else {
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	value := l.input[start:l.position]
	value = strings.TrimPrefix(value, "$")

	if isHex {
		num, err := strconv.ParseUint(value[2:], 16, 64)
		if err == nil {
			return token.Token{Kind: token.IMMEDIATE, Value: fmt.Sprintf("%d", num), Start: start, End: l.position}
		}
	} else if isOct {
		num, err := strconv.ParseInt(value[2:], 8, 64)
		if err == nil {
			return token.Token{Kind: token.IMMEDIATE, Value: fmt.Sprintf("%d", num), Start: start, End: l.position}
		}
	} else {
		num, err := strconv.Atoi(value)
		if err == nil {
			return token.Token{Kind: token.IMMEDIATE, Value: fmt.Sprintf("%d", num), Start: start, End: l.position}
		}
	}

	return token.Token{Kind: token.IMMEDIATE, Value: value, Start: start, End: l.position}
}

func (l *Lexer) readDirective() token.Token {
	start := l.position
	l.readChar()
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	value := l.input[start:l.position]
	value = strings.TrimPrefix(value, "#")
	return token.Token{Kind: token.DIRECTIVE, Value: value, Start: start, End: l.position}
}

func (l *Lexer) readHexSequence(length int) int {
	var hexStr strings.Builder
	for i := 0; i < length && isHexDigit(l.ch); i++ {
		hexStr.WriteByte(l.ch)
		l.readChar()
	}
	hexValue, err := strconv.ParseUint(hexStr.String(), 16, 64)
	if err != nil {
		return -1
	}
	return int(hexValue)
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func isInstruction(ident string) bool {
	for _, instr := range utils.Instructions {
		if instr == ident {
			return true
		}
	}
	return false
}

func isSequence(ident string) bool {
	for _, seq := range utils.Sequences {
		if seq == ident {
			return true
		}
	}
	return false
}

func isRegister(ident string) bool {
	for _, reg := range utils.Registers {
		if reg == ident {
			return true
		}
	}
	return false
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isHexDigit(ch byte) bool {
	return ('0' <= ch && ch <= '9') || ('a' <= ch && ch <= 'f') || ('A' <= ch && ch <= 'F')
}

func isOctDigit(ch byte) bool {
	return ch >= '0' && ch <= '7'
}
