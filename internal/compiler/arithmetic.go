package compiler

import (
	"fishy/pkg/ast"
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
	"fmt"
	"strconv"
)

func (c *Compiler) compileArithmetic(instruction *ast.Instruction) error {
	section := c.currentSectionBytecode()

	if len(instruction.Args) != 2 {
		return fmt.Errorf("%s expected 2 arguments", instruction.Name)
	}

	arg0 := instruction.Args[0]
	arg1 := instruction.Args[1]
	bytecode := []byte{}
	kind := "REG_LIT"

	switch a0 := arg0.(type) {
	case *ast.Register:
		switch a1 := arg1.(type) {
		case *ast.NumberLiteral:
			num, _ := strconv.ParseInt(a1.Value, 10, 32)
			kind = "REG_LIT"
			bytecode = append(bytecode, byte(a0.Value))
			bytecode = append(bytecode, utils.Bytes4(uint32(num))...)
		case *ast.Register:
			kind = "REG_REG"
			bytecode = append(bytecode, byte(a0.Value))
			bytecode = append(bytecode, byte(a1.Value))
		default:
			return fmt.Errorf("%s expected argument #2 to be REGISTER or NUMBER_LITERAL got %s", instruction.Name, a0.String())
		}
	default:
		return fmt.Errorf("%s expected argument #1 to be REGISTER got %s", instruction.Name, a0.String())
	}

	switch instruction.Name {
	case "add":
		if kind == "REG_LIT" {
			bytecode = append(utils.Bytes2(uint16(opcode.ADD_REG_LIT)), bytecode...)
		} else {
			bytecode = append(utils.Bytes2(uint16(opcode.ADD_REG_REG)), bytecode...)
		}
	case "sub":
		if kind == "REG_LIT" {
			bytecode = append(utils.Bytes2(uint16(opcode.SUB_REG_LIT)), bytecode...)
		} else {
			bytecode = append(utils.Bytes2(uint16(opcode.SUB_REG_REG)), bytecode...)
		}
	case "mul":
		if kind == "REG_LIT" {
			bytecode = append(utils.Bytes2(uint16(opcode.MUL_REG_LIT)), bytecode...)
		} else {
			bytecode = append(utils.Bytes2(uint16(opcode.MUL_REG_REG)), bytecode...)
		}
	case "div":
		if kind == "REG_LIT" {
			bytecode = append(utils.Bytes2(uint16(opcode.DIV_REG_LIT)), bytecode...)
		} else {
			bytecode = append(utils.Bytes2(uint16(opcode.DIV_REG_REG)), bytecode...)
		}
	}

	*section = append(*section, bytecode...)

	return nil
}
