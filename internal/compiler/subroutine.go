package compiler

import (
	"fishy/pkg/ast"
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
	"fmt"
	"strconv"
)

func (c *Compiler) compileCall(instruction *ast.Instruction) error {
	section := c.currentSectionBytecode()

	if len(instruction.Args) != 1 {
		return fmt.Errorf("%s expected 1 argument", instruction.Name)
	}

	arg := instruction.Args[0]

	switch a := arg.(type) {
	case *ast.NumberLiteral:
		num, _ := strconv.ParseInt(a.Value, 10, 64)
		opcode := utils.Bytes2(uint16(opcode.CALL_LIT))
		*section = append(*section, opcode...)
		*section = append(*section, utils.Bytes4(uint32(num))...)
	case *ast.Identifier:
		opcode := utils.Bytes2(uint16(opcode.CALL_LIT))
		*section = append(*section, opcode...)
		c.fixups = append(c.fixups, Fixup{
			addr:    len(*section),
			section: c.currentSection,
			label:   a.Value,
		})
		*section = append(*section, []byte{0xDE, 0xAD, 0xBE, 0xEF}...)
	default:
		return fmt.Errorf("%s expected argument #1 to be NUMBER or IDENTIFIER got %s", instruction.Name, a.String())
	}

	return nil
}

func (c *Compiler) compileRet() error {
	opcode := utils.Bytes2(uint16(opcode.RET))
	section := c.currentSectionBytecode()
	*section = append(*section, opcode...)
	return nil
}
