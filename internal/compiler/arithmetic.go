package compiler

import (
	"fishy/pkg/ast"
	"fishy/pkg/datatype"
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
	"fmt"
)

func (c *Compiler) compileArithmetic(instruction *ast.Instruction) error {
	section := c.currentSectionBytecode()

	if len(instruction.Args) != 2 {
		return fmt.Errorf("%s expected 2 arguments", instruction.Name)
	}

	bytecode, kind, err := getArithmeticArgsBytecode(instruction.DataType, instruction.Args[0], instruction.Args[1])
	if err != nil {
		return err
	}

	op, err := getArithmeticOpcode(instruction.Name, kind)
	if err != nil {
		return err
	}
	*section = append(*section, utils.Bytes2(uint16(op))...)
	*section = append(*section, byte(instruction.DataType))
	*section = append(*section, bytecode...)

	return nil
}

func getArithmeticArgsBytecode(dataType datatype.DataType, arg0, arg1 interface{}) ([]byte, string, error) {
	bytecode := []byte{}
	kind := "REG_LIT"

	switch a0 := arg0.(type) {
	case *ast.Register:
		switch a1 := arg1.(type) {
		case *ast.NumberLiteral:
			num, err := ParseStringUint(a1.Value)
			if err != nil {
				return nil, "", err
			}
			kind = "REG_LIT"
			bytecode = append(bytecode, byte(a0.Value))
			bytecode = append(bytecode, dataType.MakeBytes(num)...)
		case *ast.Register:
			kind = "REG_REG"
			bytecode = append(bytecode, byte(a0.Value))
			bytecode = append(bytecode, byte(a1.Value))
		default:
			return nil, "", fmt.Errorf("expected argument #2 to be REGISTER or NUMBER got %T", a1)
		}
	default:
		return nil, "", fmt.Errorf("expected argument #1 to be REGISTER got %T", a0)
	}

	return bytecode, kind, nil
}

func getArithmeticOpcode(name, kind string) (opcode.Opcode, error) {
	opcodes := map[string]map[string]opcode.Opcode{
		"add": {"REG_LIT": opcode.ADD_REG_LIT, "REG_REG": opcode.ADD_REG_REG},
		"sub": {"REG_LIT": opcode.SUB_REG_LIT, "REG_REG": opcode.SUB_REG_REG},
		"mul": {"REG_LIT": opcode.MUL_REG_LIT, "REG_REG": opcode.MUL_REG_REG},
		"div": {"REG_LIT": opcode.DIV_REG_LIT, "REG_REG": opcode.DIV_REG_REG},
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
