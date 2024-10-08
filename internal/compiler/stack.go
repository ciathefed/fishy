package compiler

import (
	"fishy/pkg/ast"
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
	"fmt"
	"strconv"
)

func (c *Compiler) compilePush(instruction *ast.Instruction) error {
	section := c.currentSectionBytecode()

	if len(instruction.Args) != 1 {
		return fmt.Errorf("%s expected 1 argument", instruction.Name)
	}

	arg := instruction.Args[0]

	switch a := arg.(type) {
	case *ast.Register:
		opcode := utils.Bytes2(uint16(opcode.PUSH_REG))
		*section = append(*section, opcode...)
		*section = append(*section, byte(instruction.DataType))
		*section = append(*section, byte(a.Value))
	case *ast.NumberLiteral:
		num, _ := strconv.ParseUint(a.Value, 10, 64)
		opcode := utils.Bytes2(uint16(opcode.PUSH_LIT))
		*section = append(*section, opcode...)
		*section = append(*section, byte(instruction.DataType))
		*section = append(*section, utils.Bytes8(uint64(num))...)
	case *ast.Identifier:
		opcode := utils.Bytes2(uint16(opcode.PUSH_LIT))
		*section = append(*section, opcode...)
		*section = append(*section, byte(instruction.DataType))
		c.fixups = append(c.fixups, Fixup{
			addr:     len(*section),
			section:  c.currentSection,
			label:    a.Value,
			dataType: instruction.DataType,
		})
		*section = append(*section, instruction.DataType.MakeBytes(0)...)
	case *ast.AddressOf:
		opcode := utils.Bytes2(uint16(opcode.PUSH_AOF))
		*section = append(*section, opcode...)
		*section = append(*section, byte(instruction.DataType))

		index := a.Value.Index()
		switch value := a.Value.(type) {
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
		case *ast.RegisterOffsetNumber:
			num, _ := strconv.ParseUint(value.Right.Value, 10, 64)
			*section = append(*section, byte(index))
			*section = append(*section, byte(value.Left.Value))
			*section = append(*section, byte(int(value.Operator)))
			*section = append(*section, instruction.DataType.MakeBytes(num)...)
		case *ast.RegisterOffsetRegister:
			*section = append(*section, byte(index))
			*section = append(*section, byte(value.Left.Value))
			*section = append(*section, byte(int(value.Operator)))
			*section = append(*section, byte(value.Right.Value))
		case *ast.LabelOffsetNumber:
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
		case *ast.LabelOffsetRegister:
			*section = append(*section, byte(index))
			c.fixups = append(c.fixups, Fixup{
				addr:     len(*section),
				section:  c.currentSection,
				label:    value.Left.(*ast.Identifier).Value,
				dataType: instruction.DataType,
			})
			*section = append(*section, instruction.DataType.MakeBytes(0)...)
			*section = append(*section, byte(int(value.Operator)))
			*section = append(*section, byte(value.Right.Value))
		default:
			return fmt.Errorf("push expected argument #2 to be ADDRESS_OF[REGISTER], ADDRESS_OF[NUMBER], ADDRESS_OF[IDENTIFIER], ADDRESS_OF[REGISTER_OFFSET], or ADDRESS_OF[LABEL_OFFSET] got ADDRESS_OF[%s]", value.String())
		}
	default:
		return fmt.Errorf("%s expected argument #1 to be REGISTER, NUMBER or IDENTIFIER got %s", instruction.Name, a.String())
	}

	return nil
}

func (c *Compiler) compilePop(instruction *ast.Instruction) error {
	section := c.currentSectionBytecode()

	if len(instruction.Args) != 1 {
		return fmt.Errorf("%s expected 1 argument", instruction.Name)
	}

	arg := instruction.Args[0]

	switch a := arg.(type) {
	case *ast.Register:
		opcode := utils.Bytes2(uint16(opcode.POP_REG))
		*section = append(*section, opcode...)
		*section = append(*section, byte(instruction.DataType))
		*section = append(*section, byte(a.Value))
	case *ast.AddressOf:
		opcode := utils.Bytes2(uint16(opcode.POP_AOF))
		*section = append(*section, opcode...)
		*section = append(*section, byte(instruction.DataType))

		index := a.Value.Index()
		switch value := a.Value.(type) {
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
		case *ast.RegisterOffsetNumber:
			num, _ := strconv.ParseUint(value.Right.Value, 10, 64)
			*section = append(*section, byte(index))
			*section = append(*section, byte(value.Left.Value))
			*section = append(*section, byte(int(value.Operator)))
			*section = append(*section, instruction.DataType.MakeBytes(num)...)
		case *ast.RegisterOffsetRegister:
			*section = append(*section, byte(index))
			*section = append(*section, byte(value.Left.Value))
			*section = append(*section, byte(int(value.Operator)))
			*section = append(*section, byte(value.Right.Value))
		case *ast.LabelOffsetNumber:
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
		case *ast.LabelOffsetRegister:
			*section = append(*section, byte(index))
			c.fixups = append(c.fixups, Fixup{
				addr:     len(*section),
				section:  c.currentSection,
				label:    value.Left.(*ast.Identifier).Value,
				dataType: instruction.DataType,
			})
			*section = append(*section, instruction.DataType.MakeBytes(0)...)
			*section = append(*section, byte(int(value.Operator)))
			*section = append(*section, byte(value.Right.Value))
		default:
			return fmt.Errorf("pop expected argument #2 to be ADDRESS_OF[REGISTER], ADDRESS_OF[NUMBER], ADDRESS_OF[IDENTIFIER], ADDRESS_OF[REGISTER_OFFSET], or ADDRESS_OF[LABEL_OFFSET] got ADDRESS_OF[%s]", value.String())
		}
	default:
		return fmt.Errorf("%s expected argument #1 to be REGISTER or ADDRESS_OF got %s", instruction.Name, a.String())
	}

	return nil
}
