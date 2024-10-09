package compiler

import (
	"fishy/pkg/ast"
	"fishy/pkg/datatype"
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
	"fmt"
)

func (c *Compiler) compileCall(instruction *ast.Instruction) error {
	section := c.currentSectionBytecode()

	if len(instruction.Args) != 1 {
		return fmt.Errorf("%s expected 1 argument", instruction.Name)
	}

	arg := instruction.Args[0]

	switch a := arg.(type) {
	case *ast.NumberLiteral:
		num, err := ParseStringUint(a.Value)
		if err != nil {
			return err
		}
		opcode := utils.Bytes2(uint16(opcode.CALL_LIT))
		*section = append(*section, opcode...)
		*section = append(*section, utils.Bytes8(num)...)
	case *ast.Identifier:
		opcode := utils.Bytes2(uint16(opcode.CALL_LIT))
		*section = append(*section, opcode...)
		c.fixups = append(c.fixups, Fixup{
			addr:     len(*section),
			section:  c.currentSection,
			label:    a.Value,
			dataType: datatype.UNSET,
		})
		*section = append(*section, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}...)
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
