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
			*section = append(*section, byte(instruction.DataType))
			*section = append(*section, byte(a0.Value))
			*section = append(*section, byte(a1.Value))
		case *ast.NumberLiteral:
			num, _ := strconv.ParseUint(a1.Value, 10, 64)
			opcode := utils.Bytes2(uint16(opcode.MOV_REG_LIT))
			*section = append(*section, opcode...)
			*section = append(*section, byte(instruction.DataType))
			*section = append(*section, byte(a0.Value))
			*section = append(*section, instruction.DataType.MakeBytes(num)...)
		case *ast.Identifier:
			opcode := utils.Bytes2(uint16(opcode.MOV_REG_ADR))
			*section = append(*section, opcode...)
			*section = append(*section, byte(instruction.DataType))
			*section = append(*section, byte(a0.Value))
			c.fixups = append(c.fixups, Fixup{
				addr:     len(*section),
				section:  c.currentSection,
				label:    a1.Value,
				dataType: instruction.DataType,
			})
			*section = append(*section, instruction.DataType.MakeBytes(0)...)
		case *ast.AddressOf:
			opcode := utils.Bytes2(uint16(opcode.MOV_REG_AOF))
			*section = append(*section, opcode...)
			*section = append(*section, byte(instruction.DataType))
			*section = append(*section, byte(a0.Value))

			index := a1.Value.Index()
			switch value := a1.Value.(type) {
			case *ast.NumberLiteral:
				num, _ := strconv.ParseUint(value.Value, 10, 64)
				*section = append(*section, byte(index))
				*section = append(*section, instruction.DataType.MakeBytes(num)...)
			case *ast.Register:
				*section = append(*section, byte(index))
				*section = append(*section, byte(value.Value))
			case *ast.Identifier:
				*section = append(*section, byte(index))
				c.fixups = append(c.fixups, Fixup{
					addr:     len(*section),
					section:  c.currentSection,
					label:    value.Value,
					dataType: instruction.DataType,
				})
				*section = append(*section, instruction.DataType.MakeBytes(0)...)
			case *ast.RegisterOffset:
				num, _ := strconv.ParseUint(value.Right.Value, 10, 64)
				*section = append(*section, byte(index))
				*section = append(*section, byte(value.Left.Value))
				*section = append(*section, byte(int(value.Operator)))
				*section = append(*section, instruction.DataType.MakeBytes(num)...)
			case *ast.LabelOffset:
				num, _ := strconv.ParseUint(value.Right.Value, 10, 64)
				*section = append(*section, byte(index))
				c.fixups = append(c.fixups, Fixup{
					addr:     len(*section),
					section:  c.currentSection,
					label:    value.Left.(*ast.Identifier).Value,
					dataType: instruction.DataType,
				})
				*section = append(*section, instruction.DataType.MakeBytes(0)...)
				*section = append(*section, byte(int(value.Operator)))
				*section = append(*section, instruction.DataType.MakeBytes(num)...)
			default:
				return fmt.Errorf("mov expected argument #2 to be ADDRESS_OF[REGISTER], ADDRESS_OF[NUMBER], ADDRESS_OF[IDENTIFIER], ADDRESS_OF[REGISTER_OFFSET], or ADDRESS_OF[LABEL_OFFSET] got ADDRESS_OF[%s]", value.String())
			}
		default:
			return fmt.Errorf("mov expected argument #2 to be REGISTER, NUMBER, IDENTIFIER, or ADDRESS_OF got %s", arg1.String())
		}

	case *ast.AddressOf:
		bytecode := []byte{}
		switch a1 := arg1.(type) {
		case *ast.Register:
			opcode := utils.Bytes2(uint16(opcode.MOV_AOF_REG))
			*section = append(*section, opcode...)
			*section = append(*section, byte(instruction.DataType))
			bytecode = append(bytecode, byte(a1.Value))
		case *ast.NumberLiteral:
			opcode := utils.Bytes2(uint16(opcode.MOV_AOF_LIT))
			*section = append(*section, opcode...)
			*section = append(*section, byte(instruction.DataType))
			num, _ := strconv.ParseUint(a1.Value, 10, 64)
			bytecode = append(bytecode, instruction.DataType.MakeBytes(num)...)
		default:
			return fmt.Errorf("mov expected argument #2 to be REGISTER or NUMBER got %s", a1.String())
		}

		index := a0.Value.Index()
		switch value := a0.Value.(type) {
		case *ast.NumberLiteral:
			num, _ := strconv.ParseUint(value.Value, 10, 64)
			*section = append(*section, byte(index))
			*section = append(*section, instruction.DataType.MakeBytes(num)...)
		case *ast.Register:
			*section = append(*section, byte(index))
			*section = append(*section, byte(value.Value))
		case *ast.Identifier:
			*section = append(*section, byte(index))
			c.fixups = append(c.fixups, Fixup{
				addr:     len(*section),
				section:  c.currentSection,
				label:    value.Value,
				dataType: instruction.DataType,
			})
			*section = append(*section, instruction.DataType.MakeBytes(0)...)
		case *ast.RegisterOffset:
			num, _ := strconv.ParseUint(value.Right.Value, 10, 64)
			*section = append(*section, byte(index))
			*section = append(*section, byte(value.Left.Value))
			*section = append(*section, byte(int(value.Operator)))
			*section = append(*section, instruction.DataType.MakeBytes(num)...)
		case *ast.LabelOffset:
			num, _ := strconv.ParseUint(value.Right.Value, 10, 64)
			*section = append(*section, byte(index))
			c.fixups = append(c.fixups, Fixup{
				addr:     len(*section),
				section:  c.currentSection,
				label:    value.Left.(*ast.Identifier).Value,
				dataType: instruction.DataType,
			})
			*section = append(*section, instruction.DataType.MakeBytes(0)...)
			*section = append(*section, byte(int(value.Operator)))
			*section = append(*section, instruction.DataType.MakeBytes(num)...)
		default:
			return fmt.Errorf("mov expected argument #1 to be ADDRESS_OF[REGISTER], ADDRESS_OF[NUMBER], ADDRESS_OF[IDENTIFIER], ADDRESS_OF[REGISTER_OFFSET], or ADDRESS_OF[LABEL_OFFSET] got ADDRESS_OF[%s]", value.String())
		}

		*section = append(*section, bytecode...)
	default:
		return fmt.Errorf("mov expected argument #1 to be REGISTER or ADDRESS_OF got %s", arg0.String())
	}
	return nil
}
