package vm

import (
	"fishy/pkg/datatype"
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
)

func (m *Machine) handleArithmetic(op opcode.Opcode) {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	pos := m.position()
	rdt := datatype.DataType(m.decodeNumber("u8", pos))
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	switch op {
	case opcode.ADD_REG_LIT:
		m.applyRegLitArithmetic(rdt, func(reg, lit uint64) uint64 { return reg + lit })
	case opcode.ADD_REG_REG:
		m.applyRegRegArithmetic(func(reg0, reg1 uint64) uint64 { return reg0 + reg1 })
	case opcode.SUB_REG_LIT:
		m.applyRegLitArithmetic(rdt, func(reg, lit uint64) uint64 { return reg - lit })
	case opcode.SUB_REG_REG:
		m.applyRegRegArithmetic(func(reg0, reg1 uint64) uint64 { return reg0 - reg1 })
	case opcode.MUL_REG_LIT:
		m.applyRegLitArithmetic(rdt, func(reg, lit uint64) uint64 { return reg * lit })
	case opcode.MUL_REG_REG:
		m.applyRegRegArithmetic(func(reg0, reg1 uint64) uint64 { return reg0 * reg1 })
	case opcode.DIV_REG_LIT:
		m.applyRegLitArithmetic(rdt, func(reg, lit uint64) uint64 { return reg / lit })
	case opcode.DIV_REG_REG:
		m.applyRegRegArithmetic(func(reg0, reg1 uint64) uint64 { return reg0 / reg1 })
	}
}

func (m *Machine) applyRegLitArithmetic(dataType datatype.DataType, operation func(uint64, uint64) uint64) {
	reg := m.readRegister()
	lit := m.readLiteral(dataType)
	temp := m.getRegister(reg)
	m.setRegister(reg, operation(temp, lit))
}

func (m *Machine) applyRegRegArithmetic(operation func(uint64, uint64) uint64) {
	reg0 := m.readRegister()
	reg1 := m.readRegister()
	temp0 := m.getRegister(reg0)
	temp1 := m.getRegister(reg1)
	m.setRegister(reg0, operation(temp0, temp1))
}
