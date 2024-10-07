package compiler

import (
	"fishy/pkg/ast"
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
	"fmt"
	"strconv"
)

func (c *Compiler) compileCompare(instruction *ast.Instruction) error {
	section := c.currentSectionBytecode()

	if len(instruction.Args) != 2 {
		return fmt.Errorf("%s expected 2 arguments", instruction.Name)
	}

	arg0, arg1 := instruction.Args[0], instruction.Args[1]

	bytecode, err := getCompareBytecode(arg0, arg1)
	if err != nil {
		return err
	}

	*section = append(*section, bytecode...)

	return nil
}

func getCompareBytecode(arg0, arg1 ast.Value) ([]byte, error) {
	switch a0 := arg0.(type) {
	case *ast.Register:
		switch a1 := arg1.(type) {
		case *ast.NumberLiteral:
			num, err := strconv.ParseInt(a1.Value, 10, 64)
			if err != nil {
				return nil, err
			}
			return buildBytecode(opcode.CMP_REG_LIT, byte(a0.Value), utils.Bytes4(uint32(num)))
		case *ast.Register:
			return buildBytecode(opcode.CMP_REG_REG, byte(a0.Value), byte(a1.Value))
		default:
			return nil, fmt.Errorf("cmp expected argument #2 to be NUMBER or REGISTER got %s", a1.String())
		}
	default:
		return nil, fmt.Errorf("cmp expected argument #1 to be REGISTER got %s", a0.String())
	}
}

func buildBytecode(op opcode.Opcode, args ...interface{}) ([]byte, error) {
	bytecode := utils.Bytes2(uint16(op))
	for _, arg := range args {
		switch v := arg.(type) {
		case byte:
			bytecode = append(bytecode, v)
		case []byte:
			bytecode = append(bytecode, v...)
		}
	}
	return bytecode, nil
}
