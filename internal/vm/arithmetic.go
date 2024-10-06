package vm

import (
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
)

func (m *Machine) handleArithmetic(op opcode.Opcode) {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	switch op {
	case opcode.ADD_REG_LIT:
		m.applyRegLitArithmetic(func(reg, lit uint32) uint32 { return reg + lit })
	case opcode.ADD_REG_REG:
		m.applyRegRegArithmetic(func(reg0, reg1 uint32) uint32 { return reg0 + reg1 })
	case opcode.SUB_REG_LIT:
		m.applyRegLitArithmetic(func(reg, lit uint32) uint32 { return reg - lit })
	case opcode.SUB_REG_REG:
		m.applyRegRegArithmetic(func(reg0, reg1 uint32) uint32 { return reg0 - reg1 })
	case opcode.MUL_REG_LIT:
		m.applyRegLitArithmetic(func(reg, lit uint32) uint32 { return reg * lit })
	case opcode.MUL_REG_REG:
		m.applyRegRegArithmetic(func(reg0, reg1 uint32) uint32 { return reg0 * reg1 })
	case opcode.DIV_REG_LIT:
		m.applyRegLitArithmetic(func(reg, lit uint32) uint32 { return reg / lit })
	case opcode.DIV_REG_REG:
		m.applyRegRegArithmetic(func(reg0, reg1 uint32) uint32 { return reg0 / reg1 })
	}
}

func (m *Machine) applyRegLitArithmetic(operation func(uint32, uint32) uint32) {
	reg := m.readRegister()
	lit := m.readLiteral()
	temp := m.getRegister(reg)
	m.setRegister(reg, operation(temp, lit))
}

func (m *Machine) applyRegRegArithmetic(operation func(uint32, uint32) uint32) {
	reg0 := m.readRegister()
	reg1 := m.readRegister()
	temp0 := m.getRegister(reg0)
	temp1 := m.getRegister(reg1)
	m.setRegister(reg0, operation(temp0, temp1))
}
