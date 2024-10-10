package vm

import (
	"fishy/pkg/datatype"
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
)

func (m *Machine) handleBitwise(thread *Thread, op opcode.Opcode) {
	m.incRegister(thread, utils.RegisterToIndex("ip"), 2)

	pos := m.position(thread)
	rdt := datatype.DataType(m.decodeNumber("byte", pos))
	m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

	switch op {
	case opcode.AND_REG_LIT:
		m.applyRegLitBitwise(thread, rdt, func(reg, lit uint64) uint64 { return reg & lit })
	case opcode.AND_REG_REG:
		m.applyRegRegBitwise(thread, func(reg0, reg1 uint64) uint64 { return reg0 & reg1 })
	case opcode.OR_REG_LIT:
		m.applyRegLitBitwise(thread, rdt, func(reg, lit uint64) uint64 { return reg | lit })
	case opcode.OR_REG_REG:
		m.applyRegRegBitwise(thread, func(reg0, reg1 uint64) uint64 { return reg0 | reg1 })
	case opcode.XOR_REG_LIT:
		m.applyRegLitBitwise(thread, rdt, func(reg, lit uint64) uint64 { return reg ^ lit })
	case opcode.XOR_REG_REG:
		m.applyRegRegBitwise(thread, func(reg0, reg1 uint64) uint64 { return reg0 ^ reg1 })
	case opcode.SHL_REG_LIT:
		m.applyRegLitBitwise(thread, rdt, func(reg, lit uint64) uint64 { return reg << lit })
	case opcode.SHL_REG_REG:
		m.applyRegRegBitwise(thread, func(reg0, reg1 uint64) uint64 { return reg0 << reg1 })
	case opcode.SHR_REG_LIT:
		m.applyRegLitBitwise(thread, rdt, func(reg, lit uint64) uint64 { return reg >> lit })
	case opcode.SHR_REG_REG:
		m.applyRegRegBitwise(thread, func(reg0, reg1 uint64) uint64 { return reg0 >> reg1 })
	}
}

func (m *Machine) applyRegLitBitwise(thread *Thread, dataType datatype.DataType, operation func(uint64, uint64) uint64) {
	reg := m.readRegister(thread)
	lit := m.readLiteral(thread, dataType)
	temp := m.getRegister(thread, reg)
	m.setRegister(thread, reg, operation(temp, lit))
}

func (m *Machine) applyRegRegBitwise(thread *Thread, operation func(uint64, uint64) uint64) {
	reg0 := m.readRegister(thread)
	reg1 := m.readRegister(thread)
	temp0 := m.getRegister(thread, reg0)
	temp1 := m.getRegister(thread, reg1)
	m.setRegister(thread, reg0, operation(temp0, temp1))
}
