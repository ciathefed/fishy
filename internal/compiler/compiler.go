package compiler

import (
	"fishy/pkg/ast"
	"fishy/pkg/datatype"
	"fishy/pkg/log"
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
	"fmt"
	"strconv"
)

type Section int

const (
	SectionText Section = iota
	SectionData
	SectionBSS
)

type Fixup struct {
	addr     int
	section  Section
	label    string
	dataType datatype.DataType
}

type Compiler struct {
	statements []ast.Statement
	lastLabel  ast.Statement

	headerStart       []byte
	headerSymbolTable []byte

	text []byte
	data []byte
	bss  []byte

	symbolTable *SymbolTable
	fixups      []Fixup

	currentSection Section
}

func New(statements []ast.Statement) *Compiler {
	return &Compiler{
		statements:        statements,
		lastLabel:         nil,
		headerStart:       make([]byte, 8),
		headerSymbolTable: make([]byte, 0),
		text:              make([]byte, 0),
		data:              make([]byte, 0),
		bss:               make([]byte, 0),
		symbolTable:       NewSymbolTable(),
		fixups:            make([]Fixup, 0),
		currentSection:    SectionText,
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
			c.lastLabel = stmt
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

	c.writeHeaderStart()

	bytecode := c.headerStart

	stStart := len(bytecode) + 16
	stEnd := stStart + len(c.headerSymbolTable)

	bytecode = append(bytecode, utils.Bytes8(uint64(stStart))...)
	bytecode = append(bytecode, utils.Bytes8(uint64(stEnd))...)

	bytecode = append(bytecode, c.headerSymbolTable...)

	bytecode = append(bytecode, c.text...)
	bytecode = append(bytecode, c.data...)
	bytecode = append(bytecode, c.bss...)
	return bytecode, nil
}

func (c *Compiler) compileLabel(label *ast.Label) error {
	addr := uint64(len(*c.currentSectionBytecode()))
	c.symbolTable.Set(label.Name, &Symbol{
		name:     label.Name,
		dataType: datatype.UNSET,
		addr:     addr,
		section:  c.currentSection,
	})
	return nil
}

func (c *Compiler) compileInstruction(instruction *ast.Instruction) error {
	switch instruction.Name {
	case ".section":
		if len(instruction.Args) != 1 {
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
		return c.compileMov(instruction)
	case "add", "sub", "mul", "div":
		return c.compileArithmetic(instruction)
	case "and", "or", "xor", "shl", "shr":
		return c.compileBitwise(instruction)
	case "cmp":
		return c.compileCompare(instruction)
	case "jmp", "jeq", "jne", "jlt", "jgt", "jle", "jge", "jz":
		return c.compileJump(instruction)
	case "push":
		return c.compilePush(instruction)
	case "pop":
		return c.compilePop(instruction)
	case "call":
		return c.compileCall(instruction)
	case "ret":
		return c.compileRet()
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

		c.updateSymbolDataType(sequence.Name, datatype.U8)

		for i, value := range sequence.Values {
			switch v := value.(type) {
			case *ast.StringLiteral:
				bytecode = append(bytecode, []byte(v.Value)...)
			case *ast.NumberLiteral:
				num, _ := strconv.ParseUint(v.Value, 10, 8)
				bytecode = append(bytecode, byte(num))
			default:
				return fmt.Errorf("%s expected argument #%d to be NUMBER got %s", sequence.Name, i, v.String())
			}
		}

		*section = append(*section, bytecode...)

	case "dw":
		section := c.currentSectionBytecode()
		bytecode := []byte{}

		c.updateSymbolDataType(sequence.Name, datatype.U16)

		for i, value := range sequence.Values {
			switch v := value.(type) {
			case *ast.NumberLiteral:
				num, _ := strconv.ParseUint(v.Value, 10, 16)
				bytecode = append(bytecode, utils.Bytes2(uint16(num))...)
			default:
				return fmt.Errorf("%s expected argument #%d to be NUMBER got %s", sequence.Name, i, v.String())
			}
		}

		*section = append(*section, bytecode...)

	case "dd":
		section := c.currentSectionBytecode()
		bytecode := []byte{}

		c.updateSymbolDataType(sequence.Name, datatype.U32)

		for i, value := range sequence.Values {
			switch v := value.(type) {
			case *ast.NumberLiteral:
				num, _ := strconv.ParseUint(v.Value, 10, 32)
				bytecode = append(bytecode, utils.Bytes4(uint32(num))...)
			default:
				return fmt.Errorf("%s expected argument #%d to be NUMBER got %s", sequence.Name, i, v.String())
			}
		}

		*section = append(*section, bytecode...)

	case "dq":
		section := c.currentSectionBytecode()
		bytecode := []byte{}

		c.updateSymbolDataType(sequence.Name, datatype.U64)

		for i, value := range sequence.Values {
			switch v := value.(type) {
			case *ast.NumberLiteral:
				num, _ := strconv.ParseUint(v.Value, 10, 64)
				bytecode = append(bytecode, utils.Bytes8(num)...)
			default:
				return fmt.Errorf("%s expected argument #%d to be NUMBER got %s", sequence.Name, i, v.String())
			}
		}

		*section = append(*section, bytecode...)
	case "resb":
		section := c.currentSectionBytecode()
		amount := 0

		c.updateSymbolDataType(sequence.Name, datatype.U8)

		for i, value := range sequence.Values {
			switch v := value.(type) {
			case *ast.NumberLiteral:
				num, _ := strconv.ParseUint(v.Value, 10, 64)
				amount += int(num)
			default:
				return fmt.Errorf("%s expected argument #%d to be NUMBER got %s", sequence.Name, i, v.String())
			}
		}

		*section = append(*section, make([]byte, amount)...)

	case "resw":
		section := c.currentSectionBytecode()
		amount := 0

		c.updateSymbolDataType(sequence.Name, datatype.U16)

		for i, value := range sequence.Values {
			switch v := value.(type) {
			case *ast.NumberLiteral:
				num, _ := strconv.ParseUint(v.Value, 10, 64)
				amount += int(num)
			default:
				return fmt.Errorf("%s expected argument #%d to be NUMBER got %s", sequence.Name, i, v.String())
			}
		}

		*section = append(*section, make([]byte, amount*2)...)

	case "resd":
		section := c.currentSectionBytecode()
		amount := 0

		c.updateSymbolDataType(sequence.Name, datatype.U32)

		for i, value := range sequence.Values {
			switch v := value.(type) {
			case *ast.NumberLiteral:
				num, _ := strconv.ParseUint(v.Value, 10, 64)
				amount += int(num)
			default:
				return fmt.Errorf("%s expected argument #%d to be NUMBER got %s", sequence.Name, i, v.String())
			}
		}

		*section = append(*section, make([]byte, amount*4)...)

	case "resq":
		section := c.currentSectionBytecode()
		amount := 0

		c.updateSymbolDataType(sequence.Name, datatype.U64)

		for i, value := range sequence.Values {
			switch v := value.(type) {
			case *ast.NumberLiteral:
				num, _ := strconv.ParseUint(v.Value, 10, 64)
				amount += int(num)
			default:
				return fmt.Errorf("%s expected argument #%d to be NUMBER got %s", sequence.Name, i, v.String())
			}
		}

		*section = append(*section, make([]byte, amount*8)...)

	default:
		return fmt.Errorf("unknown sequence: %s", sequence.Name)
	}
	return nil
}

func (c *Compiler) updateSymbolDataType(sequenceName string, newDataType datatype.DataType) {
	if stmt, ok := c.lastLabel.(*ast.Label); ok {
		if symbol := c.symbolTable.Get(stmt.Name); symbol != nil {
			if symbol.dataType != datatype.UNSET && symbol.dataType != newDataType {
				log.Warn("label data type being changed but was already set", "old", symbol.dataType, "new", newDataType, "setter", sequenceName, "label", symbol.name)
			}

			symbol.dataType = newDataType
		}
	}
}

func (c *Compiler) changeSection(section string) {
	switch section {
	case "text":
		c.currentSection = SectionText
	case "data":
		c.currentSection = SectionData
	case "bss":
		c.currentSection = SectionBSS
	default:
		log.Fatal("unknown section", "section", section)
	}
}

func (c *Compiler) resolveFixups() {
	for _, fixup := range c.fixups {
		if symbol := c.symbolTable.Get(fixup.label); symbol != nil {
			currentSection := fixup.section
			fixupAddr := fixup.addr

			if currentSection == symbol.section {
				bytes := fixup.dataType.MakeBytes(symbol.addr)

				kv := c.symbolTable.Compile(symbol.name, symbol.addr)
				c.headerSymbolTable = append(c.headerSymbolTable, kv...)

				if fixup.section == SectionText {
					for i := 0; i < fixup.dataType.Size(); i++ {
						c.text[(fixupAddr + i)] = bytes[i]
					}
				} else if fixup.section == SectionData {
					for i := 0; i < fixup.dataType.Size(); i++ {
						c.data[(fixupAddr + i)] = bytes[i]
					}
				} else {
					for i := 0; i < fixup.dataType.Size(); i++ {
						c.bss[(fixupAddr + i)] = bytes[i]
					}
				}
			} else {
				offset := c.getAddrOffset(symbol.addr, symbol.section)
				bytes := fixup.dataType.MakeBytes(offset)

				kv := c.symbolTable.Compile(symbol.name, offset)
				c.headerSymbolTable = append(c.headerSymbolTable, kv...)

				if fixup.section == SectionText {
					for i := 0; i < fixup.dataType.Size(); i++ {
						c.text[(fixupAddr + i)] = bytes[i]
					}
				} else if fixup.section == SectionData {
					for i := 0; i < fixup.dataType.Size(); i++ {
						c.data[(fixupAddr + i)] = bytes[i]
					}
				} else {
					for i := 0; i < fixup.dataType.Size(); i++ {
						c.bss[(fixupAddr + i)] = bytes[i]
					}
				}
			}
		} else {
			log.Fatal("label not defined", "label", fixup.label)
		}
	}
}

func (c *Compiler) writeHeaderStart() {
	if symbol := c.symbolTable.Get("_start"); symbol != nil {
		addr := c.getAddrOffset(symbol.addr, symbol.section)
		bytes := utils.Bytes8(addr)
		copy(c.headerStart, bytes[:])
	}
}

func (c *Compiler) getAddrOffset(addr uint64, section Section) uint64 {
	if section == SectionText {
		return addr
	} else if section == SectionData {
		textSectionSize := uint64(len(c.text))
		return textSectionSize + addr
	} else {
		textSectionSize := uint64(len(c.text))
		dataSectionSize := uint64(len(c.data))
		return textSectionSize + dataSectionSize + addr
	}
}

func (c *Compiler) currentSectionBytecode() *[]byte {
	switch c.currentSection {
	case SectionText:
		return &c.text
	case SectionData:
		return &c.data
	case SectionBSS:
		return &c.bss
	default:
		log.Fatal("unknown current section", "section", c.currentSection)
	}
	return nil
}
