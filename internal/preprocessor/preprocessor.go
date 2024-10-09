package preprocessor

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type PreProcessor struct {
	lines         []Line
	filePath      string
	debug         bool
	constants     map[string]*Constant
	macros        map[string]*Macro
	includedFiles map[string]*LocationInformation
}

func New(filePath string, debug bool) (*PreProcessor, error) {
	p := &PreProcessor{
		debug:         debug,
		constants:     make(map[string]*Constant),
		macros:        make(map[string]*Macro),
		includedFiles: make(map[string]*LocationInformation),
	}
	if err := p.readSourceFile(filePath); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *PreProcessor) Process() error {
	var currentMacro *Macro

	for linesIdx := 0; linesIdx < len(p.lines); linesIdx++ {
		line := &p.lines[linesIdx]

		line.data = strings.TrimSpace(line.data)

		inString, wasComment := false, false

		for i := 0; i < len(line.data); i++ {
			if line.data[i] == '"' {
				inString = !inString
			} else if !inString && line.data[i] == ';' {
				line.data = line.data[:i]
				wasComment = true
				break
			}
		}

		if wasComment {
			line.data = strings.TrimRight(line.data, " ")
			if line.data == "" {
				p.lines = append(p.lines[:linesIdx], p.lines[linesIdx+1:]...)
				linesIdx--
				continue
			}
		}

		if strings.HasPrefix(line.data, ".section") {
			continue
		}

		if len(line.data) > 0 && line.data[0] == '#' {
			i := 1
			j := i

			for i < len(line.data) && isAlphaNumeric(line.data[i]) {
				i++
			}
			directive := strings.ToLower(line.data[j:i])

			if p.debug {
				fmt.Printf("[%d:%d] DIRECTIVE: '%s'\n", line.n, j, directive)
			}

			if currentMacro != nil {
				if directive == "end" {
					if p.debug {
						fmt.Printf("\tEnd definition of macro with %d lines\n", len(currentMacro.lines))
					}
					currentMacro = nil
				} else {
					return fmt.Errorf("unknown/invalid directive in macro body: %s", directive)
				}
			} else {
				switch directive {
				case "define":
					skipAlpha(line.data, &i)
					skipWhitespace(line.data, &i)

					j = i
					skipNonWhitespace(line.data, &i)
					constant := line.data[j:i]

					if p.debug {
						fmt.Printf("\tConstant: %s", constant)
					}

					skipWhitespace(line.data, &i)
					value := line.data[i:]

					if p.debug {
						fmt.Printf("; Value = \"%s\"\n", value)
					}

					if existingConst, exists := p.constants[constant]; exists {
						if p.debug {
							fmt.Printf("re-definition of constant %s (previously defined at %d:%d)\n", constant, existingConst.line, existingConst.col)
						}
						existingConst.value = value
						existingConst.line = line.n
						existingConst.col = j
					} else {

						p.constants[constant] = &Constant{line: line.n, col: j, value: value}
					}

				case "include":
					skipNonWhitespace(line.data, &i)
					skipWhitespace(line.data, &i)
					filePath := line.data[i:]
					filePath = strings.Trim(filePath, "\"")

					if p.debug {
						fmt.Printf("\tFile path: '%s'\n", filePath)
					}

					includeData := PreProcessor{
						debug:         p.debug,
						constants:     make(map[string]*Constant),
						macros:        make(map[string]*Macro),
						includedFiles: make(map[string]*LocationInformation),
					}
					err := includeData.readSourceFile(filePath)
					if err != nil {
						return err
					}

					if _, found := p.includedFiles[filePath]; found {
						return fmt.Errorf("circular include: %s", filePath)
					}

					err = includeData.Process()
					if err != nil {
						return err
					}

					p.Merge(&includeData, linesIdx+1)
				case "macro":
					skipAlpha(line.data, &i)
					skipWhitespace(line.data, &i)

					j = i
					macroNameIndex := i
					skipNonWhitespace(line.data, &i)
					macroName := line.data[j:i]

					if p.debug {
						fmt.Printf("\tName: '%s'; Args: ", macroName)
					}

					if !isValidLabelName(macroName) {
						return fmt.Errorf("invalid macro name: %s", macroName)
					}

					macroExists, exists := p.macros[macroName]
					if exists {
						if p.debug {
							fmt.Printf("re-definition of macro %s (previously defined at %d:%d)\n", macroName, macroExists.line, macroExists.col)
						}
						macroExists.line = line.n
						macroExists.col = j
					}

					var macroParams []string
					for {
						skipWhitespace(line.data, &i)
						if i == len(line.data) {
							break
						}

						j = i
						skipNonWhitespace(line.data, &i)
						param := line.data[j:i]

						if !isValidLabelName(param) {
							return fmt.Errorf("invalid parameter name: %s", param)
						}

						macroParams = append(macroParams, param)

						if p.debug {
							fmt.Printf("%s ", param)
						}

						if i > len(line.data) {
							break
						}
					}

					if macroExists == nil {
						p.macros[macroName] = &Macro{line: line.n, col: macroNameIndex, params: macroParams}
						macroExists = p.macros[macroName]
					} else {
						macroExists.params = macroParams
					}

					currentMacro = macroExists
				default:
					return fmt.Errorf("unknown directive: %s", directive)
				}
			}

			p.lines = append(p.lines[:linesIdx], p.lines[linesIdx+1:]...)
			linesIdx--

			continue
		}

		for n, c := range p.constants {
			index := 0
			for {

				index = strings.Index(line.data, n)

				if index == -1 {
					break
				}

				if p.debug {
					fmt.Printf("[%d:%d] CONSTANT: substitute symbol %s\n", line.n, index, n)
				}

				line.data = line.data[:index] + c.value + line.data[index+len(n):]

				index += len(c.value)
			}
		}

		if currentMacro != nil {
			currentMacro.lines = append(currentMacro.lines, line.data)

			p.lines = append(p.lines[:linesIdx], p.lines[linesIdx+1:]...)
			linesIdx--

			continue
		}

		i := 0
		skipNonWhitespace(line.data, &i)
		mnemonic := line.data[:i]

		if macroExists, ok := p.macros[mnemonic]; ok {
			if p.debug {
				fmt.Printf("[%d:0] CALL TO MACRO %s\n", line.n, mnemonic)
				fmt.Printf("\tArgs: ")
			}

			var arguments []string
			var j int

			for {
				skipWhitespace(line.data, &i)

				j = i
				skipToBreak(line.data, &i)

				if i == j {
					break
				}

				argument := line.data[j:i]
				arguments = append(arguments, argument)

				if p.debug {
					fmt.Printf("%s ", argument)
				}

				if i == len(line.data) {
					break
				}

				if line.data[i] == ',' {
					i++
				}
			}

			if p.debug {
				if len(arguments) == 0 {
					fmt.Printf("(none)")
				}
				fmt.Println()
			}

			if len(macroExists.params) != len(arguments) {
				return fmt.Errorf("macro %s expects %d argument(s), received %d", mnemonic, len(macroExists.params), len(arguments))
			}

			p.lines = append(p.lines[:linesIdx], p.lines[linesIdx+1:]...)

			insertIdx := linesIdx

			for _, macroLine := range macroExists.lines {
				argIndex := 0

				for _, param := range macroExists.params {
					index := 0
					arg := arguments[argIndex]

					for {

						index = strings.Index(macroLine[index:], param)
						if index == -1 {
							break
						}

						if p.debug {
							fmt.Printf("\tCol %d: EXPANSION: substitute parameter %s with value \"%s\"\n", index, param, arg)
						}

						macroLine = strings.Replace(macroLine, param, arg, 1)
						index += len(arg)
					}

					argIndex++
				}

				p.lines = append(p.lines[:insertIdx], append([]Line{{n: line.n, data: macroLine}}, p.lines[insertIdx:]...)...)
				insertIdx++

				continue
			}
		}
	}

	return nil
}

func (p *PreProcessor) Merge(other *PreProcessor, lineIndex int) {
	if lineIndex < 0 {
		lineIndex = len(p.lines) + lineIndex
		if lineIndex < 0 {
			lineIndex = 0
		}
	}
	p.lines = append(p.lines[:lineIndex], append(other.lines, p.lines[lineIndex:]...)...)

	for k, v := range other.constants {
		p.constants[k] = v
	}

	for k, v := range other.macros {
		p.macros[k] = v
	}
}

func (p *PreProcessor) Output() string {
	var output string
	for _, line := range p.lines {
		output += line.data + "\n"
	}
	return output
}

func (p *PreProcessor) readSourceFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("cannot read file %s", filename)
	}
	defer file.Close()

	p.filePath = filename

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineContent := scanner.Text()
		if lineContent != "" {
			p.lines = append(p.lines, Line{
				n:    lineNumber,
				data: lineContent,
			})
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
