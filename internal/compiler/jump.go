package compiler

import (
	"fishy/pkg/ast"
	"fishy/pkg/datatype"
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
	"fmt"
	"strconv"
)

func (c *Compiler) compileJump(instruction *ast.Instruction) error {
	if len(instruction.Args) != 1 {
		return fmt.Errorf("%s expected 1 argument", instruction.Name)
	}

	arg := instruction.Args[0]

	op, err := getJumpOpcode(instruction.Name)
	if err != nil {
		return err
	}

	switch a := arg.(type) {
	case *ast.Register:
		return c.compileJumpReg(op.reg, a)
	case *ast.NumberLiteral:
		return c.compileJumpLit(op.lit, a)
	case *ast.Identifier:
		return c.compileJumpAdr(op.lit, a)
	default:
		return fmt.Errorf("%s expected argument #1 to be REGISTER, NUMBER, IDENTIFIER got %s", instruction.Name, a.String())
	}
}

func (c *Compiler) compileJumpLit(op opcode.Opcode, number *ast.NumberLiteral) error {
	section := c.currentSectionBytecode()
	num, _ := strconv.ParseUint(number.Value, 10, 64)
	*section = append(*section, utils.Bytes2(uint16(op))...)
	*section = append(*section, utils.Bytes8(num)...)
	return nil
}

func (c *Compiler) compileJumpReg(op opcode.Opcode, register *ast.Register) error {
	section := c.currentSectionBytecode()
	*section = append(*section, utils.Bytes2(uint16(op))...)
	*section = append(*section, byte(register.Value))
	return nil
}

func (c *Compiler) compileJumpAdr(op opcode.Opcode, identifier *ast.Identifier) error {
	section := c.currentSectionBytecode()
	*section = append(*section, utils.Bytes2(uint16(op))...)
	c.fixups = append(c.fixups, Fixup{
		addr:     len(*section),
		section:  c.currentSection,
		label:    identifier.Value,
		dataType: datatype.UNSET,
	})
	*section = append(*section, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}...)
	return nil
}

func getJumpOpcode(name string) (struct{ reg, lit opcode.Opcode }, error) {
	opcodes := map[string]struct{ reg, lit opcode.Opcode }{
		"jmp": {opcode.JMP_REG, opcode.JMP_LIT},
		"jeq": {opcode.JEQ_REG, opcode.JEQ_LIT},
		"jne": {opcode.JNE_REG, opcode.JNE_LIT},
		"jlt": {opcode.JLT_REG, opcode.JLT_LIT},
		"jgt": {opcode.JGT_REG, opcode.JGT_LIT},
		"jle": {opcode.JLE_REG, opcode.JLE_LIT},
		"jge": {opcode.JGE_REG, opcode.JGE_LIT},
	}

	op, found := opcodes[name]
	if !found {
		return op, fmt.Errorf("unknown jump instruction %s", name)
	}
	return op, nil
}
