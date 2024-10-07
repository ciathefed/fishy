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
		*section = append(*section, byte(a.Value))
	case *ast.NumberLiteral:
		num, _ := strconv.ParseInt(a.Value, 10, 64)
		opcode := utils.Bytes2(uint16(opcode.PUSH_LIT))
		*section = append(*section, opcode...)
		*section = append(*section, utils.Bytes4(uint32(num))...)
	default:
		return fmt.Errorf("%s expected argument #1 to be REGISTER or NUMBER got %s", instruction.Name, a.String())
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
		*section = append(*section, byte(a.Value))
	default:
		return fmt.Errorf("%s expected argument #1 to be REGISTER got %s", instruction.Name, a.String())
	}

	return nil
}
