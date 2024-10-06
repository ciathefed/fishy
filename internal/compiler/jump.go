package compiler

import (
	"fishy/pkg/ast"
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

	switch a := arg.(type) {
	case *ast.Register:
		return c.compileJumpReg(instruction, a)
	case *ast.NumberLiteral:
		return c.compileJumpLit(instruction, a)
	case *ast.Identifier:
		return c.compileJumpAdr(instruction, a)
	default:
		return fmt.Errorf("%s expected argument #1 to be REGISTER, NUMBER_LITERAL, IDENTIFIER got %s", instruction.Name, a.String())
	}
}

func (c *Compiler) compileJumpLit(instruction *ast.Instruction, number *ast.NumberLiteral) error {
	section := c.currentSectionBytecode()

	var op opcode.Opcode
	switch instruction.Name {
	case "jmp":
		op = opcode.JMP_LIT
	case "jeq":
		op = opcode.JEQ_LIT
	case "jne":
		op = opcode.JNE_LIT
	case "jlt":
		op = opcode.JLT_LIT
	case "jgt":
		op = opcode.JGT_LIT
	case "jle":
		op = opcode.JLE_LIT
	case "jge":
		op = opcode.JGE_LIT
	}

	num, _ := strconv.ParseInt(number.Value, 10, 32)
	*section = append(*section, utils.Bytes2(uint16(op))...)
	*section = append(*section, utils.Bytes4(uint32(num))...)

	return nil
}

func (c *Compiler) compileJumpReg(instruction *ast.Instruction, register *ast.Register) error {
	section := c.currentSectionBytecode()

	var op opcode.Opcode
	switch instruction.Name {
	case "jmp":
		op = opcode.JMP_REG
	case "jeq":
		op = opcode.JEQ_REG
	case "jne":
		op = opcode.JNE_REG
	case "jlt":
		op = opcode.JLT_REG
	case "jgt":
		op = opcode.JGT_REG
	case "jle":
		op = opcode.JLE_REG
	case "jge":
		op = opcode.JGE_REG
	}

	*section = append(*section, utils.Bytes2(uint16(op))...)
	*section = append(*section, byte(register.Value))

	return nil
}

func (c *Compiler) compileJumpAdr(instruction *ast.Instruction, identifier *ast.Identifier) error {
	section := c.currentSectionBytecode()

	var op opcode.Opcode
	switch instruction.Name {
	case "jmp":
		op = opcode.JMP_LIT
	case "jeq":
		op = opcode.JEQ_LIT
	case "jne":
		op = opcode.JNE_LIT
	case "jlt":
		op = opcode.JLT_LIT
	case "jgt":
		op = opcode.JGT_LIT
	case "jle":
		op = opcode.JLE_LIT
	case "jge":
		op = opcode.JGE_LIT
	}

	*section = append(*section, utils.Bytes2(uint16(op))...)
	c.fixups = append(c.fixups, Fixup{
		addr:    len(*section),
		section: c.currentSection,
		label:   identifier.Value,
	})
	*section = append(*section, []byte{0xDE, 0xAD, 0xBE, 0xEF}...)

	return nil
}
