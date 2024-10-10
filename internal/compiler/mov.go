package compiler

import (
	"fishy/pkg/ast"
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
	"fmt"
)

func (c *Compiler) compileMov(instruction *ast.Instruction) error {
	section := c.currentSectionBytecode()

	if len(instruction.Args) < 2 {
		return fmt.Errorf("mov expected 2 arguments")
	}

	arg0 := instruction.Args[0]
	arg1 := instruction.Args[1]

	appendOpcode := func(op opcode.Opcode) {
		*section = append(*section, utils.Bytes2(uint16(op))...)
		*section = append(*section, byte(instruction.DataType))
	}

	appendRegister := func(regValue int) {
		*section = append(*section, byte(regValue))
	}

	appendFixup := func(label string) {
		c.fixups = append(c.fixups, Fixup{
			addr:     len(*section),
			section:  c.currentSection,
			label:    label,
			dataType: instruction.DataType,
		})
		*section = append(*section, instruction.DataType.MakeBytes(0)...) // Placeholder for fixup
	}

	appendAddressOf := func(value ast.Value, index int) error {
		switch v := value.(type) {
		case *ast.NumberLiteral:
			num, err := ParseStringUint(v.Value)
			if err != nil {
				return err
			}
			*section = append(*section, byte(index))
			*section = append(*section, instruction.DataType.MakeBytes(num)...)
		case *ast.Register:
			*section = append(*section, byte(index))
			*section = append(*section, byte(v.Value))
		case *ast.Identifier:
			*section = append(*section, byte(index))
			appendFixup(v.Value)
		case *ast.RegisterOffsetNumber:
			return appendOffsetNumber(section, instruction, v, byte(index))
		case *ast.RegisterOffsetRegister:
			return appendOffsetRegister(section, v, byte(index))
		case *ast.LabelOffsetNumber, *ast.LabelOffsetRegister:
			return appendLabelOffset(section, instruction, v, appendFixup)
		default:
			return fmt.Errorf("unsupported AddressOf argument type: %s", v.String())
		}
		return nil
	}

	switch a0 := arg0.(type) {
	case *ast.Register:
		switch a1 := arg1.(type) {
		case *ast.Register:
			appendOpcode(opcode.MOV_REG_REG)
			appendRegister(a0.Value)
			appendRegister(a1.Value)
		case *ast.NumberLiteral:
			num, err := ParseStringUint(a1.Value)
			if err != nil {
				return err
			}
			appendOpcode(opcode.MOV_REG_LIT)
			appendRegister(a0.Value)
			*section = append(*section, instruction.DataType.MakeBytes(num)...)
		case *ast.Identifier:
			appendOpcode(opcode.MOV_REG_ADR)
			appendRegister(a0.Value)
			appendFixup(a1.Value)
		case *ast.AddressOf:
			appendOpcode(opcode.MOV_REG_AOF)
			appendRegister(a0.Value)
			err := appendAddressOf(a1.Value, a1.Value.Index())
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("mov expected argument #2 to be REGISTER, NUMBER, IDENTIFIER, or ADDRESS_OF got %s", arg1.String())
		}

	case *ast.AddressOf:
		var bytecode []byte
		switch a1 := arg1.(type) {
		case *ast.Register:
			appendOpcode(opcode.MOV_AOF_REG)
			bytecode = append(bytecode, byte(a1.Value))
		case *ast.NumberLiteral:
			num, err := ParseStringUint(a1.Value)
			if err != nil {
				return err
			}
			appendOpcode(opcode.MOV_AOF_LIT)
			bytecode = append(bytecode, instruction.DataType.MakeBytes(num)...)
		default:
			return fmt.Errorf("mov expected argument #2 to be REGISTER or NUMBER got %s", a1.String())
		}
		err := appendAddressOf(a0.Value, a0.Value.Index())
		if err != nil {
			return err
		}
		*section = append(*section, bytecode...)
	default:
		return fmt.Errorf("mov expected argument #1 to be REGISTER or ADDRESS_OF got %s", arg0.String())
	}
	return nil
}

func appendOffsetNumber(section *[]byte, instruction *ast.Instruction, value *ast.RegisterOffsetNumber, index byte) error {
	num, err := ParseStringUint(value.Right.Value)
	if err != nil {
		return err
	}
	*section = append(*section, byte(index), byte(value.Left.Value), byte(int(value.Operator)))
	*section = append(*section, instruction.DataType.MakeBytes(num)...)
	return nil
}

func appendOffsetRegister(section *[]byte, value *ast.RegisterOffsetRegister, index byte) error {
	*section = append(*section, byte(index), byte(value.Left.Value), byte(int(value.Operator)), byte(value.Right.Value))
	return nil
}

func appendLabelOffset(section *[]byte, instruction *ast.Instruction, value interface{}, appendFixup func(label string)) error {
	switch v := value.(type) {
	case *ast.LabelOffsetNumber:
		num, err := ParseStringUint(v.Right.Value)
		if err != nil {
			return err
		}
		appendFixup(v.Left.(*ast.Identifier).Value)
		*section = append(*section, instruction.DataType.MakeBytes(0)...)
		*section = append(*section, byte(int(v.Operator)))
		*section = append(*section, instruction.DataType.MakeBytes(num)...)
	case *ast.LabelOffsetRegister:
		appendFixup(v.Left.(*ast.Identifier).Value)
		*section = append(*section, instruction.DataType.MakeBytes(0)...)
		*section = append(*section, byte(int(v.Operator)))
	}
	return nil
}
