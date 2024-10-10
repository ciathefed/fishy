package compiler

import (
	"fishy/pkg/ast"
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
	"fmt"
)

func (c *Compiler) compileArithmetic(instruction *ast.Instruction) error {
	if len(instruction.Args) != 2 {
		return fmt.Errorf("%s expected 2 arguments", instruction.Name)
	}

	arg0 := instruction.Args[0]
	arg1 := instruction.Args[1]
	section := c.currentSectionBytecode()

	switch a0 := arg0.(type) {
	case *ast.Register:
		switch a1 := arg1.(type) {
		case *ast.NumberLiteral:
			num, err := ParseStringUint(a1.Value)
			if err != nil {
				return err
			}
			op, err := getArithmeticOpcode(instruction.Name, "REG_LIT")
			if err != nil {
				return err
			}
			*section = append(*section, utils.Bytes2(uint16(op))...)
			*section = append(*section, byte(instruction.DataType))
			*section = append(*section, byte(a0.Value))
			*section = append(*section, instruction.DataType.MakeBytes(num)...)
		case *ast.Register:
			op, err := getArithmeticOpcode(instruction.Name, "REG_REG")
			if err != nil {
				return err
			}
			*section = append(*section, utils.Bytes2(uint16(op))...)
			*section = append(*section, byte(instruction.DataType))
			*section = append(*section, byte(a0.Value))
			*section = append(*section, byte(a1.Value))
		case *ast.Identifier:
			op, err := getArithmeticOpcode(instruction.Name, "REG_LIT")
			if err != nil {
				return err
			}
			*section = append(*section, utils.Bytes2(uint16(op))...)
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
			op, err := getArithmeticOpcode(instruction.Name, "REG_AOF")
			if err != nil {
				return err
			}
			*section = append(*section, utils.Bytes2(uint16(op))...)
			*section = append(*section, byte(instruction.DataType))
			*section = append(*section, byte(a0.Value))

			index := a1.Value.Index()
			switch value := a1.Value.(type) {
			case *ast.NumberLiteral:
				num, err := ParseStringUint(value.Value)
				if err != nil {
					return err
				}
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
				num, err := ParseStringUint(value.Right.Value)
				if err != nil {
					return err
				}
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
				num, err := ParseStringUint(value.Right.Value)
				if err != nil {
					return err
				}
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
				return fmt.Errorf("%s expected argument #2 to be ADDRESS_OF[REGISTER], ADDRESS_OF[NUMBER], ADDRESS_OF[IDENTIFIER], ADDRESS_OF[REGISTER_OFFSET], or ADDRESS_OF[LABEL_OFFSET] got ADDRESS_OF[%s]", instruction.Name, value.String())
			}
		default:
			return fmt.Errorf("%s expected argument #2 to be REGISTER, NUMBER, IDENTIFIER, OR ADDRESS_OF got %T", instruction.Name, a1)
		}
	default:
		return fmt.Errorf("%s expected argument #1 to be REGISTER got %T", instruction.Name, a0)
	}

	return nil
}

func getArithmeticOpcode(name, kind string) (opcode.Opcode, error) {
	opcodes := map[string]map[string]opcode.Opcode{
		"add": {"REG_LIT": opcode.ADD_REG_LIT, "REG_REG": opcode.ADD_REG_REG, "REG_AOF": opcode.ADD_REG_AOF},
		"sub": {"REG_LIT": opcode.SUB_REG_LIT, "REG_REG": opcode.SUB_REG_REG, "REG_AOF": opcode.SUB_REG_AOF},
		"mul": {"REG_LIT": opcode.MUL_REG_LIT, "REG_REG": opcode.MUL_REG_REG, "REG_AOF": opcode.MUL_REG_AOF},
		"div": {"REG_LIT": opcode.DIV_REG_LIT, "REG_REG": opcode.DIV_REG_REG, "REG_AOF": opcode.DIV_REG_AOF},
	}

	ops, found := opcodes[name]
	if !found {
		return 0, fmt.Errorf("unknown arithmetic instruction %s", name)
	}

	op, found := ops[kind]
	if !found {
		return 0, fmt.Errorf("unknown argument combination for %s", name)
	}

	return op, nil
}
