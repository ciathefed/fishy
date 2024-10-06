package preprocessor

import (
	"bytes"
	"fishy/internal/lexer"
	"fishy/pkg/token"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Macro struct {
	Name string
	Args []string
	Body string
}

type Preprocessor struct {
	input   []byte
	tokens  []token.Token
	current int
	defines map[string]string
	macros  map[string]Macro
}

func New(input []byte) *Preprocessor {
	l := lexer.New(string(input))
	tokens := []token.Token{}
	for {
		t := l.NextToken()
		tokens = append(tokens, t)
		if t.Kind == token.EOF {
			break
		}
	}

	return &Preprocessor{
		input:   input,
		tokens:  tokens,
		current: 0,
		defines: make(map[string]string),
		macros:  make(map[string]Macro),
	}
}

func (p *Preprocessor) Process() string {
	for p.current < len(p.tokens) {
		t := p.consume()

		if t.Kind == token.COMMENT {
			continue
		}

		if t.Kind != token.DIRECTIVE {
			continue
		}

		switch t.Value {
		case "define":
			p.processDefine(t)
		case "include":
			p.processInclude(t)
		case "macro":
			p.processMacro(t)
		}
	}

	p.expandMacros()
	p.replaceDefinitions()

	return strings.TrimSpace(string(p.input))
}

func (p *Preprocessor) processDefine(t token.Token) {
	identifier := p.consume()
	value := p.consume()
	if value.Kind == token.STRING {
		value.Value = strconv.Quote(value.Value)
	}
	start := t.Start
	end := value.End

	p.defines[identifier.Value] = value.Value

	p.fillWithSpaces(start, end)
}

func (p *Preprocessor) processInclude(t token.Token) {
	fileToken := p.consume()

	filePath := strings.Trim(fileToken.Value, "\"")

	content, err := os.ReadFile(filePath)
	if err != nil {

		panic("failed to include file: " + err.Error())
	}

	start := t.Start
	end := fileToken.End
	p.input = append(p.input[:start], append(content, p.input[end:]...)...)

	p.reTokenize()
}

func (p *Preprocessor) processMacro(t token.Token) {
	start := t.Start
	macroName := p.consume()

	args := []string{}
	for p.peek().Kind == token.IDENTIFIER {
		arg := p.consume()
		args = append(args, arg.Value)
	}

	bodyStart := p.tokens[p.current].Start
	for {
		next := p.consume()
		if next.Kind == token.DIRECTIVE && next.Value == "end" {
			break
		}
	}
	bodyEnd := p.tokens[p.current-2].End

	macroBody := string(p.input[bodyStart:bodyEnd])

	p.macros[macroName.Value] = Macro{
		Name: macroName.Value,
		Args: args,
		Body: macroBody,
	}

	p.fillWithSpaces(start, p.tokens[p.current-1].End)
}

func (p *Preprocessor) expandMacros() {
	for macroName, macro := range p.macros {
		for {
			index := bytes.Index(p.input, []byte(macroName))
			if index == -1 {
				break
			}

			isComment := false
			commentStart := -1
			for i := index - 1; i >= 0; i-- {
				if p.input[i] == '\n' {
					break
				}
				if p.input[i] == ';' || p.input[i] == '/' {
					isComment = true
					commentStart = i
					break
				}
			}

			if isComment {
				endOfLine := bytes.IndexByte(p.input[commentStart:], '\n')
				if endOfLine == -1 {
					endOfLine = len(p.input)
				} else {
					endOfLine += commentStart
				}

				p.input = append(p.input[:commentStart], p.input[endOfLine:]...)
				continue
			}

			newlineIndex := bytes.IndexByte(p.input[index:], '\n')
			if newlineIndex == -1 {
				newlineIndex = len(p.input)
			} else {
				newlineIndex += index
			}

			remainingInput := string(p.input[index+len(macroName) : newlineIndex])
			argValues := []string{}
			var currentArg strings.Builder
			offset := 0

			for i := 0; i < len(remainingInput); i++ {
				char := remainingInput[i]
				offset += 1
				if char == ',' {
					argValues = append(argValues, strings.TrimSpace(currentArg.String()))
					currentArg.Reset()
				} else {
					currentArg.WriteByte(char)
				}
			}

			if currentArg.Len() > 0 {
				argValues = append(argValues, strings.TrimSpace(currentArg.String()))
			}

			if len(argValues) > len(macro.Args) {
				argValues = argValues[:len(macro.Args)]
			}

			macroBody := macro.Body
			for i, arg := range macro.Args {
				macroBody = strings.ReplaceAll(macroBody, arg, argValues[i])
			}

			replacementEnd := index + len(macroName) + offset
			p.input = append(p.input[:index], append([]byte(macroBody), p.input[replacementEnd:]...)...)
		}
	}
}

func (p *Preprocessor) replaceDefinitions() {
	keys := make([]string, 0, len(p.defines))
	for key := range p.defines {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})

	for _, key := range keys {
		p.input = bytes.ReplaceAll(p.input, []byte(key), []byte(p.defines[key]))
	}
}

func (p *Preprocessor) fillWithSpaces(start int, end int) {
	for i := start; i < end; i++ {
		p.input[i] = ' '
	}
}

func (p *Preprocessor) reTokenize() {
	l := lexer.New(string(p.input))
	p.tokens = []token.Token{}
	for {
		t := l.NextToken()
		p.tokens = append(p.tokens, t)
		if t.Kind == token.EOF {
			break
		}
	}
	p.current = 0
}

func (p *Preprocessor) peek() token.Token {
	if p.current >= len(p.tokens) {
		return token.Token{Kind: token.EOF}
	}
	return p.tokens[p.current]
}

func (p *Preprocessor) consume() token.Token {
	t := p.peek()
	p.current++
	return t
}
