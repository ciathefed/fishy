package compiler

import (
	"fishy/pkg/ast"
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
	"fmt"
	"strconv"
)

func (c *Compiler) compileMov(instruction *ast.Instruction) error {
	section := c.currentSectionBytecode()

	if len(instruction.Args) != 2 {
		return fmt.Errorf("mov expected 2 arguments, got %d", len(instruction.Args))
	}

	arg0, arg1 := instruction.Args[0], instruction.Args[1]

	appendBytecode := func(op opcode.Opcode, args ...interface{}) {
		opcodeBytes := utils.Bytes2(uint16(op))
		*section = append(*section, opcodeBytes...)
		for _, arg := range args {
			switch v := arg.(type) {
			case byte:
				*section = append(*section, v)
			case uint32:
				*section = append(*section, utils.Bytes4(v)...)
			default:
				if id, ok := v.(*ast.Identifier); ok {
					c.fixups = append(c.fixups, Fixup{
						addr:    len(*section),
						section: c.currentSection,
						label:   id.Value,
					})
					*section = append(*section, []byte{0xDE, 0xAD, 0xBE, 0xEF}...)
				}
			}
		}
	}

	switch a0 := arg0.(type) {
	case *ast.Register:
		switch a1 := arg1.(type) {
		case *ast.Register:
			appendBytecode(opcode.MOV_REG_REG, byte(a0.Value), byte(a1.Value))
		case *ast.NumberLiteral:
			num, _ := strconv.ParseInt(a1.Value, 10, 64)
			appendBytecode(opcode.MOV_REG_LIT, byte(a0.Value), uint32(num))
		case *ast.Identifier:
			appendBytecode(opcode.MOV_REG_ADR, byte(a0.Value), a1)
		case *ast.AddressOf:
			index := a1.Value.Index()
			appendBytecode(opcode.MOV_REG_AOF, byte(a0.Value), byte(index))
			switch value := a1.Value.(type) {
			case *ast.NumberLiteral:
				num, _ := strconv.ParseInt(value.Value, 10, 64)
				*section = append(*section, uint8(index))
				*section = append(*section, utils.Bytes4(uint32(num))...)
			case *ast.Register:
				*section = append(*section, uint8(index), byte(value.Value))
			case *ast.Identifier:
				appendBytecode(opcode.MOV_REG_AOF, byte(index), value)
			default:
				return fmt.Errorf("mov expected argument #2 to be ADDRESS_OF[REGISTER], ADDRESS_OF[NUMBER], or ADDRESS_OF[IDENTIFIER], got %s", value.String())
			}
		default:
			return fmt.Errorf("mov expected argument #2 to be REGISTER, NUMBER, IDENTIFIER, or ADDRESS_OF, got %s", arg1.String())
		}

	case *ast.AddressOf:
		index := a0.Value.Index()
		switch a1 := arg1.(type) {
		case *ast.Register:
			appendBytecode(opcode.MOV_AOF_REG, byte(a1.Value))
		case *ast.NumberLiteral:
			num, _ := strconv.ParseInt(a1.Value, 10, 64)
			appendBytecode(opcode.MOV_AOF_LIT, uint32(num))
		default:
			return fmt.Errorf("mov expected argument #2 to be REGISTER or NUMBER, got %s", a1.String())
		}

		switch value := a0.Value.(type) {
		case *ast.NumberLiteral:
			num, _ := strconv.ParseInt(value.Value, 10, 64)
			*section = append(*section, uint8(index))
			*section = append(*section, utils.Bytes4(uint32(num))...)
		case *ast.Register:
			*section = append(*section, uint8(index), byte(value.Value))
		case *ast.Identifier:
			appendBytecode(opcode.MOV_AOF_REG, value)
		default:
			return fmt.Errorf("mov expected argument #1 to be ADDRESS_OF[REGISTER], ADDRESS_OF[NUMBER], or ADDRESS_OF[IDENTIFIER], got %s", value.String())
		}

	default:
		return fmt.Errorf("mov expected argument #1 to be REGISTER or ADDRESS_OF, got %s", arg0.String())
	}
	return nil
}
