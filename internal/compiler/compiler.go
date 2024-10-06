package compiler

import (
	"fishy/pkg/ast"
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
	"fmt"
	"strconv"
)

type Section int

const (
	SectionText Section = iota
	SectionData
)

type Label struct {
	addr    int
	section Section
}

type Fixup struct {
	addr    int
	section Section
	label   string
}

type Compiler struct {
	statements     []ast.Statement
	header         []byte
	text           []byte
	data           []byte
	labels         map[string]Label
	fixups         []Fixup
	currentSection Section
}

func New(statements []ast.Statement) *Compiler {
	return &Compiler{
		statements:     statements,
		header:         make([]byte, 4),
		text:           make([]byte, 0),
		data:           make([]byte, 0),
		labels:         make(map[string]Label),
		fixups:         make([]Fixup, 0),
		currentSection: SectionText,
	}
}

func (c *Compiler) Compile() ([]byte, error) {
	for _, stmt := range c.statements {
		switch s := stmt.(type) {
		case *ast.Label:
			err := c.compileLabel(s)
			if err != nil {
				return nil, err
			}
		case *ast.Instruction:
			err := c.compileInstruction(s)
			if err != nil {
				return nil, err
			}
		case *ast.Sequence:
			err := c.compileSequence(s)
			if err != nil {
				return nil, err
			}
		}
	}

	c.resolveFixups()

	c.writeHeader()

	bytecode := append(c.header, c.text...)
	bytecode = append(bytecode, c.data...)
	return bytecode, nil
}

func (c *Compiler) compileLabel(label *ast.Label) error {
	addr := len(*c.currentSectionBytecode())
	c.labels[label.Name] = Label{addr: addr, section: c.currentSection}
	return nil
}

func (c *Compiler) compileInstruction(instruction *ast.Instruction) error {
	switch instruction.Name {
	case ".section":
		if len(instruction.Args) < 1 {
			return fmt.Errorf(".section expected 1 argument")
		}

		arg0 := instruction.Args[0]

		if value := arg0.(*ast.Identifier); value != nil {
			c.changeSection(value.Value)
		} else {
			return fmt.Errorf(".section expected %s got %s", "IDENTIFIER", arg0.String())
		}
	case "nop":
		opcode := utils.Bytes2(uint16(opcode.NOP))
		section := c.currentSectionBytecode()
		*section = append(*section, opcode...)
	case "hlt":
		opcode := utils.Bytes2(uint16(opcode.HLT))
		section := c.currentSectionBytecode()
		*section = append(*section, opcode...)
	case "brk":
		opcode := utils.Bytes2(uint16(opcode.BRK))
		section := c.currentSectionBytecode()
		*section = append(*section, opcode...)
	case "syscall":
		opcode := utils.Bytes2(uint16(opcode.SYSCALL))
		section := c.currentSectionBytecode()
		*section = append(*section, opcode...)
	case "mov":
		section := c.currentSectionBytecode()

		if len(instruction.Args) < 2 {
			return fmt.Errorf("mov expected 2 arguments")
		}

		arg0 := instruction.Args[0]
		arg1 := instruction.Args[1]

		switch a0 := arg0.(type) {
		case *ast.Register:
			switch a1 := arg1.(type) {
			case *ast.Register:
				opcode := utils.Bytes2(uint16(opcode.MOV_REG_REG))
				*section = append(*section, opcode...)
				*section = append(*section, byte(a0.Value))
				*section = append(*section, byte(a1.Value))
			case *ast.NumberLiteral:
				num, _ := strconv.ParseInt(a1.Value, 10, 32)
				opcode := utils.Bytes2(uint16(opcode.MOV_REG_LIT))
				*section = append(*section, opcode...)
				*section = append(*section, byte(a0.Value))
				*section = append(*section, utils.Bytes4(uint32(num))...)
			case *ast.Identifier:
				opcode := utils.Bytes2(uint16(opcode.MOV_REG_ADR))
				*section = append(*section, opcode...)
				*section = append(*section, byte(a0.Value))
				c.fixups = append(c.fixups, Fixup{
					addr:    len(*section),
					section: c.currentSection,
					label:   a1.Value,
				})
				*section = append(*section, []byte{0xDE, 0xAD, 0xBE, 0xEF}...)
			case *ast.AddressOf:
				opcode := utils.Bytes2(uint16(opcode.MOV_REG_AOF))
				*section = append(*section, opcode...)
				*section = append(*section, byte(a0.Value))

				index := a1.Value.Index()
				switch value := a1.Value.(type) {
				case *ast.NumberLiteral:
					num, _ := strconv.ParseInt(value.Value, 10, 32)
					*section = append(*section, byte(index))
					*section = append(*section, utils.Bytes4(uint32(num))...)
				case *ast.Register:
					*section = append(*section, byte(index))
					*section = append(*section, byte(value.Value))
				case *ast.Identifier:
					*section = append(*section, byte(index))
					c.fixups = append(c.fixups, Fixup{
						addr:    len(*section),
						section: c.currentSection,
						label:   value.Value,
					})
					*section = append(*section, []byte{0xDE, 0xAD, 0xBE, 0xEF}...)
				default:
					return fmt.Errorf("mov expected #2 argument to be ADDRESS_OF[REGISTER], ADDRESS_OF[NUMBER_LITERAL], or ADDRESS_OF[IDENTIFIER] got ADDRESS_OF[%s]", value.String())
				}
			default:
				return fmt.Errorf("mov expected #2 argument to be REGISTER, NUMBER_LITERAL, IDENTIFIER, or ADDRESS_OF got %s", arg1.String())
			}

		case *ast.AddressOf:
			opcode := utils.Bytes2(uint16(opcode.MOV_AOF_REG))
			*section = append(*section, opcode...)

			index := a0.Value.Index()
			switch value := a0.Value.(type) {
			case *ast.NumberLiteral:
				num, _ := strconv.ParseInt(value.Value, 10, 32)
				*section = append(*section, byte(index))
				*section = append(*section, utils.Bytes4(uint32(num))...)
			case *ast.Register:
				*section = append(*section, byte(index))
				*section = append(*section, byte(value.Value))
			case *ast.Identifier:
				*section = append(*section, byte(index))
				c.fixups = append(c.fixups, Fixup{
					addr:    len(*section),
					section: c.currentSection,
					label:   value.Value,
				})
				*section = append(*section, []byte{0xDE, 0xAD, 0xBE, 0xEF}...)
			default:
				return fmt.Errorf("mov expected #1 argument to be ADDRESS_OF[REGISTER], ADDRESS_OF[NUMBER_LITERAL], or ADDRESS_OF[IDENTIFIER] got ADDRESS_OF[%s]", value.String())
			}

			switch a1 := arg1.(type) {
			case *ast.Register:
				*section = append(*section, byte(a1.Value))
			default:
				return fmt.Errorf("mov expected #2 argument to be REGISTER got %s", a1.String())
			}
		default:
			return fmt.Errorf("mov expected #1 argument to be REGISTER or ADDRESS_OF got %s", arg0.String())
		}
	default:
		return fmt.Errorf("unknown instruction: %s", instruction.Name)
	}
	return nil
}

func (c *Compiler) compileSequence(sequence *ast.Sequence) error {
	switch sequence.Name {
	case "db":
		section := c.currentSectionBytecode()
		bytecode := []byte{}

		for i, value := range sequence.Values {
			switch v := value.(type) {
			case *ast.StringLiteral:
				bytecode = append(bytecode, []byte(v.Value)...)
			case *ast.NumberLiteral:
				num, _ := strconv.ParseInt(v.Value, 10, 8)
				bytecode = append(bytecode, byte(num))
			default:
				return fmt.Errorf("mov expected argument #%d to be STRING_LITERAL or NUMBER_LITERAL got %s", i, v.String())
			}
		}

		*section = append(*section, bytecode...)
	default:
		return fmt.Errorf("unknown sequence: %s", sequence.Name)
	}
	return nil
}

func (c *Compiler) changeSection(section string) {
	switch section {
	case "text":
		c.currentSection = SectionText
	case "data":
		c.currentSection = SectionData
	default:
		panic(fmt.Sprintf("unknown section: %s", section))
	}
}

func (c *Compiler) resolveFixups() {
	for _, fixup := range c.fixups {
		if label, ok := c.labels[fixup.label]; ok {
			currentSection := fixup.section
			fixupAddr := fixup.addr

			if currentSection == label.section {
				bytes := utils.Bytes4(uint32(label.addr))

				if fixup.section == SectionText {
					for i := 0; i < 4; i++ {
						c.text[(fixupAddr + i)] = bytes[i]
					}
				} else {
					for i := 0; i < 4; i++ {
						c.data[(fixupAddr + i)] = bytes[i]
					}
				}
			} else {
				offset := c.getAddrOffset(label.addr, label.section)
				bytes := utils.Bytes4(uint32(offset))

				if fixup.section == SectionText {
					for i := 0; i < 4; i++ {
						c.text[(fixupAddr + i)] = bytes[i]
					}
				} else {
					for i := 0; i < 4; i++ {
						c.data[(fixupAddr + i)] = bytes[i]
					}
				}
			}
		} else {
			panic(fmt.Sprintf("label not defined: %s", fixup.label))
		}
	}
}

func (c *Compiler) writeHeader() {
	if label, ok := c.labels["_start"]; ok {
		addr := c.getAddrOffset(label.addr, label.section)
		bytes := utils.Bytes4(uint32(addr))
		copy(c.header, bytes[:])
	}
}

func (c *Compiler) getAddrOffset(addr int, section Section) int {
	if section == SectionText {
		return addr
	} else {
		textSectionSize := len(c.text)
		return textSectionSize + addr
	}
}

func (c *Compiler) currentSectionBytecode() *[]byte {
	if c.currentSection == SectionText {
		return &c.text
	} else {
		return &c.data
	}
}
