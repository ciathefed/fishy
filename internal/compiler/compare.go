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

	arg0 := instruction.Args[0]
	arg1 := instruction.Args[1]

	switch a0 := arg0.(type) {
	case *ast.Register:
		switch a1 := arg1.(type) {
		case *ast.NumberLiteral:
			num, _ := strconv.ParseInt(a1.Value, 10, 32)
			opcode := utils.Bytes2(uint16(opcode.CMP_REG_LIT))
			*section = append(*section, opcode...)
			*section = append(*section, byte(a0.Value))
			*section = append(*section, utils.Bytes4(uint32(num))...)
		case *ast.Register:
			opcode := utils.Bytes2(uint16(opcode.CMP_REG_REG))
			*section = append(*section, opcode...)
			*section = append(*section, byte(a0.Value))
			*section = append(*section, byte(a1.Value))
		default:
			return fmt.Errorf("%s expected argument #2 to be NUMBER_LITERAL or REGISTER got %s", instruction.Name, a0.String())
		}
	default:
		return fmt.Errorf("%s expected argument #1 to be REGISTER got %s", instruction.Name, a0.String())
	}
	return nil
}
